package psbt

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	btcdpsbt "github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// UTXO 最少需要这些信息（优先走 segwit，legacy 需要 nonWitnessTx）
type PsbtUTXO struct {
	TxID            string // 大端 txid
	Vout            uint32
	ValueSat        int64  // 金额，sats
	ScriptPubKeyHex string // 前序输出的 pkScript(hex)

	// Legacy P2PKH/P2SH 如要签名，必须提供 non-witness 前序原始交易（hex）
	NonWitnessTxHex string // 可留空（P2WPKH/P2TR 不需要）

	// 可选，P2SH/P2WSH 的 witnessScript(hex)
	RedeemScriptHex string // 可选，P2SH/P2WSH 的 redeemScript(hex)

	// 可选，P2WSH 的 witnessScript(hex)
	WitnessScriptHex string // 可选，P2WSH 的 witnessScript(hex)
}

// 输出给调用方/前端（包含 OKX 可用的 PSBT base64）
type BuildResult struct {
	PSBTBase64      string  `json:"psbt_base64"`        // 给 OKX
	UnsignedTxHex   string  `json:"unsigned_tx_hex"`    // 调试/核对
	Packet          *Packet `json:"-"`                  // 你自定义 psbt 结构，便于回写签名/Finalize
	EstimatedVSize  int     `json:"estimated_vsize_vb"` // 估算
	FeeSat          int64   `json:"fee_sat"`
	ChangeOutputIdx int     `json:"change_output_index"` // -1 表示没有找零
}

// CreatePSBTForOKX 生成可供 OKX 钱包签名的 PSBT（v0）
func CreatePSBTForOKX(
	inputParams *types.TxInputParams,
	selectedUTXOs []PsbtUTXO,
	network *chaincfg.Params,
) (*BuildResult, error) {

	if len(inputParams.ToAddress) == 0 || len(inputParams.ToAddress) != len(inputParams.AmountBTC) {
		return nil, errors.New("to_address 与 amount 数量需一致且 > 0")
	}
	if len(selectedUTXOs) == 0 {
		return nil, errors.New("未提供可用 UTXO")
	}
	if inputParams.ChangeAddress == "" {
		return nil, errors.New("必须提供找零地址（建议与 from 同类型）")
	}

	// 1) 目标输出
	var outputs []*wire.TxOut
	var totalSend int64
	for i, to := range inputParams.ToAddress {
		addr, err := btcutil.DecodeAddress(to, network)
		if err != nil {
			return nil, fmt.Errorf("解析收款地址失败: %s", err)
		}
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, fmt.Errorf("构造收款脚本失败: %s", err)
		}
		amtSat := int64(math.Round(inputParams.AmountBTC[i] * 1e8))
		if amtSat <= 0 {
			return nil, fmt.Errorf("非法金额: %f", inputParams.AmountBTC[i])
		}
		outputs = append(outputs, &wire.TxOut{Value: amtSat, PkScript: pkScript})
		totalSend += amtSat
	}
	// 可选 OP_RETURN（纯文本）
	if inputParams.Data != "" {
		script, _ := txscript.NewScriptBuilder().
			AddOp(txscript.OP_RETURN).
			AddData([]byte(inputParams.Data)).
			Script()
		outputs = append(outputs, &wire.TxOut{Value: 0, PkScript: script})
	}

	// 2) 组装 unsigned tx（v0）
	mtx := wire.NewMsgTx(2)
	// inputs
	var inTypes []string
	var totalIn int64
	seq := uint32(0xFFFFFFFF)
	if inputParams.Replaceable {
		seq = 0xFFFFFFFD // 明确信号 RBF
	}
	for _, u := range selectedUTXOs {
		hash, err := chainhash.NewHashFromStr(u.TxID)
		if err != nil {
			return nil, fmt.Errorf("TxID 解析失败: %v", err)
		}
		mtx.AddTxIn(&wire.TxIn{
			PreviousOutPoint: wire.OutPoint{Hash: *hash, Index: u.Vout},
			Sequence:         seq,
		})
		totalIn += u.ValueSat

		// 简单识别脚本类型（仅用于 fee 估算）
		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		inTypes = append(inTypes, detectType(spk))
	}
	// outputs 先放收款；找零稍后放
	for _, o := range outputs {
		mtx.AddTxOut(o)
	}
	// locktime（秒 -> 区块时间近似：直接写入；如非需求，可保持 0）
	mtx.LockTime = uint32(inputParams.Locktime)

	// 3) 估算 vsize 与 fee
	//   估算值：p2wpkh≈68 vB；p2tr≈58 vB；p2pkh≈148 vB；p2sh-p2wpkh≈91 vB
	estVSize := estimateVSize(inTypes, mtx.TxOut)
	feeSat := int64(math.Ceil(inputParams.FeeRate * float64(estVSize)))
	changeAmt := totalIn - totalSend - feeSat
	const dust = int64(546)
	changeIdx := -1
	if changeAmt >= dust {
		// 找零
		chgAddr, err := btcutil.DecodeAddress(inputParams.ChangeAddress, network)
		if err != nil {
			return nil, fmt.Errorf("找零地址非法: %v", err)
		}
		chgScript, err := txscript.PayToAddrScript(chgAddr)
		if err != nil {
			return nil, fmt.Errorf("构造找零脚本失败: %v", err)
		}
		mtx.AddTxOut(&wire.TxOut{Value: changeAmt, PkScript: chgScript})
		changeIdx = len(mtx.TxOut) - 1
	}

	// 4) 构建你自定义包（便于回写签名/最终化）
	my := NewV0FromUnsignedTx(mtx) // v0 packet，按 UnsignedTx 初始化 I/O :contentReference[oaicite:3]{index=3}

	// 写入每个输入的 UTXO 元数据（优先 witnessUtxo；legacy 才用 nonWitnessUtxo）
	for i, u := range selectedUTXOs {
		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		typ := detectType(spk)
		switch typ {
		case "p2wpkh", "p2wsh", "p2tr":
			my.SetInputUtxo(i, &wire.TxOut{Value: u.ValueSat, PkScript: spk}, nil) // witnessUtxo 即可 :contentReference[oaicite:4]{index=4}
		default:
			// legacy 需要 NonWitnessUtxo
			if u.NonWitnessTxHex == "" {
				return nil, fmt.Errorf("utxo %s:%d 为 legacy，缺失 nonWitnessTx", u.TxID, u.Vout)
			}
			prevRaw, err := hex.DecodeString(u.NonWitnessTxHex)
			if err != nil {
				return nil, fmt.Errorf("解析 nonWitnessTxHex 失败: %v", err)
			}
			var prev wire.MsgTx
			if err := prev.Deserialize(bytes.NewReader(prevRaw)); err != nil {
				return nil, fmt.Errorf("反序列化 nonWitnessTx 失败: %v", err)
			}
			my.SetInputUtxo(i, nil, &prev)
		}
	}

	// 5) 构建标准 BIP174 PSBT（给 OKX）
	bpkt, err := btcdpsbt.NewFromUnsignedTx(mtx)
	if err != nil {
		return nil, fmt.Errorf("NewFromUnsignedTx 失败: %v", err)
	}
	for i, u := range selectedUTXOs {
		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		typ := detectType(spk)
		switch typ {
		case "p2wpkh", "p2wsh", "p2tr":
			bpkt.Inputs[i].WitnessUtxo = &wire.TxOut{Value: u.ValueSat, PkScript: spk}
		default:
			if u.NonWitnessTxHex == "" {
				return nil, fmt.Errorf("utxo %s:%d 为 legacy，缺失 nonWitnessTx", u.TxID, u.Vout)
			}
			prevRaw, _ := hex.DecodeString(u.NonWitnessTxHex)
			var prev wire.MsgTx
			if err := prev.Deserialize(bytes.NewReader(prevRaw)); err != nil {
				return nil, fmt.Errorf("反序列化 nonWitnessTx 失败: %v", err)
			}
			bpkt.Inputs[i].NonWitnessUtxo = &prev
		}
		// （可选）Sighash/RBF/派生信息等，根据需要继续填
	}

	psbtB64, err := bpkt.B64Encode()
	if err != nil {
		return nil, fmt.Errorf("PSBT B64Encode 序列化失败 失败: %w", err)
	}

	// 调试：unsigned tx hex
	var raw bytes.Buffer
	if err := mtx.Serialize(&raw); err != nil {
		return nil, fmt.Errorf("序列化交易失败: %v", err)
	}

	return &BuildResult{
		PSBTBase64:      psbtB64,
		UnsignedTxHex:   hex.EncodeToString(raw.Bytes()),
		Packet:          my,
		EstimatedVSize:  estVSize,
		FeeSat:          feeSat,
		ChangeOutputIdx: changeIdx,
	}, nil
}

