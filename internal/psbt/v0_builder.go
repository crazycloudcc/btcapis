package psbt

import (
	"bytes"
	"encoding/base64"
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
	"github.com/crazycloudcc/btcapis/types"
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

	// 1) 构造收款输出 + 可选 OP_RETURN
	var outputs []*wire.TxOut
	var totalSend int64
	for i, to := range inputParams.ToAddress {
		addr, err := btcutil.DecodeAddress(to, network)
		if err != nil {
			return nil, fmt.Errorf("解析收款地址失败: %v", err)
		}
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, fmt.Errorf("构造收款脚本失败: %v", err)
		}
		amt := int64(math.Round(inputParams.AmountBTC[i] * 1e8))
		if amt <= 0 {
			return nil, fmt.Errorf("非法金额: %f", inputParams.AmountBTC[i])
		}
		outputs = append(outputs, &wire.TxOut{Value: amt, PkScript: pkScript})
		totalSend += amt
	}
	if inputParams.Data != "" {
		script, _ := txscript.NewScriptBuilder().
			AddOp(txscript.OP_RETURN).
			AddData([]byte(inputParams.Data)).
			Script()
		outputs = append(outputs, &wire.TxOut{Value: 0, PkScript: script})
	}

	// 2) 组装未签名交易（v2 tx）
	mtx := wire.NewMsgTx(2)

	// 输入（顺序与 selectedUTXOs 对应）
	var inTypes []string
	var totalIn int64
	seq := uint32(0xFFFFFFFF)
	if inputParams.Replaceable {
		seq = 0xFFFFFFFD // 明确 RBF
	}
	for _, u := range selectedUTXOs {
		h, err := chainhash.NewHashFromStr(u.TxID)
		if err != nil {
			return nil, fmt.Errorf("TxID 解析失败: %v", err)
		}
		mtx.AddTxIn(&wire.TxIn{
			PreviousOutPoint: wire.OutPoint{Hash: *h, Index: u.Vout},
			Sequence:         seq,
		})
		totalIn += u.ValueSat

		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		inTypes = append(inTypes, detectType(spk))

		fmt.Printf("debug info ============= spk: %s\n", u.ScriptPubKeyHex)
		fmt.Printf("debug info ============= NonWitnessTxHex: %s\n", u.NonWitnessTxHex)
		fmt.Printf("debug info ============= RedeemScriptHex: %s\n", u.RedeemScriptHex)
		fmt.Printf("debug info ============= WitnessScriptHex: %s\n", u.WitnessScriptHex)
	}

	// 输出（先放收款/OP_RETURN，找零稍后）
	for _, o := range outputs {
		mtx.AddTxOut(o)
	}
	// locktime（如不使用，传0即可）
	mtx.LockTime = uint32(inputParams.Locktime)

	// 3) 估算费用并添加找零
	estVSize := estimateVSize(inTypes, mtx.TxOut)
	feeSat := int64(math.Ceil(inputParams.FeeRate * float64(estVSize)))
	change := totalIn - totalSend - feeSat
	const dust = int64(546)
	changeIdx := -1
	if change >= dust {
		chgAddr, err := btcutil.DecodeAddress(inputParams.ChangeAddress, network)
		if err != nil {
			return nil, fmt.Errorf("找零地址非法: %v", err)
		}
		chgScript, err := txscript.PayToAddrScript(chgAddr)
		if err != nil {
			return nil, fmt.Errorf("构造找零脚本失败: %v", err)
		}
		mtx.AddTxOut(&wire.TxOut{Value: change, PkScript: chgScript})
		changeIdx = len(mtx.TxOut) - 1
	}
	// （生产中建议在签名后再次精确估算 fee 并微调找零）

	// 4) 构建你自定义包（便于回写签名/最终化）
	my := NewV0FromUnsignedTx(mtx) // v0 packet，按 UnsignedTx 初始化 I/O :contentReference[oaicite:3]{index=3}

	// 5) 标准 BIP174 PSBT（给 OKX）
	bpkt, err := btcdpsbt.NewFromUnsignedTx(mtx)
	if err != nil {
		return nil, fmt.Errorf("NewFromUnsignedTx 失败: %w", err)
	}

	// 6) 为每个输入写入 UTXO 元数据（witness 优先；legacy 需 nonWitness）
	for i, u := range selectedUTXOs {
		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		typ := detectType(spk)

		var redeem []byte
		if u.RedeemScriptHex != "" {
			if rb, err := hex.DecodeString(u.RedeemScriptHex); err == nil {
				redeem = rb
			}
		}

		var wscript []byte
		if u.WitnessScriptHex != "" {
			if ws, err := hex.DecodeString(u.WitnessScriptHex); err == nil {
				wscript = ws
			}
		}

		switch typ {
		// 原生 segwit / taproot：WitnessUtxo 即可
		case "p2wpkh", "p2wsh", "p2tr":
			txout := &wire.TxOut{Value: u.ValueSat, PkScript: spk}
			bpkt.Inputs[i].WitnessUtxo = txout
			// 若是 p2wsh 还应提供 witnessScript（如是多签/脚本花费）
			if typ == "p2wsh" && len(wscript) > 0 {
				bpkt.Inputs[i].WitnessScript = wscript
			}
			// 你的包也同步（可选）
			my.SetInputUtxo(i, txout, nil) // witness 优先 :contentReference[oaicite:5]{index=5}

		case "p2sh":
			// P2SH 分两种：嵌套 segwit（需 redeemScript + witnessUtxo），或传统 P2SH（需 nonWitnessUtxo）
			if len(redeem) > 0 && isSegwitProgram(redeem) {
				// 视为 P2SH-P2WPKH / P2SH-P2WSH
				txout := &wire.TxOut{Value: u.ValueSat, PkScript: spk}
				bpkt.Inputs[i].WitnessUtxo = txout
				bpkt.Inputs[i].RedeemScript = redeem
				if len(wscript) > 0 { // P2SH-P2WSH 的脚本体
					bpkt.Inputs[i].WitnessScript = wscript
				}
				my.SetInputUtxo(i, txout, nil)
				// 如你的包支持脚本字段，也可同步设置（可选）

			} else {
				// 传统 P2SH = legacy，必须提供 NonWitnessUtxo
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
				bpkt.Inputs[i].NonWitnessUtxo = &prev
				my.SetInputUtxo(i, nil, &prev)
			}

		case "p2pkh":
			// 传统 P2PKH：必须 NonWitnessUtxo
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
			bpkt.Inputs[i].NonWitnessUtxo = &prev
			my.SetInputUtxo(i, nil, &prev)

		default:
			// 未知脚本，保守处理：若提供了 NonWitnessTx 用它，否则尝试 WitnessUtxo
			if u.NonWitnessTxHex != "" {
				prevRaw, _ := hex.DecodeString(u.NonWitnessTxHex)
				var prev wire.MsgTx
				_ = prev.Deserialize(bytes.NewReader(prevRaw))
				bpkt.Inputs[i].NonWitnessUtxo = &prev
				my.SetInputUtxo(i, nil, &prev)
			} else {
				txout := &wire.TxOut{Value: u.ValueSat, PkScript: spk}
				bpkt.Inputs[i].WitnessUtxo = txout
				my.SetInputUtxo(i, txout, nil)
			}
		}
	}

	// 7) 编码为 Base64（OKX 直接可用）
	psbtB64, err := bpkt.B64Encode()
	if err != nil {
		// 某些旧版本缺少 B64Encode 方法，可退化为 Serialize + 手动 base64
		var buf bytes.Buffer
		if err2 := bpkt.Serialize(&buf); err2 != nil {
			return nil, fmt.Errorf("PSBT 编码失败: %v / %v", err, err2)
		}
		psbtB64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	}

	// 8) 输出原始未签名交易（调试）
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

