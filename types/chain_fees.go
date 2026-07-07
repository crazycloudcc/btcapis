package types

// ChainFees 链上手续费估算（统一单位 sat/vB，保留 2 位小数）
type ChainFees struct {
	High    float64 `json:"high"`    // 高优先级费率 (sat/vB)
	Medium  float64 `json:"medium"`  // 中优先级费率 (sat/vB)
	Low     float64 `json:"low"`     // 低优先级费率 (sat/vB)
	FeeRate float64 `json:"feerate"` // 目标确认档费率 (sat/vB)
	Blocks  int64   `json:"blocks"`  // 目标确认区块数
}

// RecommendedFees 多档推荐费率（sat/vB）
type RecommendedFees struct {
	FastestFee  float64
	HalfHourFee float64
	HourFee     float64
	EconomyFee  float64
	MinimumFee  float64
}