// detectType 仅用于估算 fee（不参与签名逻辑）
func detectType(spk []byte) string {
	n := len(spk)
	if n == 22 && spk[0] == 0x00 && spk[1] == 0x14 {
		return "p2wpkh"
	}
	if n == 34 && spk[0] == 0x00 && spk[1] == 0x20 {
		return "p2wsh"
	}
	if n == 34 && spk[0] == 0x51 && spk[1] == 0x20 {
		return "p2tr"
	}
	if n == 25 && spk[0] == 0x76 && spk[1] == 0xa9 && spk[23] == 0x88 && spk[24] == 0xac {
		return "p2pkh"
	}
	if n == 23 && spk[0] == 0xa9 && spk[1] == 0x14 && spk[22] == 0x87 {
		return "p2sh"
	}
	return "unknown"
}

// 估算 vsize（简单模型）
func estimateVSize(inTypes []string, outs []*wire.TxOut) int {
	vin := 0
	for _, t := range inTypes {
		switch t {
		case "p2wpkh":
			vin += 68
		case "p2tr":
			vin += 58
		case "p2sh": // 假定 p2sh-p2wpkh
			vin += 91
		case "p2pkh":
			vin += 148
		default:
			vin += 110
		}
	}
	vout := 0
	for _, o := range outs {
		// 8(amount) + 1(len) + len(script)
		vout += 9 + len(o.PkScript)
	}
	// 10 base + inputs(41 each no-wit part) 简化：直接返回 vin + vout + 常数
	return vin + vout + 10
}
