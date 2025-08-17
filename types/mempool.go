// Package types 内存池相关类型定义
// 包含比特币内存池（mempool）中待确认交易的信息、统计数据和状态管理
package types

import (
	"time"
)

// MempoolInfo 表示内存池信息
// 包含内存池的整体状态、容量限制和手续费策略等基本信息
type MempoolInfo struct {
	// Size 交易数量
	// 当前内存池中待确认交易的总数量
	Size int `json:"size"` // 交易数量
	// Bytes 总字节数
	// 内存池中所有交易占用的总字节数
	Bytes int64 `json:"bytes"` // 总字节数
	// Usage 内存使用量
	// 内存池实际使用的内存大小（字节）
	Usage int64 `json:"usage"` // 内存使用量
	// MaxMempool 最大内存池大小
	// 内存池允许的最大大小（字节），超过此限制的交易会被拒绝
	MaxMempool int64 `json:"max_mempool"` // 最大内存池大小
	// MempoolMinFee 最小手续费率
	// 内存池接受交易的最小手续费率（BTC/kB）
	MempoolMinFee float64 `json:"mempool_min_fee"` // 最小手续费率
	// MinRelayFee 最小中继手续费率
	// 节点接受并转发交易的最小手续费率（BTC/kB）
	MinRelayFee float64 `json:"min_relay_fee"` // 最小中继手续费率
	// UnbroadcastCount 未广播交易数
	// 尚未广播到网络的其他节点的交易数量
	UnbroadcastCount int `json:"unbroadcast_count"` // 未广播交易数
	// LastUpdated 最后更新时间
	// 内存池信息最后更新的时间戳
	LastUpdated time.Time `json:"last_updated"`
	// Backend 后端存储类型
	// 内存池使用的存储后端，如 "memory" 或 "leveldb"
	Backend string `json:"backend"`
}

// MempoolEntry 表示内存池中的交易条目
// 包含单个待确认交易的详细信息、手续费和依赖关系
type MempoolEntry struct {
	// TxID 交易ID
	// 交易的唯一标识符（64字符的十六进制字符串）
	TxID string `json:"txid"`
	// Size 交易大小
	// 交易的原始大小（字节）
	Size int `json:"size"`
	// VSize 虚拟大小
	// 交易的虚拟大小（字节），用于手续费计算
	VSize int `json:"vsize"`
	// Weight 交易权重
	// 交易的权重单位，用于隔离见证手续费计算
	Weight int `json:"weight"`
	// Fee 手续费
	// 交易支付的手续费（聪）
	Fee int64 `json:"fee"`
	// FeeRate 手续费率
	// 交易的手续费率（聪/字节）
	FeeRate float64 `json:"fee_rate"`
	// Time 进入内存池时间
	// 交易首次进入内存池的时间戳
	Time time.Time `json:"time"`
	// Height 区块高度
	// 交易被包含的区块高度，-1表示尚未被包含
	Height int64 `json:"height"`
	// DescendantCount 后代交易数量
	// 依赖于此交易的未确认交易数量
	DescendantCount int `json:"descendant_count"`
	// DescendantSize 后代交易大小
	// 所有后代交易的总大小（字节）
	DescendantSize int `json:"descendant_size"`
	// DescendantFees 后代交易手续费
	// 所有后代交易的总手续费（聪）
	DescendantFees int64 `json:"descendant_fees"`
	// AncestorCount 祖先交易数量
	// 此交易依赖的未确认交易数量
	AncestorCount int `json:"ancestor_count"`
	// AncestorSize 祖先交易大小
	// 所有祖先交易的总大小（字节）
	AncestorSize int `json:"ancestor_size"`
	// AncestorFees 祖先交易手续费
	// 所有祖先交易的总手续费（聪）
	AncestorFees int64 `json:"ancestor_fees"`
}

// MempoolStats 表示内存池统计信息
// 包含内存池中交易的汇总统计数据和手续费分布信息
type MempoolStats struct {
	// TotalTxCount 总交易数量
	// 内存池中所有交易的总数
	TotalTxCount int `json:"total_tx_count"`
	// TotalVSize 总虚拟大小
	// 所有交易的虚拟大小总和（字节）
	TotalVSize int64 `json:"total_vsize"`
	// TotalFee 总手续费
	// 所有交易支付的手续费总和（聪）
	TotalFee int64 `json:"total_fee"`
	// AverageFeeRate 平均手续费率
	// 所有交易的平均手续费率（聪/字节）
	AverageFeeRate float64 `json:"average_fee_rate"`
	// MedianFeeRate 中位数手续费率
	// 所有交易手续费率的中位数（聪/字节）
	MedianFeeRate float64 `json:"median_fee_rate"`
	// MinFeeRate 最小手续费率
	// 内存池中交易的最低手续费率（聪/字节）
	MinFeeRate float64 `json:"min_fee_rate"`
	// MaxFeeRate 最大手续费率
	// 内存池中交易的最高手续费率（聪/字节）
	MaxFeeRate float64 `json:"max_fee_rate"`
	// LastUpdated 最后更新时间
	// 统计信息最后更新的时间戳
	LastUpdated time.Time `json:"last_updated"`
}
