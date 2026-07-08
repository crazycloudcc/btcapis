package types

import "context"

// ChainFeesEstimator 链上手续费数据源（供外部注入自定义实现，如 TCP ElectrumX）
type ChainFeesEstimator interface {
	Name() string
	EstimateChainFees(ctx context.Context, targetBlocks int64) (*ChainFees, error)
}
