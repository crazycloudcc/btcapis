package psbt

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"golang.org/x/crypto/ripemd160"
)

// hash160 计算RIPEMD160(SHA256(data))
func hash160(data []byte) []byte {
	sha256Hash := sha256.Sum256(data)
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256Hash[:])
	return ripemd160Hash.Sum(nil)
}

// FinalizeInput 根据已收集到的 PartialSigs/WitnessScript 等，生成最终的 ScriptSig/ScriptWitness。
// 支持 p2wpkh / p2wsh(v0) / p2tr-keypath / p2sh-wrapped-segwit / legacy p2pkh/p2sh
func (p *Packet) FinalizeInput(i int) error {
	in := p.MustInput(i)
	var pkScript []byte
	var value int64
	if in.WitnessUtxo != nil {
		pkScript = in.WitnessUtxo.PkScript
		value = in.WitnessUtxo.Value
	} else if in.NonWitnessUtxo != nil {
		// 对于legacy输入，从NonWitnessUtxo获取pkScript和value
		if p.IsV0() && p.UnsignedTx != nil {
			prev := p.UnsignedTx.TxIn[i].PreviousOutPoint
			// 允许 segwit 仅带 NonWitnessUtxo，同时严格校验 txid 与 vout
			if in.NonWitnessUtxo.TxHash() != prev.Hash {
				return errors.New("psbt: non-witness utxo txid mismatch")
			}
			if int(prev.Index) >= len(in.NonWitnessUtxo.TxOut) {
				return errors.New("psbt: non-witness utxo vout out of range")
			}
			pkScript = in.NonWitnessUtxo.TxOut[prev.Index].PkScript
			value = in.NonWitnessUtxo.TxOut[prev.Index].Value
		} else {
			return errors.New("psbt: missing utxo to finalize")
		}
	} else {
		return errors.New("psbt: missing utxo to finalize")
	}

	// 识别脚本类型
	scriptType := classifyScript(pkScript)

	// 不在此处统一强制需要 PartialSigs，由各分支自行校验（例如 tapscript 可无签名）。

	switch scriptType {
	case "p2wpkh":
		return p.finalizeP2WPKH(in, pkScript, value)
	case "p2wsh":
		return p.finalizeP2WSH(in, pkScript, value)
	case "p2tr":
		return p.finalizeP2TR(in, pkScript, value)
	case "p2sh":
		return p.finalizeP2SH(in, pkScript, value)
	case "p2pkh":
		return p.finalizeP2PKH(in, pkScript, value)
	default:
		return fmt.Errorf("psbt: unsupported script type to finalize: %s", scriptType)
	}
}

// classifyScript 识别脚本类型
func classifyScript(pkScript []byte) string {
	n := len(pkScript)
	if n == 0 {
		return "unknown"
	}

	// P2PKH: OP_DUP OP_HASH160 PUSH20 <20> OP_EQUALVERIFY OP_CHECKSIG
	if n == 25 && pkScript[0] == 0x76 && pkScript[1] == 0xa9 && pkScript[2] == 0x14 && pkScript[23] == 0x88 && pkScript[24] == 0xac {
		return "p2pkh"
	}
	// P2SH: OP_HASH160 PUSH20 <20> OP_EQUAL
	if n == 23 && pkScript[0] == 0xa9 && pkScript[1] == 0x14 && pkScript[22] == 0x87 {
		return "p2sh"
	}
	// P2WPKH v0: OP_0 PUSH20 <20>
	if n == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14 {
		return "p2wpkh"
	}
	// P2WSH v0: OP_0 PUSH32 <32>
	if n == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20 {
		return "p2wsh"
	}
	// P2TR v1: OP_1 PUSH32 <32>
	if n == 34 && pkScript[0] == 0x51 && pkScript[1] == 0x20 {
		return "p2tr"
	}
	return "unknown"
}

