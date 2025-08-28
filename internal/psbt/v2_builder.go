package psbt

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// CreatePSBTv2ForOKX 构建 PSBT v2（BIP-370），并返回 v2 base64 字符串。
func CreatePSBTv2ForOKX(
	inputParams *types.TxInputParams,
	selectedUTXOs []PsbtUTXO,
	network *chaincfg.Params,
) (*BuildResult, error) {
	if len(inputParams.ToAddress) == 0 || len(inputParams.ToAddress) != len(inputParams.AmountBTC) {
		return nil, fmt.Errorf("to_address 与 amount 数量需一致且 > 0")
	}
	if len(selectedUTXOs) == 0 {
		return nil, fmt.Errorf("未提供可用 UTXO")
	}
	if inputParams.ChangeAddress == "" {
		return nil, fmt.Errorf("必须提供找零地址")
	}

	// 1) 构造输出
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
		script, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData([]byte(inputParams.Data)).Script()
		outputs = append(outputs, &wire.TxOut{Value: 0, PkScript: script})
	}

	// 2) 临时 MsgTx 用于费率估算
	var inTypes []string
	var totalIn int64
	seq := uint32(0xFFFFFFFF)
	if inputParams.Replaceable {
		seq = 0xFFFFFFFD
	}
	tmp := wire.NewMsgTx(2)
	for _, u := range selectedUTXOs {
		h, err := chainhash.NewHashFromStr(u.TxID)
		if err != nil {
			return nil, fmt.Errorf("TxID 解析失败: %v", err)
		}
		tmp.AddTxIn(&wire.TxIn{PreviousOutPoint: wire.OutPoint{Hash: *h, Index: u.Vout}, Sequence: seq})
		totalIn += u.ValueSat
		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		inTypes = append(inTypes, detectType(spk))
	}
	for _, o := range outputs {
		tmp.AddTxOut(o)
	}
	tmp.LockTime = uint32(inputParams.Locktime)

	estVSize := estimateVSize(inTypes, tmp.TxOut)
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
		tmp.AddTxOut(&wire.TxOut{Value: change, PkScript: chgScript})
		changeIdx = len(tmp.TxOut) - 1
	}

	// 3) 构建 v2 Packet
	pkt := NewV2(2, uint32(inputParams.Locktime), len(selectedUTXOs), len(tmp.TxOut))
	for i, u := range selectedUTXOs {
		h, _ := chainhash.NewHashFromStr(u.TxID)
		pkt.SetV2InputMeta(i, *h, u.Vout, seq)
		spk, _ := hex.DecodeString(u.ScriptPubKeyHex)
		typ := detectType(spk)
		if typ == "p2wpkh" || typ == "p2wsh" || typ == "p2tr" {
			pkt.SetInputUtxo(i, &wire.TxOut{Value: u.ValueSat, PkScript: spk}, nil)
			if typ == "p2wsh" && u.WitnessScriptHex != "" {
				if ws, err := hex.DecodeString(u.WitnessScriptHex); err == nil {
					pkt.SetInputScripts(i, nil, ws)
				}
			}
		} else {
			// 其他类型：必须 NonWitnessTx
			if u.NonWitnessTxHex == "" {
				return nil, fmt.Errorf("utxo %s:%d 缺失 NonWitnessTx", u.TxID, u.Vout)
			}
			prevRaw, err := hex.DecodeString(u.NonWitnessTxHex)
			if err != nil {
				return nil, fmt.Errorf("解析 NonWitnessTxHex 失败: %v", err)
			}
			var prev wire.MsgTx
			if err := prev.Deserialize(bytes.NewReader(prevRaw)); err != nil {
				return nil, fmt.Errorf("反序列化 NonWitnessTx 失败: %v", err)
			}
			pkt.SetInputUtxo(i, nil, &prev)
		}
	}
	for i, o := range tmp.TxOut {
		pkt.SetV2OutputMeta(i, o.Value, o.PkScript)
	}

	// 4) v2 序列化
	raw, err := SerializeV2Packet(pkt)
	if err != nil {
		return nil, fmt.Errorf("v2 序列化失败: %v", err)
	}
	psbtB64 := base64.StdEncoding.EncodeToString(raw)

	return &BuildResult{PSBTBase64: psbtB64, UnsignedTxHex: unsignedBtcTxHex(tmp), Packet: pkt, EstimatedVSize: estVSize, FeeSat: feeSat, ChangeOutputIdx: changeIdx}, nil
}

