package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
)

type TestClient struct {
	bitcoindrpcClient *bitcoindrpc.Client
	mempoolapisClient *mempoolapis.Client
}

// 获取节点网络信息
func (c *TestClient) GetNetworkInfo(ctx context.Context) (*bitcoindrpc.NetworkInfoDTO, error) {
	return c.bitcoindrpcClient.GetNetworkInfo(ctx)
}

// 获取链信息
func (c *TestClient) GetBlockChainInfo(ctx context.Context) (*bitcoindrpc.ChainInfoDTO, error) {
	return c.bitcoindrpcClient.GetBlockChainInfo(ctx)
}

// 获取区块统计信息
func (c *TestClient) GetBlockStats(ctx context.Context, height int64) (*bitcoindrpc.BlockStatsDTO, error) {
	return c.bitcoindrpcClient.GetBlockStats(ctx, height)
}

// 获取链顶信息
func (c *TestClient) GetChainTips(ctx context.Context) ([]bitcoindrpc.ChainTipDTO, error) {
	return c.bitcoindrpcClient.GetChainTips(ctx)
}

// GetBlockCount
func (c *TestClient) GetBlockCount(ctx context.Context) (int64, error) {
	return c.bitcoindrpcClient.GetBlockCount(ctx)
}