// finalizeP2WPKH 最终化P2WPKH输入
func (p *Packet) finalizeP2WPKH(in *Input, pkScript []byte, value int64) error {
	// 尝试从BIP32获取pubkey，如果没有则从PartialSigs的key反解
	var pubkey []byte
	var sig []byte

	if len(in.BIP32) > 0 {
		for _, d := range in.BIP32 {
			keyHex := fmt.Sprintf("%x", d.PubKey)
			if s, ok := in.PartialSigs[keyHex]; ok {
				pubkey = d.PubKey
				sig = s
				break
			}
		}
	}

	// 如果BIP32中没有找到，尝试从PartialSigs的key反解
	if len(pubkey) == 0 {
		for keyHex, s := range in.PartialSigs {
			// 尝试解析hex格式的pubkey
			if len(keyHex) == 66 { // 33字节压缩公钥的hex长度
				pubkeyBytes, err := hex.DecodeString(keyHex)
				if err == nil && len(pubkeyBytes) == 33 {
					pubkey = pubkeyBytes
					sig = s
					break
				}
			}
		}
	}

	if len(pubkey) == 0 || len(sig) == 0 {
		return errors.New("psbt: matching pubkey/sig not found for p2wpkh")
	}

	// 验证pubkey hash是否匹配
	expectedHash := pkScript[2:] // 跳过OP_0和PUSH20
	actualHash := hash160(pubkey)
	if !bytes.Equal(expectedHash, actualHash) {
		return fmt.Errorf("psbt: pubkey hash mismatch for p2wpkh: expected %x, got %x", expectedHash, actualHash)
	}

	in.FinalScriptWitness = wire.TxWitness{append([]byte(nil), sig...), append([]byte(nil), pubkey...)}
	in.FinalScriptSig = nil
	p.cleanupInput(in)
	return nil
}

// finalizeP2WSH 最终化P2WSH输入
func (p *Packet) finalizeP2WSH(in *Input, pkScript []byte, value int64) error {
	if len(in.WitnessScript) == 0 {
		return errors.New("psbt: missing witnessScript for p2wsh finalize")
	}

	// 验证witnessScript hash是否匹配
	expectedHash := pkScript[2:] // 跳过OP_0和PUSH32
	actualHash := sha256.Sum256(in.WitnessScript)
	if !bytes.Equal(expectedHash, actualHash[:]) {
		return fmt.Errorf("psbt: witnessScript hash mismatch for p2wsh: expected %x, got %x", expectedHash, actualHash[:])
	}

	// 检测多签：是否包含 OP_CHECKMULTISIG/VERIFY (0xae/0xaf)
	ws := in.WitnessScript
	isMultisig := false
	for _, b := range ws {
		if b == 0xae || b == 0xaf {
			isMultisig = true
			break
		}
	}

	// 解析witnessScript中的公钥顺序
	pubkeyOrder, err := p.parseWitnessScriptPubkeys(ws)
	if err != nil {
		return fmt.Errorf("psbt: failed to parse witnessScript pubkeys: %v", err)
	}

	// 按witnessScript中公钥顺序构建签名栈
	stack := make([][]byte, 0, len(in.PartialSigs)+2)
	if isMultisig {
		stack = append(stack, []byte{}) // 占位空元素
	}

	// 按witnessScript中公钥顺序添加签名
	for _, pubkey := range pubkeyOrder {
		keyHex := fmt.Sprintf("%x", pubkey)
		if sig, ok := in.PartialSigs[keyHex]; ok {
			stack = append(stack, append([]byte(nil), sig...))
		}
	}

	stack = append(stack, append([]byte(nil), ws...))
	if len(stack) < 2 {
		return errors.New("psbt: not enough sigs to finalize p2wsh")
	}

	in.FinalScriptWitness = wire.TxWitness(stack)
	in.FinalScriptSig = nil
	p.cleanupInput(in)
	// 可选择清理或保留 WitnessScript；此处保留供上层检查
	return nil
}

// parseWitnessScriptPubkeys 解析witnessScript中的公钥顺序
// 支持常见的多签脚本格式：OP_2 <pubkey1> <pubkey2> <pubkey3> OP_3 OP_CHECKMULTISIG
func (p *Packet) parseWitnessScriptPubkeys(script []byte) ([][]byte, error) {
	var pubkeys [][]byte
	i := 0

	// 跳过开头的OP_2, OP_3等操作码
	for i < len(script) && (script[i] >= 0x51 && script[i] <= 0x60) { // OP_1 to OP_16
		i++
	}

	// 解析公钥
	for i < len(script) {
		if script[i] == 0xae || script[i] == 0xaf { // OP_CHECKMULTISIG or OP_CHECKMULTISIGVERIFY
			break
		}

		// 检查是否是PUSH操作码
		if script[i] == 0x21 { // PUSH33
			if i+1+33 > len(script) {
				return nil, fmt.Errorf("psbt: invalid script: unexpected end after PUSH33")
			}
			pubkey := script[i+1 : i+1+33]
			pubkeys = append(pubkeys, append([]byte(nil), pubkey...))
			i += 1 + 33
		} else if script[i] == 0x41 { // PUSH65 (未压缩公钥)
			if i+1+65 > len(script) {
				return nil, fmt.Errorf("psbt: invalid script: unexpected end after PUSH65")
			}
			pubkey := script[i+1 : i+1+65]
			pubkeys = append(pubkeys, append([]byte(nil), pubkey...))
			i += 1 + 65
		} else {
			// 跳过其他操作码
			i++
		}
	}

	if len(pubkeys) == 0 {
		return nil, fmt.Errorf("psbt: no pubkeys found in witnessScript")
	}

	return pubkeys, nil
}

