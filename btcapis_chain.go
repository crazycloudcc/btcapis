package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/extensions/okx"
	"github.com/crazycloudcc/btcapis/extensions/unisat"
	"github.com/crazycloudcc/btcapis/internal/chain"
	"github.com/crazycloudcc/btcapis/types"
)

var (
	ErrBlockNotFound    = chain.ErrBlockNotFound
	ErrInvalidBlockData = chain.ErrInvalidBlockData
)

// DefaultFeesProviderOptions 默认费率 Provider 选项
type DefaultFeesProviderOptions struct {
	// ElectrumX 可选自定义实现（如 TCP ElectrumX），优先于内置 HTTP ElectrumX
	ElectrumX types.ChainFeesEstimator
}

// EstimateChainFees 使用默认节点 Provider 链估算手续费（mempool.space → electrumX → bitcoin core）
func (c *Client) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	return c.EstimateChainFeesWithOptions(ctx, targetBlocks, DefaultFeesProviderOptions{})
}

// EstimateChainFeesWithOptions 使用默认 Provider 链估算手续费，支持自定义 ElectrumX
func (c *Client) EstimateChainFeesWithOptions(ctx context.Context, targetBlocks int64, opts DefaultFeesProviderOptions) (*types.ChainFees, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrProviderUnavailable
	}
	return c.chainClient.EstimateChainFeesDefault(ctx, targetBlocks, opts.ElectrumX)
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
	// ElectrumX 可选自定义实现（如 TCP ElectrumX），优先于内置 HTTP ElectrumX
	ElectrumX types.ChainFeesEstimator
}

// EstimateChainFeesForAPI 使用 apis 完整降级链估算手续费
// mempool.space → unisat → okx → electrumX → bitcoin core
func (c *Client) EstimateChainFeesForAPI(ctx context.Context, targetBlocks int64, opts APIFeesProviderOptions) (*types.ChainFees, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrProviderUnavailable
	}
	clients := c.chainClient.NewFeesClients()
	providers := chain.APIFeesProviders(clients.Mempool, opts.Unisat, opts.OKX, clients.Electrum, clients.Bitcoind, opts.ElectrumX)
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

// GetNetworkInfo 获取节点网络信息
func (c *Client) GetNetworkInfo(ctx context.Context) (*types.ChainNetworkInfo, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetNetworkInfo(ctx)
}

// GetBlockStats 获取区块统计
func (c *Client) GetBlockStats(ctx context.Context) (*types.ChainStats, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetBlockStats(ctx)
}

// GetChainTips 获取链分叉信息
func (c *Client) GetChainTips(ctx context.Context) ([]types.ChainTip, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetChainTips(ctx)
}

// GetBlockChainInfo 获取链状态
func (c *Client) GetBlockChainInfo(ctx context.Context) (*types.BlockChainInfo, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetBlockChainInfo(ctx)
}

// GetBlockHeader 使用区块 hash 查询区块头
func (c *Client) GetBlockHeader(ctx context.Context, blockHash string) (*types.BlockHeader, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetBlockHeader(ctx, blockHash)
}

// GetBlock 使用区块 hash 查询区块（verbosity=1）
func (c *Client) GetBlock(ctx context.Context, blockHash string) (*types.Block, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetBlock(ctx, blockHash)
}

// GetBlockTransactions 使用区块 hash 查询紧凑有序交易（verbosity=2）。
func (c *Client) GetBlockTransactions(ctx context.Context, blockHash string) (*types.BlockTransactions, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetBlockTransactions(ctx, blockHash)
}

// GetBlockHashByHeight 使用区块高度查询 hash
func (c *Client) GetBlockHashByHeight(ctx context.Context, height int64) (string, error) {
	if c == nil || c.chainClient == nil {
		return "", chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetBlockHashByHeight(ctx, height)
}
