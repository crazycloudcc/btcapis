// Package types 提供跨模块共享的轻量类型定义
package types

type AddressType string

const (
	AddrP2PKH   AddressType = "p2pkh"
	AddrP2SH    AddressType = "p2sh"
	AddrP2WPKH  AddressType = "p2wpkh"
	AddrP2WSH   AddressType = "p2wsh"
	AddrP2TR    AddressType = "p2tr"
	AddrUnknown AddressType = "unknown"
)

type AddressInfo struct {
	Type         AddressType
	Network      Network
	Program      []byte   // witness program 或 hash160，可留空
	ScriptPubKey []byte   // 若从地址派生可填，否则留空
	Address      string   // 规范化地址（如果能解析出来）
	Tags         []string // 可选：额外标记
}
