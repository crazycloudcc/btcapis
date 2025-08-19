package types

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

// AddressInfo 结构体：存储地址解析后的信息
type AddressInfo struct {
	Network   Network     // 网络
	Typ       AddressType // 地址类型
	ReqSigs   int         // 需要签名数（多签时有意义）
	Addresses []string    // 可能为 0/1/N
}

// AddressScriptInfo 结构体：存储地址解析后的脚本信息
// 包含脚本类型、各种哈希值、见证版本等关键信息
type AddressScriptInfo struct {
	ScriptType        AddressType // 脚本类型：P2PKH / P2SH / P2WPKH / P2WSH / P2TR
	PubKeyHash        []byte      // 公钥哈希：20字节（P2PKH）或32字节（P2TR）
	RedeemScriptHash  []byte      // 赎回脚本哈希：20字节（P2SH）
	WitnessScriptHash []byte      // 见证脚本哈希：20字节（P2WPKH）或32字节（P2WSH）
	WitnessVersion    int         // 见证版本：0（SegWit v0）、1（Taproot）、-1（非SegWit）
	WitnessLen        int         // 见证数据长度：20字节或32字节
	TaprootKey        []byte      // Taproot调整后的公钥：32字节
}