// ====== PSBT v2 序列化 ======
const (
	psbtMagic1 = 0x70
	psbtMagic2 = 0x73
	psbtMagic3 = 0x62
	psbtMagic4 = 0x74
	psbtSep    = 0xff

	PSBT_GLOBAL_TX_VERSION    = 0x02
	PSBT_GLOBAL_FALLBACK_LOCK = 0x03
	PSBT_GLOBAL_INPUT_COUNT   = 0x04
	PSBT_GLOBAL_OUTPUT_COUNT  = 0x05
	PSBT_GLOBAL_VERSION       = 0xfb

	PSBT_IN_NON_WITNESS_UTXO = 0x00
	PSBT_IN_WITNESS_UTXO     = 0x01
	PSBT_IN_REDEEM_SCRIPT    = 0x04
	PSBT_IN_WITNESS_SCRIPT   = 0x05
	PSBT_IN_PREVIOUS_TXID    = 0x0e
	PSBT_IN_OUTPUT_INDEX     = 0x0f
	PSBT_IN_SEQUENCE         = 0x10

	PSBT_OUT_REDEEM_SCRIPT  = 0x00
	PSBT_OUT_WITNESS_SCRIPT = 0x01
	PSBT_OUT_AMOUNT         = 0x03
	PSBT_OUT_SCRIPT         = 0x04
)

func SerializeV2Packet(p *Packet) ([]byte, error) {
	if p == nil || !p.IsV2() {
		return nil, fmt.Errorf("psbt: SerializeV2Packet 仅支持 v2")
	}
	if len(p.Inputs) == 0 || len(p.Outputs) == 0 {
		return nil, fmt.Errorf("psbt: v2 至少1个输入和输出")
	}

	var b bytes.Buffer
	// magic
	b.Write([]byte{'p', 's', 'b', 't'})
	b.WriteByte(psbtSep)

	// global
	writeKVU32(&b, PSBT_GLOBAL_VERSION, 2) // 写死 2，避免使用 VersionV2 未定义
	writeKVI32(&b, PSBT_GLOBAL_TX_VERSION, p.TxVersion)
	if p.LockTime != 0 {
		writeKVU32(&b, PSBT_GLOBAL_FALLBACK_LOCK, p.LockTime)
	}
	writeKVU32(&b, PSBT_GLOBAL_INPUT_COUNT, uint32(len(p.Inputs)))
	writeKVU32(&b, PSBT_GLOBAL_OUTPUT_COUNT, uint32(len(p.Outputs)))
	b.WriteByte(0x00)

	for i, in := range p.Inputs {
		// PSBTv2 要求 prev txid 采用与交易序列化一致的字节序（LE）
		{
			h := in.PrevTxID
			le := h
			// 显式反转为 LE（btcd 的 Hash 为 BE 表示）
			for i, j := 0, len(le)-1; i < j; i, j = i+1, j-1 {
				le[i], le[j] = le[j], le[i]
			}
			writeKV(&b, []byte{PSBT_IN_PREVIOUS_TXID}, le[:])
		}
		writeKVU32(&b, PSBT_IN_OUTPUT_INDEX, in.PrevIndex)
		writeKVU32(&b, PSBT_IN_SEQUENCE, in.Sequence)
		if in.WitnessUtxo != nil {
			writeKV(&b, []byte{PSBT_IN_WITNESS_UTXO}, serializeTxOut(in.WitnessUtxo))
		} else if in.NonWitnessUtxo != nil {
			var nb bytes.Buffer
			_ = in.NonWitnessUtxo.Serialize(&nb)
			writeKV(&b, []byte{PSBT_IN_NON_WITNESS_UTXO}, nb.Bytes())
		}
		if len(in.RedeemScript) > 0 {
			writeKV(&b, []byte{PSBT_IN_REDEEM_SCRIPT}, in.RedeemScript)
		}
		if len(in.WitnessScript) > 0 {
			writeKV(&b, []byte{PSBT_IN_WITNESS_SCRIPT}, in.WitnessScript)
		}
		b.WriteByte(0x00)
		_ = i
	}

	for i, out := range p.Outputs {
		writeKVI64(&b, PSBT_OUT_AMOUNT, out.Value)
		writeKV(&b, []byte{PSBT_OUT_SCRIPT}, out.ScriptPubKey)
		if len(out.RedeemScript) > 0 {
			writeKV(&b, []byte{PSBT_OUT_REDEEM_SCRIPT}, out.RedeemScript)
		}
		if len(out.WitnessScript) > 0 {
			writeKV(&b, []byte{PSBT_OUT_WITNESS_SCRIPT}, out.WitnessScript)
		}
		b.WriteByte(0x00)
		_ = i
	}

	return b.Bytes(), nil
}

