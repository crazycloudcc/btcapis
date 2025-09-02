// rpc原生数据结构定义
package bitcoindrpc

// UTXO数据结构
type UTXODTO struct {
	TxID         string  `json:"txid"`         // 交易ID
	Vout         uint32  `json:"vout"`         // 输出索引
	ScriptPubKey string  `json:"scriptPubKey"` // 脚本公钥
	Desc         string  `json:"desc"`         // 描述
	AmountBTC    float64 `json:"amount"`       // 金额
	Height       int64   `json:"height"`       // 高度
}

// 扫描UTXO结果数据结构(模块内部使用, 不需要对外暴露)
type scanResult struct {
	Success       bool      `json:"success"`        // 是否成功
	Bestblock     string    `json:"bestblock"`      // 最佳区块
	Height        int64     `json:"height"`         // 高度
	SearchedItems int64     `json:"searched_items"` // 搜索项
	Unspents      []UTXODTO `json:"unspents"`       // UTXO集
}

// 估算交易费率数据结构
type FeeRateSmartDTO struct {
	Feerate float64 `json:"feerate"` // 交易费率(BTC/KB)
	Blocks  int     `json:"blocks"`  // 目标区块数
}

// 区块头数据结构
type BlockHeaderDTO struct {
	Hash              string  `json:"hash"`              // 区块哈希
	Confirmations     int     `json:"confirmations"`     // 确认数
	Height            int     `json:"height"`            // 高度
	Version           int     `json:"version"`           // 版本
	VersionHex        string  `json:"versionHex"`        // 版本十六进制
	MerkleRoot        string  `json:"merkleroot"`        // 默克尔根
	Time              int     `json:"time"`              // 时间
	MedianTime        int     `json:"mediantime"`        // 中位时间
	Nonce             int     `json:"nonce"`             // 随机数
	Bits              string  `json:"bits"`              // 难度位
	Difficulty        float64 `json:"difficulty"`        // 难度
	Chainwork         string  `json:"chainwork"`         // 工作量
	NTx               int     `json:"nTx"`               // 交易数
	PreviousBlockHash string  `json:"previousblockhash"` // 上一个区块哈希
	NextBlockHash     string  `json:"nextblockhash"`     // 下一个区块哈希
}

// 区块数据结构
type BlockDTO struct {
	Hash              string   `json:"hash"`              // 区块哈希
	Confirmations     int      `json:"confirmations"`     // 确认数
	Height            int      `json:"height"`            // 高度
	Version           int      `json:"version"`           // 版本
	VersionHex        string   `json:"versionHex"`        // 版本十六进制
	MerkleRoot        string   `json:"merkleroot"`        // 默克尔根
	Time              int      `json:"time"`              // 时间
	MedianTime        int      `json:"mediantime"`        // 中位时间
	Nonce             int      `json:"nonce"`             // 随机数
	Bits              string   `json:"bits"`              // 难度位
	Difficulty        float64  `json:"difficulty"`        // 难度
	Chainwork         string   `json:"chainwork"`         // 工作量
	NTx               int      `json:"nTx"`               // 交易数
	PreviousBlockHash string   `json:"previousblockhash"` // 上一个区块哈希
	NextBlockHash     string   `json:"nextblockhash"`     // 下一个区块哈希
	StrippedSize      int      `json:"strippedsize"`      // 剥离大小
	Size              int      `json:"size"`              // 大小
	Weight            int      `json:"weight"`            // 权重
	Tx                []string `json:"tx"`                // 交易
}

// 地址信息数据结构
type AddressInfoDTO struct {
	Address        string   `json:"address"`         // 地址
	ScriptPubKey   string   `json:"scriptPubKey"`    // 脚本公钥
	IsMine         bool     `json:"ismine"`          // 是否属于钱包
	Solvable       bool     `json:"solvable"`        // 是否可解
	Desc           string   `json:"desc"`            // 描述
	ParentDesc     string   `json:"parent_desc"`     // 父描述
	IsWatchOnly    bool     `json:"iswatchonly"`     // 是否只读
	IsScript       bool     `json:"isscript"`        // 是否脚本
	IsWitness      bool     `json:"iswitness"`       // 是否见证
	WitnessVersion int      `json:"witness_version"` // 见证版本
	WitnessProgram string   `json:"witness_program"` // 见证程序
	IsChange       bool     `json:"ischange"`        // 是否改变
	Labels         []string `json:"labels"`          // 标签
}

type ValidateAddressDTO struct {
	IsValid        bool   `json:"isvalid"`         // 是否有效
	Address        string `json:"address"`         // 地址
	ScriptPubKey   string `json:"scriptPubKey"`    // 脚本公钥
	IsScript       bool   `json:"isscript"`        // 是否脚本
	IsWitness      bool   `json:"iswitness"`       // 是否见证
	WitnessVersion int    `json:"witness_version"` // 见证版本
	WitnessProgram string `json:"witness_program"` // 见证程序
}

// 经过bitcoin core的finalizepsbt处理后的交易数据结构
type SignedTxDTO struct {
	Hex      string `json:"hex"`      // 交易十六进制
	Complete bool   `json:"complete"` // 是否完成
}
