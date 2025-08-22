package tx

import (
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
)

type Client struct {
	bitcoindrpcClient *bitcoindrpc.Client
	mempoolapisClient *mempoolapis.Client
}

func New(bitcoindrpcClient *bitcoindrpc.Client, mempoolapisClient *mempoolapis.Client) *Client {
	return &Client{
		bitcoindrpcClient: bitcoindrpcClient,
		mempoolapisClient: mempoolapisClient,
	}
}
