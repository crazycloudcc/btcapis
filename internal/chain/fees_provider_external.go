package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// ExternalFeesProvider 将 types.ChainFeesEstimator 适配为内置 FeesProvider
type ExternalFeesProvider struct {
	Estimator types.ChainFeesEstimator
}

func (p ExternalFeesProvider) Name() string {
	if p.Estimator == nil {
		return "external"
	}
	return p.Estimator.Name()
}

func (p ExternalFeesProvider) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	if p.Estimator == nil {
		return nil, ErrProviderUnavailable
	}
	fees, err := p.Estimator.EstimateChainFees(ctx, targetBlocks)
	if err != nil {
		return nil, err
	}
	return NormalizeChainFees(fees), nil
}

func electrumFeesProvider(builtin *ElectrumXFeesProvider, override types.ChainFeesEstimator) FeesProvider {
	if override != nil {
		return ExternalFeesProvider{Estimator: override}
	}
	if builtin != nil && builtin.Client != nil {
		return builtin
	}
	return nil
}
