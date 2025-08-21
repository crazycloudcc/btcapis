package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/adapters/decoders"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// AnalyzeInput 解析某个输入的 witness/script 为 OP code 列表；仅在需要时调用。
func AnalyzeTxInWithIdx(t *types.Tx, idx int) (*types.TapscriptInfo, error) {
	if idx < 0 || idx >= len(t.TxIn) {
		return nil, fmt.Errorf("vin index out of range")
	}
	in := t.TxIn[idx]

	return AnalyzeTxIn(t, &in)
}

func AnalyzeTxIn(t *types.Tx, in *types.TxIn) (*types.TapscriptInfo, error) {
	// 1) 尝试识别 Taproot 脚本路径
	if stack, scr, cb, ok := decoders.ExtractTapScriptPath(in.Witness); ok {
		ops, asm, err := decoders.DisasmScript(scr)
		if err == nil {
			ctrl, err2 := decoders.ParseControlBlock(cb)
			if err2 != nil {
				// 极少数场景：长度过关但语义不对——当作非 taproot
			} else {
				ss := make([]string, len(stack))
				for i := range stack {
					ss[i] = hex.EncodeToString(stack[i])
				}

				info := &types.TapscriptInfo{
					Path:      "p2tr-script",
					ScriptHex: hex.EncodeToString(scr),
					ASM:       asm,
					Ops:       ops,
					Control:   ctrl,
					StackHex:  ss,
				}
				// Ordinals TLV（宽容：解析失败不报错，仅不填 Ord）
				if env, found, err3 := decoders.ParseOrdinalEnvelope(scr); found && err3 == nil {
					info.Ord = env
				}
				return info, nil
			}
		}
		// 反汇编失败/控制块语义失败：退回到非 Taproot 流程
	}

	// 2) 非 Taproot：尝试识别 P2WPKH（witness=[sig,pubkey33]）
	if len(in.Witness) == 2 && len(in.Witness[1]) == 33 {
		ss := make([]string, 2)
		ss[0] = hex.EncodeToString(in.Witness[0]) // sig
		ss[1] = hex.EncodeToString(in.Witness[1]) // pubkey
		return &types.TapscriptInfo{
			Path:     "p2wpkh", // 标注路径类型，字段仍用同一结构返回
			StackHex: ss,
		}, nil
	}

	// 3) 其他（P2WSH / P2SH-P2WPKH / 未识别）
	ss := make([]string, len(in.Witness))
	for i := range in.Witness {
		ss[i] = hex.EncodeToString(in.Witness[i])
	}
	return &types.TapscriptInfo{Path: "unknown", StackHex: ss}, nil
}
