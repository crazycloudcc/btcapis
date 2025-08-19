package script

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/crazycloudcc/btcapis/types"
)

// DisasmScript 将原始脚本字节反汇编为 ops + ASM 字符串；遇到非法 push 会返回错误。
func DisasmScript(b []byte) (ops []types.ScriptOp, asm string, err error) {
	var out []types.ScriptOp
	var asmParts []string
	i := 0
	for i < len(b) {
		op := b[i]
		i++

		// 小于等于 0x4b：直接 push N 字节
		if op >= 0x01 && op <= 0x4b {
			n := int(op)
			if i+n > len(b) {
				return nil, "", fmt.Errorf("short push: need=%d have=%d", n, len(b)-i)
			}
			data := b[i : i+n]
			i += n
			h := hex.EncodeToString(data)
			out = append(out, types.ScriptOp{Op: fmt.Sprintf("OP_PUSH_%d", n), DataHex: h, DataLen: n})
			asmParts = append(asmParts, h)
			continue
		}

		// OP_PUSHDATA1/2/4
		if op == 0x4c || op == 0x4d || op == 0x4e {
			var n, need int
			switch op {
			case 0x4c:
				if i+1 > len(b) {
					return nil, "", fmt.Errorf("PUSHDATA1 short")
				}
				n = int(b[i])
				i++
			case 0x4d:
				if i+2 > len(b) {
					return nil, "", fmt.Errorf("PUSHDATA2 short")
				}
				n = int(b[i]) | int(b[i+1])<<8
				i += 2
			case 0x4e:
				if i+4 > len(b) {
					return nil, "", fmt.Errorf("PUSHDATA4 short")
				}
				n = int(b[i]) | int(b[i+1])<<8 | int(b[i+2])<<16 | int(b[i+3])<<24
				i += 4
			}
			need = n
			if n < 0 || i+n > len(b) {
				return nil, "", fmt.Errorf("PUSHDATA need=%d have=%d", need, len(b)-i)
			}
			data := b[i : i+n]
			i += n
			h := hex.EncodeToString(data)
			out = append(out, types.ScriptOp{Op: types.OpcodeName(op), DataHex: h, DataLen: n})
			asmParts = append(asmParts, h)
			continue
		}

		// 常规 opcode
		nm := types.OpcodeName(op)
		out = append(out, types.ScriptOp{Op: nm})
		asmParts = append(asmParts, nm)
	}
	return out, strings.Join(asmParts, " "), nil
}
