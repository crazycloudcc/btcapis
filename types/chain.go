package types

// // 区块头数据结构
// type ChainBlockHeader struct {
// 	Hash              string  `json:"hash"`              // 区块哈希
// 	Confirmations     int     `json:"confirmations"`     // 确认数
// 	Height            int     `json:"height"`            // 高度
// 	Version           int     `json:"version"`           // 版本
// 	VersionHex        string  `json:"versionHex"`        // 版本十六进制
// 	MerkleRoot        string  `json:"merkleroot"`        // 默克尔根
// 	Time              int     `json:"time"`              // 时间
// 	MedianTime        int     `json:"mediantime"`        // 中位时间
// 	Nonce             int     `json:"nonce"`             // 随机数
// 	Bits              string  `json:"bits"`              // 难度位
// 	Difficulty        float64 `json:"difficulty"`        // 难度
// 	Chainwork         string  `json:"chainwork"`         // 工作量
// 	NTx               int     `json:"nTx"`               // 交易数
// 	PreviousBlockHash string  `json:"previousblockhash"` // 上一个区块哈希
// 	NextBlockHash     string  `json:"nextblockhash"`     // 下一个区块哈希
// }

// 区块数据结构
type ChainBlock struct {
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

	// 区块链头数据不包含以下字段
	StrippedSize int      `json:"strippedsize"` // 剥离大小
	Size         int      `json:"size"`         // 大小
	Weight       int      `json:"weight"`       // 权重
	Tx           []string `json:"tx"`           // 交易
}
