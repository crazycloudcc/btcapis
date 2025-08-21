package decoders

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/crazycloudcc/btcapis/src/types"
)

// ParseOrdinalEnvelope 尝试在 tapscript 中寻找并解析 Ordinals TLV。
// 规则：匹配子序列 OP_0 OP_IF "ord" <kv kv kv...> OP_ENDIF。
// 返回 (env, found, err)。found=false 表示不是 ord envelope；err!=nil 表示语法损坏。
func ParseOrdinalEnvelope(scr []byte) (*types.OrdinalsEnvelope, bool, error) {
	ops, _, err := DisasmScript(scr)
	if err != nil {
		return nil, false, fmt.Errorf("disasm: %w", err)
	}
	// 扫描 OP_0 OP_IF
	for i := 0; i+3 < len(ops); i++ {
		if ops[i].Op != "OP_0" || ops[i+1].Op != "OP_IF" {
			continue
		}
		// 期望后面紧跟 "ord"
		j := i + 2
		if !isPush(ops[j]) || !eqHex(ops[j].DataHex, "6f7264") { // "ord"
			continue
		}
		// 解析 TLV，直到遇到 OP_ENDIF
		k := j + 1
		env := &types.OrdinalsEnvelope{Records: make([]types.OrdinalsRecord, 0, 8)}
		var body []byte
		for k < len(ops) {
			if ops[k].Op == "OP_ENDIF" {
				// 结束
				if len(body) > 0 {
					env.BodyHex = hex.EncodeToString(body)
				}
				return env, true, nil
			}

			keyByte, ok := parseKey(ops[k])
			if !ok {
				return nil, true, fmt.Errorf("ord: expect key at %d, got %s", k, ops[k].Op)
			}

			if k+1 >= len(ops) {
				return nil, true, fmt.Errorf("ord: missing value after key at %d", k)
			}

			valHex, vok := parseValue(ops[k+1])
			if !vok {
				return nil, true, fmt.Errorf("ord: expect value at %d, got %s", k+1, ops[k+1].Op)
			}

			keyHex := fmt.Sprintf("%02x", keyByte)
			env.Records = append(env.Records, types.OrdinalsRecord{
				KeyHex:   keyHex,
				ValueHex: valHex,
			})

			// 常用键处理
			if keyByte == 0x00 { // Body（可分片）
				if valHex != "" {
					vb, _ := hex.DecodeString(valHex)
					body = append(body, vb...)
				}
			} else if keyByte == 0x01 { // Content-Type
				env.ContentType = asciiOf(valHex)
			}
			k += 2

			env.Records = append(env.Records, types.OrdinalsRecord{KeyHex: keyHex, ValueHex: valHex})

			// 识别常用键：0=body(可分片)，1=content-type（ASCII）
			if keyIsByte(keyHex, 0x00) {
				vb, _ := hex.DecodeString(valHex)
				body = append(body, vb...)
			} else if keyIsByte(keyHex, 0x01) {
				env.ContentType = asciiOf(valHex)
			}
			k += 2
		}
		// 未遇到 ENDIF
		return nil, true, fmt.Errorf("ord: missing OP_ENDIF")
	}
	return nil, false, nil
}

// —— 辅助 ——

// isPush: 简单判断是否 data-push（我们在 DisasmScript 里把所有 push 标成 OP_PUSH_* 或 OP_PUSHDATA*）
func isPush(op types.ScriptOp) bool {
	return strings.HasPrefix(op.Op, "OP_PUSH_") || op.Op == "OP_PUSHDATA1" || op.Op == "OP_PUSHDATA2" || op.Op == "OP_PUSHDATA4"
}

// keyIsByte: key 恰好 1 字节并等于 v
func keyIsByte(hexStr string, v byte) bool {
	if len(hexStr) != 2 {
		return false
	}
	b, err := hex.DecodeString(hexStr)
	return err == nil && b[0] == v
}

func eqHex(h, want string) bool { return strings.EqualFold(h, want) }

// 把十六进制按 ASCII 解码（非 ASCII 字节会变为 U+FFFD 替代符）
func asciiOf(hexStr string) string {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return ""
	}
	// 简单返回 string(b) 即可；内容类型一般 ASCII。
	return string(b)
}

// 读取 key：支持 OP_0..OP_16，或 push(1-byte)
func parseKey(op types.ScriptOp) (byte, bool) {
	// 小整数：OP_0..OP_16
	switch op.Op {
	case "OP_0":
		return 0x00, true
	case "OP_1":
		return 0x01, true
	case "OP_2":
		return 0x02, true
	case "OP_3":
		return 0x03, true
	case "OP_4":
		return 0x04, true
	case "OP_5":
		return 0x05, true
	case "OP_6":
		return 0x06, true
	case "OP_7":
		return 0x07, true
	case "OP_8":
		return 0x08, true
	case "OP_9":
		return 0x09, true
	case "OP_10":
		return 0x0a, true
	case "OP_11":
		return 0x0b, true
	case "OP_12":
		return 0x0c, true
	case "OP_13":
		return 0x0d, true
	case "OP_14":
		return 0x0e, true
	case "OP_15":
		return 0x0f, true
	case "OP_16":
		return 0x10, true
	}
	// 单字节 push：DataLen==1
	if strings.HasPrefix(op.Op, "OP_PUSH_") || op.Op == "OP_PUSHDATA1" || op.Op == "OP_PUSHDATA2" || op.Op == "OP_PUSHDATA4" {
		if op.DataLen == 1 && len(op.DataHex) == 2 {
			b, err := hex.DecodeString(op.DataHex)
			if err == nil {
				return b[0], true
			}
		}
	}
	return 0, false
}

// 读取 value：允许 OP_0（空串）或 pushdata
func parseValue(op types.ScriptOp) (hexStr string, ok bool) {
	if op.Op == "OP_0" {
		return "", true // 空值
	}
	if strings.HasPrefix(op.Op, "OP_PUSH_") || op.Op == "OP_PUSHDATA1" || op.Op == "OP_PUSHDATA2" || op.Op == "OP_PUSHDATA4" {
		return op.DataHex, true
	}
	return "", false
}
