package types

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
)

// ScriptOp 表示一条脚本指令；若为数据推送，DataHex/DataLen 会被填充。
type ScriptOp struct {
	Op      string `json:"op"`                 // 操作码
	DataHex string `json:"data_hex,omitempty"` // 数据十六进制
	DataLen int    `json:"data_len,omitempty"` // 数据长度，单位：字节
}

var dictOp2Name map[byte]string

// 将 opcode 映射为可读名称
func OpcodeName(op byte) string {

	if dictOp2Name == nil {
		dictOp2Name = make(map[byte]string, len(txscript.OpcodeByName))
		for name, op := range txscript.OpcodeByName {
			dictOp2Name[byte(op)] = name
		}
	}

	if name, ok := dictOp2Name[op]; ok {
		return name
	}

	return fmt.Sprintf("OP_%d", int(op))
}
