package btcapis

import (
	"github.com/crazycloudcc/btcapis/src/decoders"
	"github.com/crazycloudcc/btcapis/src/types"
)

type Decoders interface {
	// 解码地址
	DecodeAddress(addr string) (*types.AddressScriptInfo, error)
	// 反汇编脚本
	DisasmScript(b []byte) (ops []types.ScriptOp, asm string, err error)
	// 解码原始交易
	DecodeRawTx(raw []byte) (*types.Tx, error)
	// 解码公钥脚本
	DecodePkScript(pkScript []byte) (*types.AddressInfo, error)
	// 解析 Ordinals 信封
	ParseOrdinalEnvelope(scr []byte) (*types.OrdinalsEnvelope, bool, error)
	// 解析 Taproot 控制块
	ParseControlBlock(cb []byte) (types.TapControlBlock, error)
	// 提取 Taproot 脚本路径
	ExtractTapScriptPath(w [][]byte) (stack [][]byte, script []byte, control []byte, ok bool)
	// 计算 Taproot 叶子哈希
	TapLeafHash(leafVersion byte, script []byte) [32]byte
}

// DecodeAddress 解码地址
func (c *Client) DecodeAddress(addr string) (*types.AddressScriptInfo, error) {
	return decoders.DecodeAddress(addr)
}

// DisasmScript 反汇编脚本
func (c *Client) DisasmScript(b []byte) (ops []types.ScriptOp, asm string, err error) {
	return decoders.DisasmScript(b)
}

// DecodeRawTx 解码原始交易
func (c *Client) DecodeRawTx(raw []byte) (*types.Tx, error) {
	return decoders.DecodeRawTx(raw)
}

// DecodePkScript 解码公钥脚本
func (c *Client) DecodePkScript(pkScript []byte) (*types.AddressInfo, error) {
	return decoders.DecodePkScript(pkScript)
}