// finalizeP2TR 最终化P2TR输入
func (p *Packet) finalizeP2TR(in *Input, pkScript []byte, value int64) error {
	// 先尝试脚本路径：需要 TapLeafScript + TapControlBlock
	if len(in.TapLeafScript) > 0 && len(in.TapControlBlock) > 0 {
		stack := make([][]byte, 0, 4+len(in.TapScriptStack)+len(in.PartialSigs))
		// 可选 annex：若存在且首字节为 0x50，则作为 witness[0]
		if len(in.TapAnnex) > 0 {
			if in.TapAnnex[0] != 0x50 {
				return errors.New("psbt: invalid taproot annex (must start with 0x50)")
			}
			stack = append(stack, append([]byte(nil), in.TapAnnex...))
		}
		if len(in.TapScriptStack) > 0 {
			for _, v := range in.TapScriptStack {
				stack = append(stack, append([]byte(nil), v...))
			}
		} else {
			// 若未显式提供栈，则回退为按 BIP32 顺序附加 PartialSigs
			for _, d := range in.BIP32 {
				keyHex := fmt.Sprintf("%x", d.PubKey)
				if s, ok := in.PartialSigs[keyHex]; ok {
					stack = append(stack, append([]byte(nil), s...))
				}
			}
		}
		// 附上 tapscript 与 control block
		stack = append(stack, append([]byte(nil), in.TapLeafScript...))
		stack = append(stack, append([]byte(nil), in.TapControlBlock...))
		if len(stack) < 2 {
			return errors.New("psbt: not enough elements to finalize p2tr script path")
		}
		in.FinalScriptWitness = wire.TxWitness(stack)
		in.FinalScriptSig = nil
		p.cleanupInput(in)
		return nil
	}

	// 否则回退为 keypath：最终 witness = [schnorr_sig]
	// 如果有annex，需要包含
	if len(in.TapAnnex) > 0 {
		if in.TapAnnex[0] != 0x50 {
			return errors.New("psbt: invalid taproot annex (must start with 0x50)")
		}
		for _, s := range in.PartialSigs {
			in.FinalScriptWitness = wire.TxWitness{append([]byte(nil), in.TapAnnex...), append([]byte(nil), s...)}
			in.FinalScriptSig = nil
			p.cleanupInput(in)
			return nil
		}
	} else {
		for _, s := range in.PartialSigs {
			in.FinalScriptWitness = wire.TxWitness{append([]byte(nil), s...)}
			in.FinalScriptSig = nil
			p.cleanupInput(in)
			return nil
		}
	}
	return errors.New("psbt: no schnorr sig for p2tr finalize")
}

// finalizeP2SH 最终化P2SH输入（包括包裹的SegWit）
func (p *Packet) finalizeP2SH(in *Input, pkScript []byte, value int64) error {
	if len(in.RedeemScript) == 0 {
		return errors.New("psbt: missing redeemScript for p2sh finalize")
	}

	// 检查redeemScript类型
	redeemType := classifyScript(in.RedeemScript)

	switch redeemType {
	case "p2wpkh":
		// P2SH-P2WPKH: FinalScriptSig = PUSH(redeemScript), FinalWitness = [sig, pubkey]
		return p.finalizeP2SHWrappedP2WPKH(in, pkScript, value)
	case "p2wsh":
		// P2SH-P2WSH: FinalScriptSig = PUSH(redeemScript), FinalWitness = [..., witnessScript]
		return p.finalizeP2SHWrappedP2WSH(in, pkScript, value)
	default:
		// 普通P2SH（非见证程序）
		return p.finalizeLegacyP2SH(in, pkScript, value)
	}
}

