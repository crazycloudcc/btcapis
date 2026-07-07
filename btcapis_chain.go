package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/extensions/okx"
	"github.com/crazycloudcc/btcapis/extensions/unisat"
	"github.com/crazycloudcc/btcapis/internal/chain"
	"github.com/crazycloudcc/btcapis/types"
)

// EstimateChainFees 使用默认节点 Provider 链估算手续费（mempool.space → electrumX → bitcoin core）
func (c *Client) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrProviderUnavailable
	}
	return c.chainClient.EstimateChainFeesDefault(ctx, targetBlocks)
}

// EstimateChainFeesWithProviders 自定义 Provider 顺序
func (c *Client) EstimateChainFeesWithProviders(ctx context.Context, targetBlocks int64, providers []chain.FeesProvider) (*types.ChainFees, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrProviderUnavailable
	}
	return c.chainClient.EstimateChainFeesWithProviders(ctx, targetBlocks, providers)
}

// APIFeesProviderOptions apis 服务费率 Provider 配置
type APIFeesProviderOptions struct {
	Unisat *unisat.Client
	OKX    *okx.Client
}

// EstimateChainFeesForAPI 使用 apis 完整降级链估算手续费
// mempool.space → unisat → okx → electrumX → bitcoin core
func (c *Client) EstimateChainFeesForAPI(ctx context.Context, targetBlocks int64, opts APIFeesProviderOptions) (*types.ChainFees, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrProviderUnavailable
	}
	clients := c.chainClient.NewFeesClients()
	providers := chain.APIFeesProviders(clients.Mempool, opts.Unisat, opts.OKX, clients.Electrum, clients.Bitcoind)
	return c.chainClient.EstimateChainFeesWithProviders(ctx, targetBlocks, providers)
}

// EstimateFeeRate 兼容旧接口：返回高优先级费率与低优先级费率（sat/vB）
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, float64, error) {
	fees, err := c.EstimateChainFees(ctx, int64(targetBlocks))
	if err != nil {
		return 0, 0, err
	}
	return fees.High, fees.Low, nil
}
