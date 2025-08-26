package tx

import (
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/address"
)

type Client struct {
	bitcoindrpcClient *bitcoindrpc.Client
	mempoolapisClient *mempoolapis.Client
	addressClient     *address.Client
}

func New(bitcoindrpcClient *bitcoindrpc.Client, mempoolapisClient *mempoolapis.Client, addressClient *address.Client) *Client {
	return &Client{
		bitcoindrpcClient: bitcoindrpcClient,
		mempoolapisClient: mempoolapisClient,
		addressClient:     addressClient,
	}
}
