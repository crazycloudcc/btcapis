package psbt

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
)

// FinalizeInput 根据已收集到的 PartialSigs/WitnessScript 等，生成最终的 ScriptSig/ScriptWitness。
// 支持 p2wpkh / p2wsh(v0) / p2tr-keypath。复杂脚本（多签顺序占位等）后续增强。
func (p *Packet) FinalizeInput(i int) error {
	in := p.MustInput(i)
	var pkScript []byte
	var value int64
	if in.WitnessUtxo != nil {
		pkScript = in.WitnessUtxo.PkScript
		value = in.WitnessUtxo.Value
	} else if in.NonWitnessUtxo != nil {
		return errors.New("psbt: finalize requires witnessUtxo for segwit inputs")
	} else {
		return errors.New("psbt: missing utxo to finalize")
	}

	isP2WPKH := len(pkScript) == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14
	isP2WSH := len(pkScript) == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20
	isP2TR := len(pkScript) == 34 && pkScript[0] == 0x51 && pkScript[1] == 0x20

	if len(in.PartialSigs) == 0 {
		return errors.New("psbt: need partial sig to finalize")
	}

	if isP2WPKH {
		if len(in.BIP32) == 0 {
			return errors.New("psbt: need pubkey to finalize p2wpkh")
		}
		var pubkey []byte
		var sig []byte
		for _, d := range in.BIP32 {
			keyHex := fmt.Sprintf("%x", d.PubKey)
			if s, ok := in.PartialSigs[keyHex]; ok {
				pubkey = d.PubKey
				sig = s
				break
			}
		}
		if len(pubkey) == 0 || len(sig) == 0 {
			return errors.New("psbt: matching pubkey/sig not found")
		}
		in.FinalScriptWitness = wire.TxWitness{append([]byte(nil), sig...), append([]byte(nil), pubkey...)}
		in.FinalScriptSig = nil
		in.PartialSigs = nil
		in.RedeemScript = nil
		in.WitnessScript = nil
		_ = value
		return nil
	}

	if isP2WSH {
		if len(in.WitnessScript) == 0 {
			return errors.New("psbt: missing witnessScript for p2wsh finalize")
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
		stack := make([][]byte, 0, len(in.PartialSigs)+2)
		if isMultisig {
			stack = append(stack, []byte{}) // 占位空元素
		}
		for _, d := range in.BIP32 {
			keyHex := fmt.Sprintf("%x", d.PubKey)
			if s, ok := in.PartialSigs[keyHex]; ok {
				stack = append(stack, append([]byte(nil), s...))
			}
		}
		stack = append(stack, append([]byte(nil), ws...))
		if len(stack) < 2 {
			return errors.New("psbt: not enough sigs to finalize p2wsh")
		}
		in.FinalScriptWitness = wire.TxWitness(stack)
		in.FinalScriptSig = nil
		in.PartialSigs = nil
		in.RedeemScript = nil
		// 可选择清理或保留 WitnessScript；此处保留供上层检查
		return nil
	}

	if isP2TR {
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
			in.PartialSigs = nil
			in.RedeemScript = nil
			in.WitnessScript = nil
			return nil
		}

		// 否则回退为 keypath：最终 witness = [schnorr_sig]
		for _, s := range in.PartialSigs {
			in.FinalScriptWitness = wire.TxWitness{append([]byte(nil), s...)}
			in.FinalScriptSig = nil
			in.PartialSigs = nil
			in.RedeemScript = nil
			in.WitnessScript = nil
			return nil
		}
		return errors.New("psbt: no schnorr sig for p2tr finalize")
	}

	return errors.New("psbt: unsupported script type to finalize")
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
			return nil, errors.New("psbt: input not finalized")
		}
		m.TxIn[i].Witness = in.FinalScriptWitness
		m.TxIn[i].SignatureScript = append([]byte(nil), in.FinalScriptSig...)
	}
	return m, nil
}
