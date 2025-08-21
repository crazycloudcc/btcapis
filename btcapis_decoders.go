package btcapis

import (
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// DecodeAddress 解码地址
func DecodeAddress(addr string) (*types.AddressScriptInfo, error) {
	return decoders.DecodeAddress(addr)
}

// DecodePkScript 解码公钥脚本
func DecodePkScript(pkScript []byte) (*types.AddressInfo, error) {
	return decoders.DecodePkScript(pkScript)
}

// DecodeAsmScript 解码脚本
func DecodeAsmScript(pkScript []byte) (ops []types.ScriptOp, asm string, err error) {
	return decoders.DecodeAsmScript(pkScript)
}

// DisasmScript 反汇编脚本
func DisasmScript(b []byte) (ops []types.ScriptOp, asm string, err error) {
	return decoders.DisasmScript(b)
}

// DecodeRawTx 解码原始交易
func DecodeRawTx(raw []byte) (*types.Tx, error) {
	return decoders.DecodeRawTx(raw)
}
