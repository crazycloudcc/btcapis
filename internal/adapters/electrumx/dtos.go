// ElectrumX 原生数据结构定义
package electrumx

// BalanceDTO 地址余额数据结构
type BalanceDTO struct {
	Confirmed   int64 `json:"confirmed"`   // 已确认余额（单位：聪）
	Unconfirmed int64 `json:"unconfirmed"` // 未确认余额（单位：聪）
}

// UTXODTO UTXO数据结构
type UTXODTO struct {
	TxHash string `json:"tx_hash"` // 交易哈希
	TxPos  uint32 `json:"tx_pos"`  // 交易输出索引
	Value  int64  `json:"value"`   // 金额（单位：聪）
	Height int64  `json:"height"`  // 区块高度（0表示未确认）
}

// HistoryDTO 交易历史数据结构
type HistoryDTO struct {
	TxHash string `json:"tx_hash"` // 交易哈希
	Height int64  `json:"height"`  // 区块高度（0表示未确认，-1表示未广播）
	Fee    int64  `json:"fee"`     // 交易费用（单位：聪）
}

// TransactionDTO 交易详情数据结构
type TransactionDTO struct {
	Txid          string `json:"txid"`          // 交易ID
	Hash          string `json:"hash"`          // 交易哈希（包含见证数据）
	Version       int32  `json:"version"`       // 交易版本
	Size          int64  `json:"size"`          // 交易大小（字节）
	Vsize         int64  `json:"vsize"`         // 虚拟大小（vbytes）
	Weight        int64  `json:"weight"`        // 权重
	Locktime      uint32 `json:"locktime"`      // 锁定时间
	Hex           string `json:"hex"`           // 原始交易十六进制
	Confirmations int64  `json:"confirmations"` // 确认数
	BlockHash     string `json:"blockhash"`     // 区块哈希
	BlockTime     int64  `json:"blocktime"`     // 区块时间
	Time          int64  `json:"time"`          // 交易时间
	Vin           []Vin  `json:"vin"`           // 输入列表
	Vout          []Vout `json:"vout"`          // 输出列表
}

// Vin 交易输入数据结构
type Vin struct {
	Txid      string    `json:"txid"`                  // 引用的交易ID
	Vout      uint32    `json:"vout"`                  // 引用的输出索引
	ScriptSig ScriptSig `json:"scriptSig"`             // 签名脚本
	Sequence  uint32    `json:"sequence"`              // 序列号
	Witness   []string  `json:"txinwitness,omitempty"` // 见证数据
}

// Vout 交易输出数据结构
type Vout struct {
	Value        float64      `json:"value"`        // 金额（单位：BTC）
	N            uint32       `json:"n"`            // 输出索引
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"` // 锁定脚本
}

// ScriptSig 签名脚本数据结构
type ScriptSig struct {
	Asm string `json:"asm"` // 汇编格式
	Hex string `json:"hex"` // 十六进制格式
}

// ScriptPubKey 锁定脚本数据结构
type ScriptPubKey struct {
	Asm       string   `json:"asm"`       // 汇编格式
	Hex       string   `json:"hex"`       // 十六进制格式
	ReqSigs   int      `json:"reqSigs"`   // 需要的签名数
	Type      string   `json:"type"`      // 脚本类型
	Addresses []string `json:"addresses"` // 地址列表
}

// FeeEstimateDTO 手续费估算数据结构
type FeeEstimateDTO struct {
	FeeRate float64 `json:"feerate"` // 手续费率（单位：BTC/KB）
}

// BlockHeaderDTO 区块头数据结构
type BlockHeaderDTO struct {
	Version       int32  `json:"version"`         // 版本
	PrevBlockHash string `json:"prev_block_hash"` // 上一个区块哈希
	MerkleRoot    string `json:"merkle_root"`     // 默克尔根
	Timestamp     int64  `json:"timestamp"`       // 时间戳
	Bits          uint32 `json:"bits"`            // 难度目标
	Nonce         uint32 `json:"nonce"`           // 随机数
	BlockHeight   int64  `json:"block_height"`    // 区块高度
}

// ServerVersionDTO 服务器版本信息
type ServerVersionDTO struct {
	ServerVersion   string `json:"server_version"`   // 服务器版本
	ProtocolVersion string `json:"protocol_version"` // 协议版本
}

// ServerFeaturesDTO 服务器功能信息
type ServerFeaturesDTO struct {
	GenesisHash   string                 `json:"genesis_hash"`   // 创世区块哈希
	Hosts         map[string]interface{} `json:"hosts"`          // 服务器主机信息
	ProtocolMax   string                 `json:"protocol_max"`   // 最大协议版本
	ProtocolMin   string                 `json:"protocol_min"`   // 最小协议版本
	Pruning       interface{}            `json:"pruning"`        // 修剪信息
	ServerVersion string                 `json:"server_version"` // 服务器版本
	HashFunction  string                 `json:"hash_function"`  // 哈希函数
}

// MempoolDTO 内存池交易数据结构
type MempoolDTO struct {
	TxHash string `json:"tx_hash"` // 交易哈希
	Height int64  `json:"height"`  // 高度（-1表示在内存池中）
	Fee    int64  `json:"fee"`     // 交易费用
}
