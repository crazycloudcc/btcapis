package types

import "github.com/btcsuite/btcd/txscript"

// AddressType 标注常见地址/脚本族群。
type AddressType string

const (
	AddrP2PK    AddressType = "p2pk"
	AddrP2PKH   AddressType = "p2pkh"
	AddrP2SH    AddressType = "p2sh"
	AddrP2WPKH  AddressType = "p2wpkh"
	AddrP2WSH   AddressType = "p2wsh"
	AddrP2TR    AddressType = "p2tr"
	AddrUnknown AddressType = "unknown"
)

// AddressBalanceInfo 地址余额信息
type AddressBalanceInfo struct {
	Address     string // 地址
	Confirmed   int64  // 已确认余额（聪）
	Unconfirmed int64  // 未确认余额（聪）
	Total       int64  // 总余额（聪）
	Error       error  // 查询错误（如果有）
}

// AddressInfo 结构体：存储地址解析后的信息
type AddressInfo struct {
	PKScript  []byte               // 原始脚本
	Typ       AddressType          // 地址类型
	Cls       txscript.ScriptClass // 脚本类型
	ReqSigs   int                  // 需要签名数（多签时有意义）
	Addresses []string             // 可能为 0/1/N
}

// AddressScriptInfo 结构体：存储地址解析后的脚本信息
// 包含脚本类型、各种哈希值、见证版本等关键信息
type AddressScriptInfo struct {
	Address         string               // 地址
	Typ             AddressType          // 地址类型
	Cls             txscript.ScriptClass // 脚本类型
	ScriptPubKeyHex []byte               // 脚本哈希 => PKScript
	ScriptAsm       string               // 脚本汇编

	// “哈希/程序”层面（地址能直接给出的）
	PubKeyHashHex       []byte // P2PKH: 20 bytes
	RedeemScriptHashHex []byte // P2SH:  20 bytes

	// SegWit
	IsWitness         bool   // 是否为见证地址
	WitnessVersion    int    // 见证版本：0（SegWit v0）、1（Taproot）、-1（非SegWit）
	WitnessProgramHex []byte // v0: 20/32 bytes; v1+: 32 bytes
	WitnessProgramLen int    // 见证数据长度：20字节或32字节
	BechEncoding      string // bech32 / bech32m

	// Taproot
	TaprootOutputKeyHex []byte // 等同于 witness program (v=1, 32B x-only pubkey)
}
