// Package types 脚本相关类型定义
package types

// ScriptType 表示脚本类型
type ScriptType string

const (
	ScriptTypeP2PKH    ScriptType = "p2pkh"
	ScriptTypeP2SH     ScriptType = "p2sh"
	ScriptTypeP2WPKH   ScriptType = "p2wpkh"
	ScriptTypeP2WSH    ScriptType = "p2wsh"
	ScriptTypeP2TR     ScriptType = "p2tr"
	ScriptTypeMultiSig ScriptType = "multisig"
	ScriptTypeOpReturn ScriptType = "op_return"
	ScriptTypeUnknown  ScriptType = "unknown"
)

// ScriptInfo 表示脚本信息
type ScriptInfo struct {
	Type         ScriptType `json:"type"`
	Hex          string     `json:"hex"`
	ASM          string     `json:"asm"`
	Addresses    []string   `json:"addresses"`
	RequiredSigs int        `json:"required_sigs"`
	TotalSigs    int        `json:"total_sigs"`
	IsValid      bool       `json:"is_valid"`
}

// OpCode 表示操作码
type OpCode struct {
	Code        byte   `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Length      int    `json:"length"`
}

// ScriptTemplate 表示脚本模板
type ScriptTemplate struct {
	Name        string            `json:"name"`
	Type        ScriptType        `json:"type"`
	Template    string            `json:"template"`
	Parameters  []ScriptParameter `json:"parameters"`
	Description string            `json:"description"`
}

// ScriptParameter 表示脚本参数
type ScriptParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}
