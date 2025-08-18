package psbt

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// —— Updater ——

// SetInputUtxo 为第 i 个输入设置 UTXO（优先 witnessUtxo；若为 legacy 建议提供 nonWitnessUtxo）
func (p *Packet) SetInputUtxo(i int, witness *wire.TxOut, nonWitness *wire.MsgTx) {
	in := p.MustInput(i)
	in.WitnessUtxo = witness
	in.NonWitnessUtxo = nonWitness
}

// SetInputScripts 设置 redeem/witness 脚本
func (p *Packet) SetInputScripts(i int, redeem, witness []byte) {
	in := p.MustInput(i)
	in.RedeemScript = append([]byte(nil), redeem...)
	in.WitnessScript = append([]byte(nil), witness...)
}

// SetInputTapScriptPath 设置 Taproot 脚本路径花费所需数据
// script: tapscript 脚本
// control: 控制块
// annex: 可选 annex（若非空，首字节应为 0x50），可传 nil
// stack: 额外入栈元素（不含脚本与控制块，不含 annex/sig），可为空
func (p *Packet) SetInputTapScriptPath(i int, script, control, annex []byte, stack ...[]byte) {
	in := p.MustInput(i)
	in.TapLeafScript = append([]byte(nil), script...)
	in.TapControlBlock = append([]byte(nil), control...)
	in.TapAnnex = append([]byte(nil), annex...)
	if len(stack) > 0 {
		in.TapScriptStack = make([][]byte, 0, len(stack))
		for _, v := range stack {
			in.TapScriptStack = append(in.TapScriptStack, append([]byte(nil), v...))
		}
	}
}

// SetInputDerivation 设置单钥的 BIP32 派生信息
func (p *Packet) SetInputDerivation(i int, d BIP32Derivation) {
	in := p.MustInput(i)
	in.BIP32 = append(in.BIP32, d)
}

// SetV2InputMeta 设置 v2 输入的前序 outpoint/sequence
func (p *Packet) SetV2InputMeta(i int, prev chainhash.Hash, vout uint32, seq uint32) {
	in := p.MustInput(i)
	in.PrevTxID = prev
	in.PrevIndex = vout
	in.Sequence = seq
}

// SetV2OutputMeta 设置 v2 输出金额与脚本
func (p *Packet) SetV2OutputMeta(i int, value int64, spk []byte) {
	out := p.MustOutput(i)
	out.Value = value
	out.ScriptPubKey = append([]byte(nil), spk...)
}

// —— Combiner ——

// Combine 合并多份 PSBT（只合并相同交易的输入签名/元数据；冲突时保留已有并忽略冲突项）
func (p *Packet) Combine(others ...*Packet) error {
	for _, q := range others {
		if q == nil {
			continue
		}
		if err := ensureSameTemplate(p, q); err != nil {
			return err
		}
		for i := range p.Inputs {
			mergeInput(p.Inputs[i], q.Inputs[i])
		}
		for i := range p.Outputs {
			mergeOutput(p.Outputs[i], q.Outputs[i])
		}
	}
	return nil
}

func ensureSameTemplate(a, b *Packet) error {
	if a.Version != b.Version {
		return fmt.Errorf("psbt: version mismatch: %d vs %d", a.Version, b.Version)
	}
	if len(a.Inputs) != len(b.Inputs) || len(a.Outputs) != len(b.Outputs) {
		return errors.New("psbt: io count mismatch")
	}
	if a.IsV0() {
		if a.UnsignedTx == nil || b.UnsignedTx == nil {
			return errors.New("psbt: missing unsigned tx in v0")
		}
		if a.UnsignedTx.TxHash() != b.UnsignedTx.TxHash() {
			return errors.New("psbt: unsigned tx mismatch")
		}
	} else {
		if a.TxVersion != b.TxVersion || a.LockTime != b.LockTime {
			return errors.New("psbt: v2 meta mismatch")
		}
		for i := range a.Inputs {
			if a.Inputs[i].PrevIndex != b.Inputs[i].PrevIndex || a.Inputs[i].PrevTxID != b.Inputs[i].PrevTxID {
				return errors.New("psbt: v2 input meta mismatch")
			}
		}
		for i := range a.Outputs {
			if a.Outputs[i].Value != b.Outputs[i].Value || !bytes.Equal(a.Outputs[i].ScriptPubKey, b.Outputs[i].ScriptPubKey) {
				return errors.New("psbt: v2 output meta mismatch")
			}
		}
	}
	return nil
}

