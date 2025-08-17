// Package types 内存池相关类型定义
package types

import (
	"time"
)

// MempoolInfo 表示内存池信息
type MempoolInfo struct {
	Size             int       `json:"size"`              // 交易数量
	Bytes            int64     `json:"bytes"`             // 总字节数
	Usage            int64     `json:"usage"`             // 内存使用量
	MaxMempool       int64     `json:"max_mempool"`       // 最大内存池大小
	MempoolMinFee    float64   `json:"mempool_min_fee"`   // 最小手续费率
	MinRelayFee      float64   `json:"min_relay_fee"`     // 最小中继手续费率
	UnbroadcastCount int       `json:"unbroadcast_count"` // 未广播交易数
	LastUpdated      time.Time `json:"last_updated"`
	Backend          string    `json:"backend"`
}

// MempoolEntry 表示内存池中的交易条目
type MempoolEntry struct {
	TxID            string    `json:"txid"`
	Size            int       `json:"size"`
	VSize           int       `json:"vsize"`
	Weight          int       `json:"weight"`
	Fee             int64     `json:"fee"`
	FeeRate         float64   `json:"fee_rate"`
	Time            time.Time `json:"time"`
	Height          int64     `json:"height"`
	DescendantCount int       `json:"descendant_count"`
	DescendantSize  int       `json:"descendant_size"`
	DescendantFees  int64     `json:"descendant_fees"`
	AncestorCount   int       `json:"ancestor_count"`
	AncestorSize    int       `json:"ancestor_size"`
	AncestorFees    int64     `json:"ancestor_fees"`
}

// MempoolStats 表示内存池统计信息
type MempoolStats struct {
	TotalTxCount   int       `json:"total_tx_count"`
	TotalVSize     int64     `json:"total_vsize"`
	TotalFee       int64     `json:"total_fee"`
	AverageFeeRate float64   `json:"average_fee_rate"`
	MedianFeeRate  float64   `json:"median_fee_rate"`
	MinFeeRate     float64   `json:"min_fee_rate"`
	MaxFeeRate     float64   `json:"max_fee_rate"`
	LastUpdated    time.Time `json:"last_updated"`
}
