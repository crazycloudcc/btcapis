package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/crazycloudcc/btcapis/script"
	"github.com/crazycloudcc/btcapis/types"
)

// AnalyzeInput 解析某个输入的 witness/script 为 OP code 列表；仅在需要时调用。
func AnalyzeInput(t *types.Tx, idx int) (*types.TapscriptInfo, error) {
	if idx < 0 || idx >= len(t.Vin) {
		return nil, fmt.Errorf("vin index out of range")
	}
	in := t.Vin[idx]

	// P2TR？
	if t.Vout != nil { /* 无需依赖输出推断，这里只看 witness 形态 */
	}
	if stack, scr, cb, ok := script.ExtractTapScriptPath(in.Witness); ok {
		ops, asm, err := script.DisasmScript(scr)
		if err != nil {
			return nil, err
		}
		ctrl, err := script.ParseControlBlock(cb)
		if err != nil {
			return nil, err
		}
		ss := make([]string, len(stack))
		for i := range stack {
			ss[i] = hex.EncodeToString(stack[i])
		}
		return &types.TapscriptInfo{
			Path:      "p2tr-script",
			ScriptHex: hex.EncodeToString(scr),
			ASM:       asm,
			Ops:       ops,
			Control:   ctrl,
			StackHex:  ss,
		}, nil
	}

	// P2TR key-path: witness 通常只有一个 64/65B schnorr sig（+可选 annex）
	if len(in.Witness) >= 1 && len(in.Witness[len(in.Witness)-1]) == 64 || len(in.Witness[len(in.Witness)-1]) == 65 {
		// 简化返回：没有脚本与控制块
		ss := make([]string, len(in.Witness))
		for i := range in.Witness {
			ss[i] = hex.EncodeToString(in.Witness[i])
		}
		return &types.TapscriptInfo{
			Path:     "p2tr-key",
			StackHex: ss,
		}, nil
	}

	// P2WPKH/P2WSH 等（可按需扩展）
	ss := make([]string, len(in.Witness))
	for i := range in.Witness {
		ss[i] = hex.EncodeToString(in.Witness[i])
	}
	return &types.TapscriptInfo{
		Path:     "unknown",
		StackHex: ss,
	}, nil
}