// finalizeP2SHWrappedP2WPKH 最终化P2SH包裹的P2WPKH
func (p *Packet) finalizeP2SHWrappedP2WPKH(in *Input, pkScript []byte, value int64) error {
	// 验证redeemScript hash是否匹配
	expectedHash := pkScript[2:] // 跳过OP_HASH160和PUSH20
	actualHash := hash160(in.RedeemScript)
	if !bytes.Equal(expectedHash, actualHash) {
		return fmt.Errorf("psbt: redeemScript hash mismatch for p2sh-p2wpkh: expected %x, got %x", expectedHash, actualHash)
	}

	// 从redeemScript获取pubkey hash
	redeemPkScript := in.RedeemScript
	if len(redeemPkScript) != 22 || redeemPkScript[0] != 0x00 || redeemPkScript[1] != 0x14 {
		return errors.New("psbt: invalid redeemScript for p2sh-p2wpkh")
	}

	// 查找对应的签名和公钥
	var pubkey []byte
	var sig []byte

	if len(in.BIP32) > 0 {
		for _, d := range in.BIP32 {
			keyHex := fmt.Sprintf("%x", d.PubKey)
			if s, ok := in.PartialSigs[keyHex]; ok {
				// 验证公钥hash是否匹配
				if bytes.Equal(hash160(d.PubKey), redeemPkScript[2:]) {
					pubkey = d.PubKey
					sig = s
					break
				}
			}
		}
	}

	if len(pubkey) == 0 || len(sig) == 0 {
		return errors.New("psbt: matching pubkey/sig not found for p2sh-p2wpkh")
	}

	// FinalScriptSig = PUSH(redeemScript)
	in.FinalScriptSig = p.buildScriptSig([][]byte{append([]byte(nil), in.RedeemScript...)})
	// FinalWitness = [sig, pubkey]
	in.FinalScriptWitness = wire.TxWitness{append([]byte(nil), sig...), append([]byte(nil), pubkey...)}

	p.cleanupInput(in)
	return nil
}

// finalizeP2SHWrappedP2WSH 最终化P2SH包裹的P2WSH
func (p *Packet) finalizeP2SHWrappedP2WSH(in *Input, pkScript []byte, value int64) error {
	// 验证redeemScript hash是否匹配
	expectedHash := pkScript[2:] // 跳过OP_HASH160和PUSH20
	actualHash := hash160(in.RedeemScript)
	if !bytes.Equal(expectedHash, actualHash) {
		return fmt.Errorf("psbt: redeemScript hash mismatch for p2sh-p2wsh: expected %x, got %x", expectedHash, actualHash)
	}

	if len(in.WitnessScript) == 0 {
		return errors.New("psbt: missing witnessScript for p2sh-p2wsh finalize")
	}

	// 验证witnessScript hash是否匹配
	redeemPkScript := in.RedeemScript
	if len(redeemPkScript) != 34 || redeemPkScript[0] != 0x00 || redeemPkScript[1] != 0x20 {
		return errors.New("psbt: invalid redeemScript for p2sh-p2wsh")
	}

	expectedWitnessHash := redeemPkScript[2:] // 跳过OP_0和PUSH32
	actualWitnessHash := sha256.Sum256(in.WitnessScript)
	if !bytes.Equal(expectedWitnessHash, actualWitnessHash[:]) {
		return fmt.Errorf("psbt: witnessScript hash mismatch for p2sh-p2wsh: expected %x, got %x", expectedWitnessHash, actualWitnessHash[:])
	}

	// 构建见证栈
	stack := make([][]byte, 0, len(in.PartialSigs)+2)

	// 检测多签
	ws := in.WitnessScript
	isMultisig := false
	for _, b := range ws {
		if b == 0xae || b == 0xaf {
			isMultisig = true
			break
		}
	}

	if isMultisig {
		stack = append(stack, []byte{}) // 占位空元素
	}

	// 按BIP32顺序添加签名
	for _, d := range in.BIP32 {
		keyHex := fmt.Sprintf("%x", d.PubKey)
		if s, ok := in.PartialSigs[keyHex]; ok {
			stack = append(stack, append([]byte(nil), s...))
		}
	}

	stack = append(stack, append([]byte(nil), ws...))

	// FinalScriptSig = PUSH(redeemScript)
	in.FinalScriptSig = p.buildScriptSig([][]byte{append([]byte(nil), in.RedeemScript...)})
	// FinalWitness = [..., witnessScript]
	in.FinalScriptWitness = wire.TxWitness(stack)

	p.cleanupInput(in)
	return nil
}