func mergeInput(dst, src *Input) {
	if dst == nil || src == nil {
		return
	}
	if dst.WitnessUtxo == nil && src.WitnessUtxo != nil {
		dst.WitnessUtxo = src.WitnessUtxo
	}
	if dst.NonWitnessUtxo == nil && src.NonWitnessUtxo != nil {
		dst.NonWitnessUtxo = src.NonWitnessUtxo
	}
	if len(dst.RedeemScript) == 0 && len(src.RedeemScript) > 0 {
		dst.RedeemScript = append([]byte(nil), src.RedeemScript...)
	}
	if len(dst.WitnessScript) == 0 && len(src.WitnessScript) > 0 {
		dst.WitnessScript = append([]byte(nil), src.WitnessScript...)
	}
	if dst.SighashType == 0 {
		dst.SighashType = src.SighashType
	}
	if dst.PartialSigs == nil {
		dst.PartialSigs = make(map[string][]byte)
	}
	for k, v := range src.PartialSigs {
		if _, exists := dst.PartialSigs[k]; !exists {
			dst.PartialSigs[k] = append([]byte(nil), v...)
		}
	}

	// Taproot 脚本路径字段
	if len(dst.TapLeafScript) == 0 && len(src.TapLeafScript) > 0 {
		dst.TapLeafScript = append([]byte(nil), src.TapLeafScript...)
	}
	if len(dst.TapControlBlock) == 0 && len(src.TapControlBlock) > 0 {
		dst.TapControlBlock = append([]byte(nil), src.TapControlBlock...)
	}
	if len(dst.TapAnnex) == 0 && len(src.TapAnnex) > 0 {
		dst.TapAnnex = append([]byte(nil), src.TapAnnex...)
	}
	if len(dst.TapScriptStack) == 0 && len(src.TapScriptStack) > 0 {
		dst.TapScriptStack = make([][]byte, 0, len(src.TapScriptStack))
		for _, v := range src.TapScriptStack {
			dst.TapScriptStack = append(dst.TapScriptStack, append([]byte(nil), v...))
		}
	}
}

func mergeOutput(dst, src *Output) {
	if dst == nil || src == nil {
		return
	}
	if len(dst.RedeemScript) == 0 && len(src.RedeemScript) > 0 {
		dst.RedeemScript = append([]byte(nil), src.RedeemScript...)
	}
	if len(dst.WitnessScript) == 0 && len(src.WitnessScript) > 0 {
		dst.WitnessScript = append([]byte(nil), src.WitnessScript...)
	}
	if len(dst.BIP32) == 0 && len(src.BIP32) > 0 {
		dst.BIP32 = append([]BIP32Derivation(nil), src.BIP32...)
	}
}

// —— Signer ——

