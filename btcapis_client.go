package btcapis

import (
	"net/http"
	"time"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// Client, 包含所有的provider client, 便于外部统一调用
type Client struct {
	bitcoindrpcClient  *bitcoindrpc.Client
	mempoolspaceClient *mempoolapis.Client
	// TODO: 添加其他provider
}

func NewClient(network string, rpc_url, rpc_user, rpc_pass string) *Client {
	types.SetCurrentNetwork(network)

	client := &Client{}

	if rpc_url != "" {
		client.bitcoindrpcClient = bitcoindrpc.New(rpc_url, rpc_user, rpc_pass,
			bitcoindrpc.WithHTTPClient(&http.Client{
				Timeout: 10 * time.Second,
			}),
		)
	}

	mempool_rpc_url := ""
	if network == "mainnet" {
		mempool_rpc_url = "https://mempool.space"
	} else if network == "testnet" {
		mempool_rpc_url = "https://mempool.space/signet"
	}

	if mempool_rpc_url != "" {
		client.mempoolspaceClient = mempoolapis.New(mempool_rpc_url,
			mempoolapis.WithHTTPClient(&http.Client{
				Timeout: 10 * time.Second,
			}),
		)
	}

	return client
}
