// Package types 手续费相关类型定义
package types

import (
	"time"
)

// FeeEstimate 表示手续费估算结果
type FeeEstimate struct {
	TargetBlocks int       `json:"target_blocks"`
	FeeRate      float64   `json:"fee_rate"`     // sat/vB
	FeeRateBTC   float64   `json:"fee_rate_btc"` // BTC/kB
	Confidence   float64   `json:"confidence"`   // 置信度 0-1
	EstimatedAt  time.Time `json:"estimated_at"`
	Backend      string    `json:"backend"` // 数据来源后端
}

// FeeHistory 表示手续费历史数据
type FeeHistory struct {
	TargetBlocks []int       `json:"target_blocks"`
	FeeRates     []float64   `json:"fee_rates"`
	Timestamps   []time.Time `json:"timestamps"`
	Backend      string      `json:"backend"`
}

// MempoolFee 表示内存池中的手续费信息
type MempoolFee struct {
	Count    int     `json:"count"`     // 交易数量
	VSize    int64   `json:"vsize"`     // 虚拟大小
	FeeRate  float64 `json:"fee_rate"`  // 手续费率
	TotalFee int64   `json:"total_fee"` // 总手续费
}

// FeeTarget 表示手续费目标
type FeeTarget struct {
	Blocks     int     `json:"blocks"`     // 目标确认区块数
	FeeRate    float64 `json:"fee_rate"`   // 对应手续费率
	Confidence float64 `json:"confidence"` // 置信度
}
