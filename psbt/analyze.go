package psbt

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// AnalyzeInputReport 描述单个输入的自检结果
type AnalyzeInputReport struct {
	Index      int
	ScriptType string
	HasUtxo    bool
	IsFinal    bool
	Missing    []string
}

// AnalyzeReport 汇总整笔 PSBT 的自检结果
type AnalyzeReport struct {
	Version     int
	InputCount  int
	OutputCount int
	Inputs      []AnalyzeInputReport
	CanFinalize bool
}

// Analyze 对当前 PSBT 做自检，给出每个输入缺失项与可最终化评估
func (p *Packet) Analyze() AnalyzeReport {
	res := AnalyzeReport{Version: p.Version, InputCount: len(p.Inputs), OutputCount: len(p.Outputs)}
	res.Inputs = make([]AnalyzeInputReport, 0, len(p.Inputs))

	allReady := true
	for i := range p.Inputs {
		in := p.Inputs[i]
		ir := AnalyzeInputReport{Index: i}

		// 已最终化
		if len(in.FinalScriptWitness) > 0 || len(in.FinalScriptSig) > 0 {
			ir.IsFinal = true
			// 推断脚本类型：尽力从已知字段/UTXO获得
			if spk := inputPkScriptForAnalyze(p, i); len(spk) > 0 {
				ir.ScriptType = classifyScript(spk)
			}
			ir.HasUtxo = in.WitnessUtxo != nil || in.NonWitnessUtxo != nil
			res.Inputs = append(res.Inputs, ir)
			continue
		}

		var pkScript []byte
		// 推断 UTXO/脚本
		if in.WitnessUtxo != nil {
			pkScript = in.WitnessUtxo.PkScript
			ir.HasUtxo = true
		} else if in.NonWitnessUtxo != nil {
			// 若为 v0，可校验 txid/vout
			if p.IsV0() && p.UnsignedTx != nil {
				prev := p.UnsignedTx.TxIn[i].PreviousOutPoint
				if in.NonWitnessUtxo.TxHash() == prev.Hash && int(prev.Index) < len(in.NonWitnessUtxo.TxOut) {
					pkScript = in.NonWitnessUtxo.TxOut[prev.Index].PkScript
					ir.HasUtxo = true
				} else {
					ir.Missing = append(ir.Missing, "non_witness_utxo_mismatch")
				}
			} else {
				// 无法校验，只记为有UTXO
				ir.HasUtxo = true
			}
		} else {
			ir.HasUtxo = false
			ir.Missing = append(ir.Missing, "utxo")
		}

		// 识别脚本类型
		sType := classifyScript(pkScript)
		// P2SH 可进一步根据 RedeemScript 细分
		if sType == "p2sh" && len(in.RedeemScript) > 0 {
			inner := classifyScript(in.RedeemScript)
			sType = fmt.Sprintf("p2sh(%s)", inner)
		}
		ir.ScriptType = sType

		// 根据脚本类型判断缺失项（偏向最终化阶段需求）
		switch sType {
		case "p2wpkh":
			if !ir.HasUtxo {
				// segwit 最终化允许仅 NonWitnessUtxo，但 analyze 只标记 utxo
			}
			if len(in.PartialSigs) == 0 {
				ir.Missing = append(ir.Missing, "partial_sig")
			}
			// 校验 pubkey hash 是否可匹配
			if len(pkScript) == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14 {
				if !hasMatchingP2WPKHPubkey(in, pkScript[2:]) {
					ir.Missing = append(ir.Missing, "pubkey_match")
				}
			}
		case "p2wsh":
			if len(in.WitnessScript) == 0 {
				ir.Missing = append(ir.Missing, "witness_script")
			} else {
				// 校验 hash 一致
				if !(len(pkScript) == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20) ||
					!bytes.Equal(pkScript[2:], sha256Bytes(in.WitnessScript)) {
					ir.Missing = append(ir.Missing, "witness_script_mismatch")
				}
			}
			if len(in.PartialSigs) == 0 {
				ir.Missing = append(ir.Missing, "partial_sig")
			}
		case "p2pkh":
			if in.NonWitnessUtxo == nil {
				ir.Missing = append(ir.Missing, "non_witness_utxo")
			}
			if len(in.PartialSigs) == 0 {
				ir.Missing = append(ir.Missing, "partial_sig")
			}
		case "p2sh":
			if len(in.RedeemScript) == 0 {
				ir.Missing = append(ir.Missing, "redeem_script")
			} else {
				inner := classifyScript(in.RedeemScript)
				switch inner {
				case "p2wpkh":
					if len(in.PartialSigs) == 0 {
						ir.Missing = append(ir.Missing, "partial_sig")
					}
				case "p2wsh":
					if len(in.WitnessScript) == 0 {
						ir.Missing = append(ir.Missing, "witness_script")
					}
					if len(in.PartialSigs) == 0 {
						ir.Missing = append(ir.Missing, "partial_sig")
					}
				default:
					if len(in.PartialSigs) == 0 {
						ir.Missing = append(ir.Missing, "partial_sig")
					}
				}
			}
		case "p2tr":
			// keypath or scriptpath
			if len(in.TapLeafScript) > 0 || len(in.TapControlBlock) > 0 {
				if len(in.TapLeafScript) == 0 || len(in.TapControlBlock) == 0 {
					ir.Missing = append(ir.Missing, "tap_script_path_data")
				}
				// 允许无签名脚本；若需要签名，仍检查
				if len(in.PartialSigs) == 0 && len(in.TapScriptSigs) == 0 {
					// 脚本可无需签名，提示为可选
				}
			} else {
				if len(in.PartialSigs) == 0 && len(in.TapKeySig) == 0 {
					ir.Missing = append(ir.Missing, "tap_key_sig")
				}
			}
		default:
			// 无法识别脚本，保守要求有签名或脚本材料
			if len(in.PartialSigs) == 0 && len(in.FinalScriptSig) == 0 && len(in.FinalScriptWitness) == 0 {
				ir.Missing = append(ir.Missing, "partial_sig_or_final")
			}
		}

		ir.IsFinal = len(ir.Missing) == 0 && ir.HasUtxo
		if !ir.IsFinal {
			allReady = false
		}
		res.Inputs = append(res.Inputs, ir)
	}

	res.CanFinalize = allReady
	return res
}

