// Package script 脚本构建器
package script

import (
	"fmt"
)

// OpCode 操作码定义
type OpCode struct {
	Code byte
	Name string
	Data []byte
}

// Builder 脚本构建器
type Builder struct {
	ops []OpCode
}

// NewBuilder 创建新的脚本构建器
func NewBuilder() *Builder {
	return &Builder{
		ops: make([]OpCode, 0),
	}
}

// AddOp 添加操作码
func (b *Builder) AddOp(op OpCode) *Builder {
	b.ops = append(b.ops, op)
	return b
}

// AddData 添加数据
func (b *Builder) AddData(data []byte) *Builder {
	// 根据数据长度选择适当的操作码
	if len(data) == 0 {
		b.ops = append(b.ops, OpCode{Code: 0x00, Name: "OP_0"})
	} else if len(data) <= 75 {
		b.ops = append(b.ops, OpCode{Code: byte(len(data)), Name: "OP_DATA"})
		b.ops = append(b.ops, OpCode{Code: 0x00, Name: "DATA", Data: data})
	} else if len(data) <= 255 {
		b.ops = append(b.ops, OpCode{Code: 0x4c, Name: "OP_PUSHDATA1"})
		b.ops = append(b.ops, OpCode{Code: byte(len(data)), Name: "LENGTH"})
		b.ops = append(b.ops, OpCode{Code: 0x00, Name: "DATA", Data: data})
	} else {
		b.ops = append(b.ops, OpCode{Code: 0x4d, Name: "OP_PUSHDATA2"})
		b.ops = append(b.ops, OpCode{Code: byte(len(data) & 0xff), Name: "LENGTH_LOW"})
		b.ops = append(b.ops, OpCode{Code: byte(len(data) >> 8), Name: "LENGTH_HIGH"})
		b.ops = append(b.ops, OpCode{Code: 0x00, Name: "DATA", Data: data})
	}
	return b
}

// Build 构建脚本
func (b *Builder) Build() ([]byte, error) {
	// TODO: 实现脚本构建逻辑
	return nil, fmt.Errorf("not implemented")
}

// BuildP2PKH 构建P2PKH脚本
func (b *Builder) BuildP2PKH(pubKeyHash []byte) ([]byte, error) {
	if len(pubKeyHash) != 20 {
		return nil, fmt.Errorf("invalid public key hash length: %d", len(pubKeyHash))
	}

	builder := NewBuilder()
	builder.AddOp(OpCode{Code: 0x76, Name: "OP_DUP"})
	builder.AddOp(OpCode{Code: 0xa9, Name: "OP_HASH160"})
	builder.AddData(pubKeyHash)
	builder.AddOp(OpCode{Code: 0x88, Name: "OP_EQUALVERIFY"})
	builder.AddOp(OpCode{Code: 0xac, Name: "OP_CHECKSIG"})

	return builder.Build()
}

// BuildP2SH 构建P2SH脚本
func (b *Builder) BuildP2SH(scriptHash []byte) ([]byte, error) {
	if len(scriptHash) != 20 {
		return nil, fmt.Errorf("invalid script hash length: %d", len(scriptHash))
	}

	builder := NewBuilder()
	builder.AddOp(OpCode{Code: 0xa9, Name: "OP_HASH160"})
	builder.AddData(scriptHash)
	builder.AddOp(OpCode{Code: 0x87, Name: "OP_EQUAL"})

	return builder.Build()
}

// BuildP2WPKH 构建P2WPKH脚本
func (b *Builder) BuildP2WPKH(pubKeyHash []byte) ([]byte, error) {
	if len(pubKeyHash) != 20 {
		return nil, fmt.Errorf("invalid public key hash length: %d", len(pubKeyHash))
	}

	builder := NewBuilder()
	builder.AddOp(OpCode{Code: 0x00, Name: "OP_0"})
	builder.AddData(pubKeyHash)

	return builder.Build()
}

// BuildP2WSH 构建P2WSH脚本
func (b *Builder) BuildP2WSH(scriptHash []byte) ([]byte, error) {
	if len(scriptHash) != 32 {
		return nil, fmt.Errorf("invalid script hash length: %d", len(scriptHash))
	}

	builder := NewBuilder()
	builder.AddOp(OpCode{Code: 0x00, Name: "OP_0"})
	builder.AddData(scriptHash)

	return builder.Build()
}

// BuildP2TR 构建P2TR脚本
func (b *Builder) BuildP2TR(pubKey []byte) ([]byte, error) {
	if len(pubKey) != 32 {
		return nil, fmt.Errorf("invalid public key length: %d", len(pubKey))
	}

	builder := NewBuilder()
	builder.AddOp(OpCode{Code: 0x51, Name: "OP_1"})
	builder.AddData(pubKey)

	return builder.Build()
}
