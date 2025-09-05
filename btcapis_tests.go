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
	res, err := c.bitcoindrpcClient.GetNetworkInfo(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取链信息
func (c *TestClient) GetBlockChainInfo(ctx context.Context) (*bitcoindrpc.ChainInfoDTO, error) {
	res, err := c.bitcoindrpcClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取区块统计信息
func (c *TestClient) GetBlockStats(ctx context.Context, height int64) (*bitcoindrpc.BlockStatsDTO, error) {
	res, err := c.bitcoindrpcClient.GetBlockStats(ctx, height)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取链顶信息
func (c *TestClient) GetChainTips(ctx context.Context) ([]bitcoindrpc.ChainTipDTO, error) {
	res, err := c.bitcoindrpcClient.GetChainTips(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}