func writeKV(buf *bytes.Buffer, key []byte, val []byte) {
	writeVarInt(buf, uint64(len(key)))
	buf.Write(key)
	writeVarInt(buf, uint64(len(val)))
	buf.Write(val)
}

func writeKVU32(buf *bytes.Buffer, keyType byte, v uint32) {
	var tmp [4]byte
	tmp[0] = byte(v)
	tmp[1] = byte(v >> 8)
	tmp[2] = byte(v >> 16)
	tmp[3] = byte(v >> 24)
	writeKV(buf, []byte{keyType}, tmp[:])
}

func writeKVI32(buf *bytes.Buffer, keyType byte, v int32) { writeKVU32(buf, keyType, uint32(v)) }

func writeKVI64(buf *bytes.Buffer, keyType byte, v int64) {
	var tmp [8]byte
	uv := uint64(v)
	for i := 0; i < 8; i++ {
		tmp[i] = byte(uv >> (8 * i))
	}
	writeKV(buf, []byte{keyType}, tmp[:])
}

func writeKVVarInt(buf *bytes.Buffer, keyType byte, v uint64) {
	vb := new(bytes.Buffer)
	writeVarInt(vb, v)
	writeKV(buf, []byte{keyType}, vb.Bytes())
}

func writeVarInt(buf *bytes.Buffer, v uint64) {
	switch {
	case v < 0xfd:
		buf.WriteByte(byte(v))
	case v <= 0xffff:
		buf.WriteByte(0xfd)
		buf.WriteByte(byte(v))
		buf.WriteByte(byte(v >> 8))
	case v <= 0xffffffff:
		buf.WriteByte(0xfe)
		for i := 0; i < 4; i++ {
			buf.WriteByte(byte(v >> (8 * i)))
		}
	default:
		buf.WriteByte(0xff)
		for i := 0; i < 8; i++ {
			buf.WriteByte(byte(v >> (8 * i)))
		}
	}
}

func serializeTxOut(o *wire.TxOut) []byte {
	var b bytes.Buffer
	var tmp [8]byte
	uv := uint64(o.Value)
	for i := 0; i < 8; i++ {
		tmp[i] = byte(uv >> (8 * i))
	}
	b.Write(tmp[:])
	writeVarInt(&b, uint64(len(o.PkScript)))
	b.Write(o.PkScript)
	return b.Bytes()
}

func unsignedBtcTxHex(m *wire.MsgTx) string {
	var buf bytes.Buffer
	_ = m.Serialize(&buf)
	return hex.EncodeToString(buf.Bytes())
}
