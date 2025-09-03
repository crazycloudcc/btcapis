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

func (c *TestClient) GetNetworkInfo(ctx context.Context) (*bitcoindrpc.NetworkInfoDTO, error) {
	res, err := c.bitcoindrpcClient.GetNetworkInfo(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