// SignInput 依据 PSBT 中的 UTXO 信息尝试为索引 i 的输入签名。
// 仅当数据充分时才签名；legacy 需校验 non-witness utxo 的 txid；缺失则返回错误。
// privSign 为回调：接收 digest（双 SHA256）并返回 DER(sig)+hashtype（或 schnorr+hashtype）。
func (p *Packet) SignInput(i int, pubkey33 []byte, sighash txscript.SigHashType, privSign func(digest []byte) ([]byte, error)) error {
	in := p.MustInput(i)
	// 基本约束
	if in.WitnessUtxo == nil && in.NonWitnessUtxo == nil {
		return errors.New("psbt: missing utxo for signing")
	}

	// 构造 sighash
	var (
		msgTx    *wire.MsgTx
		pkScript []byte
		value    int64
		err      error
	)
	if p.IsV0() {
		if p.UnsignedTx == nil {
			return errors.New("psbt: v0 missing unsigned tx")
		}
		msgTx = p.UnsignedTx
	} else {
		// v2 需要临时拼装 MsgTx 计算哈希（不含签名脚本）
		msgTx, err = p.buildMsgTxSkeleton()
		if err != nil {
			return err
		}
	}
	// 选择脚本与金额
	if in.WitnessUtxo != nil {
		pkScript = in.WitnessUtxo.PkScript
		value = in.WitnessUtxo.Value
	} else {
		// legacy：校验 txid 匹配
		prev := msgTx.TxIn[i].PreviousOutPoint
		if in.NonWitnessUtxo == nil {
			return errors.New("psbt: missing non-witness utxo for legacy input")
		}
		if in.NonWitnessUtxo.TxHash() != prev.Hash {
			return errors.New("psbt: non-witness utxo txid mismatch")
		}
		if int(prev.Index) >= len(in.NonWitnessUtxo.TxOut) {
			return errors.New("psbt: non-witness utxo vout out of range")
		}
		pkScript = in.NonWitnessUtxo.TxOut[prev.Index].PkScript
		value = in.NonWitnessUtxo.TxOut[prev.Index].Value
	}

	// 支持：P2WPKH / P2WSH(v0) / P2TR-KeyPath
	isP2WPKH := len(pkScript) == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14
	isP2WSH := len(pkScript) == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20
	isP2TR := len(pkScript) == 34 && pkScript[0] == 0x51 && pkScript[1] == 0x20

	prevFetcher := txscript.NewCannedPrevOutputFetcher(pkScript, value)
	sighashes := txscript.NewTxSigHashes(msgTx, prevFetcher)

	var digest []byte
	if isP2WPKH {
		// scriptCode = OP_DUP OP_HASH160 PUSH20 <20> OP_EQUALVERIFY OP_CHECKSIG
		scriptCode, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
			AddData(pkScript[2:]).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).Script()
		var derr error
		digest, derr = txscript.CalcWitnessSigHash(scriptCode, sighashes, sighash, msgTx, i, value)
		if derr != nil {
			return derr
		}
	} else if isP2WSH {
		if len(in.WitnessScript) == 0 {
			return errors.New("psbt: missing witnessScript for p2wsh")
		}
		var derr error
		digest, derr = txscript.CalcWitnessSigHash(in.WitnessScript, sighashes, sighash, msgTx, i, value)
		if derr != nil {
			return derr
		}
	} else if isP2TR {
		// Taproot: 若提供了脚本路径信息（TapLeafScript/ControlBlock），此处仍生成 keypath 的 digest，
		// 由调用方决定使用何种签名与最终化方式。后续可扩展脚本路径的 sighash 计算。
		var derr error
		digest, derr = txscript.CalcTaprootSignatureHash(sighashes, sighash, msgTx, i, prevFetcher)
		if derr != nil {
			return derr
		}
	} else {
		return errors.New("psbt: unsupported script type for signing")
	}

	sig, err := privSign(digest)
	if err != nil {
		return err
	}
	if in.PartialSigs == nil {
		in.PartialSigs = make(map[string][]byte)
	}
	keyHex := fmt.Sprintf("%x", pubkey33)
	in.PartialSigs[keyHex] = append([]byte(nil), sig...)
	return nil
}

// buildMsgTxSkeleton 把 v2 Packet 拼装为 MsgTx（不含签名脚本/见证），便于 sighash 计算
func (p *Packet) buildMsgTxSkeleton() (*wire.MsgTx, error) {
	if !p.IsV2() {
		return nil, errors.New("psbt: build skeleton only for v2")
	}
	m := &wire.MsgTx{Version: p.TxVersion, LockTime: p.LockTime}
	m.TxIn = make([]*wire.TxIn, len(p.Inputs))
	for i, in := range p.Inputs {
		m.TxIn[i] = &wire.TxIn{PreviousOutPoint: wire.OutPoint{Hash: in.PrevTxID, Index: in.PrevIndex}, Sequence: in.Sequence}
	}
	m.TxOut = make([]*wire.TxOut, len(p.Outputs))
	for i, out := range p.Outputs {
		m.TxOut[i] = &wire.TxOut{Value: out.Value, PkScript: append([]byte(nil), out.ScriptPubKey...)}
	}
	return m, nil
}