// -------------------- 辅助函数 --------------------

// detectType：粗略识别常见输出脚本类型（仅用于估算/分支）
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
	if n == 25 && spk[0] == 0x76 && spk[1] == 0xa9 && spk[2] == 0x14 && spk[23] == 0x88 && spk[24] == 0xac {
		return "p2pkh"
	}
	if n == 23 && spk[0] == 0xa9 && spk[1] == 0x14 && spk[22] == 0x87 {
		return "p2sh"
	}
	return "unknown"
}

// 估算 vsize（简化模型，足以先行算费）
func estimateVSize(inTypes []string, outs []*wire.TxOut) int {
	vin := 0
	for _, t := range inTypes {
		switch t {
		case "p2wpkh":
			vin += 68
		case "p2tr":
			vin += 58
		case "p2wsh":
			vin += 104 // 视脚本大小而变，这里粗估
		case "p2pkh":
			vin += 148
		case "p2sh": // 可能是 p2sh-p2wpkh（~91）或传统（~148），取中位保守
			vin += 120
		default:
			vin += 110
		}
	}
	vout := 0
	for _, o := range outs {
		vout += 9 + len(o.PkScript) // 8(amount)+1(len)+script
	}
	// 粗略常数开销
	return vin + vout + 10
}

// 判定传统 P2SH：OP_HASH160 0x14 <20> OP_EQUAL
func isP2SH(spk []byte) bool {
	return len(spk) == 23 && spk[0] == 0xa9 && spk[1] == 0x14 && spk[22] == 0x87
}

// 判定传统 P2PKH：OP_DUP OP_HASH160 0x14 <20> OP_EQUALVERIFY OP_CHECKSIG
func isP2PKH(spk []byte) bool {
	return len(spk) == 25 && spk[0] == 0x76 && spk[1] == 0xa9 && spk[2] == 0x14 && spk[23] == 0x88 && spk[24] == 0xac
}

// legacy = P2PKH 或 传统 P2SH（是否嵌套由 redeemScript 再判定）
func isLegacy(spk []byte) bool { return isP2PKH(spk) || isP2SH(spk) }

// 判定“SegWit 程序”本体（用于 P2SH redeemScript：P2WPKH/P2WSH/Taproot）
func isSegwitProgram(script []byte) bool {
	n := len(script)
	if n < 4 || n > 42 {
		return false
	}
	ver := script[0]
	push := int(script[1])
	if ver != 0x00 && (ver < 0x51 || ver > 0x60) {
		return false
	}
	if 2+push != n {
		return false
	}
	// 仅接受 20/32（v0-pkh / v0-wsh / v1-taproot）的典型长度
	return push == 20 || push == 32
}
