package types

// MempoolStats mempool.space 内存池统计
type MempoolStats struct {
	Count        int64       `json:"count"`
	Vsize        int64       `json:"vsize"`
	TotalFee     int64       `json:"total_fee"`
	FeeHistogram [][]float64 `json:"fee_histogram"`
}

// MempoolTxStatus 交易确认状态
type MempoolTxStatus struct {
	TxID          string  `json:"txid"`
	Confirmed     bool    `json:"confirmed"`
	BlockHeight   int64   `json:"block_height,omitempty"`
	BlockHash     string  `json:"block_hash,omitempty"`
	BlockTime     int64   `json:"block_time,omitempty"`
	InMempool     bool    `json:"in_mempool"`
	Replaceable   bool    `json:"replaceable"`
	Confirmations int64   `json:"confirmations"`
	Fee           int64   `json:"fee,omitempty"`
	FeeRateSatVB  float64 `json:"fee_rate_sat_vb,omitempty"`
}

// FeesRecommend 多档推荐费率
type FeesRecommend struct {
	Blocks1 ChainFees `json:"blocks_1"`
	Blocks3 ChainFees `json:"blocks_3"`
	Blocks6 ChainFees `json:"blocks_6"`
	Source  string    `json:"source"`
}
