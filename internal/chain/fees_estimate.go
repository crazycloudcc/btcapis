package chain

import (
	"context"
	"errors"

	"github.com/crazycloudcc/btcapis/pkg/logger"
	"github.com/crazycloudcc/btcapis/types"
)

var (
	// ErrProviderUnavailable Provider 未配置或不可用
	ErrProviderUnavailable = errors.New("fees provider unavailable")
	// ErrAllProvidersFailed 所有 Provider 均失败
	ErrAllProvidersFailed = errors.New("all fees providers failed")
)

// EstimateChainFees 按 Provider 顺序估算链上手续费，失败时记录日志并回退
func EstimateChainFees(ctx context.Context, targetBlocks int64, providers []FeesProvider) (*types.ChainFees, error) {
	const method = "EstimateChainFees"
	if len(providers) == 0 {
		logger.Warn("[%s] 未配置任何 Provider", method)
		return nil, ErrAllProvidersFailed
	}

	for i, provider := range providers {
		if provider == nil {
			logger.Warn("[%s] Provider[%d] 为空，跳过", method, i)
			continue
		}
		fees, err := provider.EstimateChainFees(ctx, targetBlocks)
		if err == nil {
			logger.Debug("[%s] source=%s blocks=%d", method, provider.Name(), targetBlocks)
			return NormalizeChainFees(fees), nil
		}
		if i < len(providers)-1 {
			logger.Warn("[%s] %s 获取失败: %v，回退下一数据源", method, provider.Name(), err)
		} else {
			logger.Warn("[%s] %s 获取失败: %v", method, provider.Name(), err)
		}
	}

	logger.Warn("[%s] 所有数据源均不可用", method)
	return nil, ErrAllProvidersFailed
}

// DefaultFeesProviders 默认节点 Provider：mempool.space → electrumX → bitcoin core
func DefaultFeesProviders(mempool *MempoolFeesProvider, electrum *ElectrumXFeesProvider, bitcoind *BitcoindFeesProvider) []FeesProvider {
	providers := make([]FeesProvider, 0, 3)
	if mempool != nil && mempool.Client != nil {
		providers = append(providers, mempool)
	}
	if electrum != nil && electrum.Client != nil {
		providers = append(providers, electrum)
	}
	if bitcoind != nil && bitcoind.Client != nil {
		providers = append(providers, bitcoind)
	}
	return providers
}

// APIFeesProviders apis 服务完整降级链：mempool → unisat → okx → electrumX → bitcoin core
func APIFeesProviders(
	mempool *MempoolFeesProvider,
	unisat RecommendedFeesProvider,
	okx RecommendedFeesProvider,
	electrum *ElectrumXFeesProvider,
	bitcoind *BitcoindFeesProvider,
) []FeesProvider {
	providers := make([]FeesProvider, 0, 5)
	if mempool != nil && mempool.Client != nil {
		providers = append(providers, mempool)
	}
	if unisat != nil {
		providers = append(providers, RecommendedFeesAdapter{Provider: unisat})
	}
	if okx != nil {
		providers = append(providers, RecommendedFeesAdapter{Provider: okx})
	}
	if electrum != nil && electrum.Client != nil {
		providers = append(providers, electrum)
	}
	if bitcoind != nil && bitcoind.Client != nil {
		providers = append(providers, bitcoind)
	}
	return providers
}