func inputPkScriptForAnalyze(p *Packet, i int) []byte {
	in := p.Inputs[i]
	if in.WitnessUtxo != nil {
		return in.WitnessUtxo.PkScript
	}
	if in.NonWitnessUtxo != nil {
		if p.IsV0() && p.UnsignedTx != nil {
			prev := p.UnsignedTx.TxIn[i].PreviousOutPoint
			if in.NonWitnessUtxo.TxHash() == prev.Hash && int(prev.Index) < len(in.NonWitnessUtxo.TxOut) {
				return in.NonWitnessUtxo.TxOut[prev.Index].PkScript
			}
		}
	}
	return nil
}

func hasMatchingP2WPKHPubkey(in *Input, expectedHash []byte) bool {
	// 1) 从 BIP32 公钥匹配
	for _, d := range in.BIP32 {
		if bytes.Equal(hash160(d.PubKey), expectedHash) {
			return true
		}
	}
	// 2) 从 PartialSigs 的 key 解析压缩公钥并匹配
	for keyHex := range in.PartialSigs {
		if len(keyHex) == 66 { // 33字节压缩公钥hex
			pubkey, err := hex.DecodeString(keyHex)
			if err == nil && len(pubkey) == 33 && bytes.Equal(hash160(pubkey), expectedHash) {
				return true
			}
		}
	}
	return false
}

func sha256Bytes(b []byte) []byte {
	sum := sha256.Sum256(b)
	return sum[:]
}