// finalizeLegacyP2SH 最终化普通P2SH（非见证程序）
func (p *Packet) finalizeLegacyP2SH(in *Input, pkScript []byte, value int64) error {
	// 验证redeemScript hash是否匹配
	expectedHash := pkScript[2:] // 跳过OP_HASH160和PUSH20
	actualHash := hash160(in.RedeemScript)
	if !bytes.Equal(expectedHash, actualHash) {
		return fmt.Errorf("psbt: redeemScript hash mismatch for legacy p2sh: expected %x, got %x", expectedHash, actualHash)
	}

	// 构建FinalScriptSig: [sig1, sig2, ..., redeemScript]
	stack := make([][]byte, 0, len(in.PartialSigs)+1)

	// 添加签名（按BIP32顺序）
	for _, d := range in.BIP32 {
		keyHex := fmt.Sprintf("%x", d.PubKey)
		if s, ok := in.PartialSigs[keyHex]; ok {
			stack = append(stack, append([]byte(nil), s...))
		}
	}

	// 添加redeemScript
	stack = append(stack, append([]byte(nil), in.RedeemScript...))

	// 构建FinalScriptSig
	in.FinalScriptSig = p.buildScriptSig(stack)
	in.FinalScriptWitness = nil

	p.cleanupInput(in)
	return nil
}

// finalizeP2PKH 最终化P2PKH输入
func (p *Packet) finalizeP2PKH(in *Input, pkScript []byte, value int64) error {
	// 查找对应的签名和公钥
	var pubkey []byte
	var sig []byte

	if len(in.BIP32) > 0 {
		for _, d := range in.BIP32 {
			keyHex := fmt.Sprintf("%x", d.PubKey)
			if s, ok := in.PartialSigs[keyHex]; ok {
				// 验证公钥hash是否匹配
				if bytes.Equal(hash160(d.PubKey), pkScript[3:23]) { // 跳过OP_DUP OP_HASH160 PUSH20
					pubkey = d.PubKey
					sig = s
					break
				}
			}
		}
	}

	if len(pubkey) == 0 || len(sig) == 0 {
		return errors.New("psbt: matching pubkey/sig not found for p2pkh")
	}

	// 构建FinalScriptSig: [sig, pubkey]
	stack := [][]byte{append([]byte(nil), sig...), append([]byte(nil), pubkey...)}
	in.FinalScriptSig = p.buildScriptSig(stack)
	in.FinalScriptWitness = nil

	p.cleanupInput(in)
	return nil
}

// buildScriptSig 构建ScriptSig
func (p *Packet) buildScriptSig(stack [][]byte) []byte {
	var result []byte
	for _, item := range stack {
		if len(item) <= 75 {
			result = append(result, byte(len(item)))
		} else if len(item) <= 255 {
			result = append(result, 0x4c, byte(len(item)))
		} else if len(item) <= 65535 {
			result = append(result, 0x4d, byte(len(item)&0xff), byte(len(item)>>8))
		} else {
			result = append(result, 0x4e, byte(len(item)&0xff), byte(len(item)>>8), byte(len(item)>>16), byte(len(item)>>24))
		}
		result = append(result, item...)
	}
	return result
}

// cleanupInput 清理输入字段
func (p *Packet) cleanupInput(in *Input) {
	// 最小化通用临时字段
	in.PartialSigs = nil
	in.BIP32 = nil
	in.TapScriptStack = nil

	// 配置化清理策略
	if !p.Cleanup.KeepWitnessScript {
		in.WitnessScript = nil
		in.TapLeafScript = nil
		in.TapControlBlock = nil
		in.TapAnnex = nil
	}
	if !p.Cleanup.KeepRedeemScript {
		in.RedeemScript = nil
	}
	if !p.Cleanup.KeepUTXO {
		in.WitnessUtxo = nil
		in.NonWitnessUtxo = nil
	}
}

// FinalizeAll 尝试最终化所有输入
func (p *Packet) FinalizeAll() error {
	for i := range p.Inputs {
		if err := p.FinalizeInput(i); err != nil {
			return err
		}
	}
	return nil
}

// Extract 当所有输入都最终化后，导出可广播的原始交易（wire.MsgTx）
func (p *Packet) Extract() (*wire.MsgTx, error) {
	var m *wire.MsgTx
	if p.IsV0() {
		if p.UnsignedTx == nil {
			return nil, errors.New("psbt: v0 missing unsigned tx")
		}
		m = p.UnsignedTx.Copy()
	} else {
		var err error
		m, err = p.buildMsgTxSkeleton()
		if err != nil {
			return nil, err
		}
	}
	// 写入最终脚本
	for i := range p.Inputs {
		in := p.Inputs[i]
		if len(in.FinalScriptWitness) == 0 && len(in.FinalScriptSig) == 0 {
			return nil, fmt.Errorf("psbt: input %d not finalized", i)
		}
		m.TxIn[i].Witness = in.FinalScriptWitness
		m.TxIn[i].SignatureScript = append([]byte(nil), in.FinalScriptSig...)
	}
	return m, nil
}
