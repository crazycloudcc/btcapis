package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// FeesClients 节点费率客户端集合
type FeesClients struct {
	Mempool  *MempoolFeesProvider
	Electrum *ElectrumXFeesProvider
	Bitcoind *BitcoindFeesProvider
}

// NewFeesClients 从 chain 客户端构建内置 Provider
func (c *Client) NewFeesClients() FeesClients {
	return FeesClients{
		Mempool:  &MempoolFeesProvider{Client: c.mempoolapisClient},
		Electrum: &ElectrumXFeesProvider{Client: c.electrumxClient},
		Bitcoind: &BitcoindFeesProvider{Client: c.bitcoindrpcClient},
	}
}

// EstimateChainFeesDefault 使用默认 Provider 链估算费率
func (c *Client) EstimateChainFeesDefault(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	clients := c.NewFeesClients()
	return EstimateChainFees(ctx, targetBlocks, DefaultFeesProviders(clients.Mempool, clients.Electrum, clients.Bitcoind))
}

// EstimateChainFeesWithProviders 自定义 Provider 顺序
func (c *Client) EstimateChainFeesWithProviders(ctx context.Context, targetBlocks int64, providers []FeesProvider) (*types.ChainFees, error) {
	return EstimateChainFees(ctx, targetBlocks, providers)
}
