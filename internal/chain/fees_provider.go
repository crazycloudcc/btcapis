package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// FeesProvider 链上手续费数据源
type FeesProvider interface {
	Name() string
	EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error)
}

// RecommendedFeesProvider 返回多档推荐费率的数据源（如 unisat / okx）
type RecommendedFeesProvider interface {
	Name() string
	GetRecommendedFees(ctx context.Context) (types.RecommendedFees, error)
}

// RecommendedFeesAdapter 将 RecommendedFeesProvider 适配为 FeesProvider
type RecommendedFeesAdapter struct {
	Provider RecommendedFeesProvider
}

func (a RecommendedFeesAdapter) Name() string {
	if a.Provider == nil {
		return "recommended"
	}
	return a.Provider.Name()
}

func (a RecommendedFeesAdapter) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	if a.Provider == nil {
		return nil, ErrProviderUnavailable
	}
	recommended, err := a.Provider.GetRecommendedFees(ctx)
	if err != nil {
		return nil, err
	}
	return ChainFeesFromRecommended(targetBlocks, recommended), nil
}
