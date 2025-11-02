package address

import (
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/electrumx"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
)

type Client struct {
	bitcoindrpcClient *bitcoindrpc.Client
	mempoolapisClient *mempoolapis.Client
	electrumxClient   *electrumx.Client
}

func New(bitcoindrpcClient *bitcoindrpc.Client, mempoolapisClient *mempoolapis.Client, electrumxClient *electrumx.Client) *Client {
	return &Client{
		bitcoindrpcClient: bitcoindrpcClient,
		mempoolapisClient: mempoolapisClient,
		electrumxClient:   electrumxClient,
	}
}
