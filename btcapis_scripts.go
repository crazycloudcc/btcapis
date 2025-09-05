package btcapis

import (
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// 通过地址获取地址的详细信息
func (c *Client) DecodeAddressToScriptInfo(addr string) (*types.AddressScriptInfo, error) {
	return decoders.DecodeAddress(addr)
}

// 通过地址获取锁定脚本
func (c *Client) DecodeAddressToPkScript(addr string) ([]byte, error) {
	return decoders.AddressToPkScript(addr)
}

// 通过地址获取类型
func (c *Client) DecodeAddressToType(addr string) (types.AddressType, error) {
	return decoders.AddressToType(addr)
}

// 通过脚本获取地址信息
func (c *Client) DecodePkScriptToAddressInfo(pkScript []byte) (*types.AddressInfo, error) {
	return decoders.DecodePkScript(pkScript)
}

// 通过脚本获取类型
func (c *Client) DecodePKScriptToType(pkScript []byte) (types.AddressType, error) {
	return decoders.PKScriptToType(pkScript)
}

// 解析脚本为操作码
func (c *Client) DecodePkScriptToAsmString(pkScript []byte) (ops []types.ScriptOp, asm string, err error) {
	return decoders.DecodeAsmScript(pkScript)
}

// 解析一笔交易元数据 => 适用于外部直接输入交易元数据解析结构
func (c *Client) DecodeRawTx(rawtx []byte) (*types.Tx, error) {
	return decoders.DecodeRawTx(rawtx)
}

// 解析一笔交易元数据 => 适用于外部直接输入交易元数据解析结构(十六进制字符串)
func (c *Client) DecodeRawTxString(rawHex string) (*types.Tx, error) {
	return decoders.DecodeRawTxString(rawHex)
}
