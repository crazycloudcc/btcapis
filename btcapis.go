package btcapis

import (
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/address"
	"github.com/crazycloudcc/btcapis/internal/chain"
	"github.com/crazycloudcc/btcapis/internal/tx"
	"github.com/crazycloudcc/btcapis/internal/types"
)

type Client struct {
	// bitcoindrpcClient *bitcoindrpc.Client // bitcoindrpc接口调用集合.
	// mempoolapisClient *mempoolapis.Client // mempool.space接口调用集合.
	addressClient *address.Client // 钱包地址操作
	txClient      *tx.Client      // 交易操作
	chainClient   *chain.Client   // 链操作 - 无法归类到钱包和交易类的其他链上操作
}

func New(network string, rpc_url, rpc_user, rpc_pass string, timeout int) (*Client, *TestClient) {
	types.SetCurrentNetwork(network)

	var bitcoindrpcClient *bitcoindrpc.Client = nil
	var mempoolapisClient *mempoolapis.Client = nil

	if rpc_url != "" {
		bitcoindrpcClient = bitcoindrpc.New(rpc_url, rpc_user, rpc_pass, timeout)
	}

	mempool_rpc_url := ""
	if network == "mainnet" {
		mempool_rpc_url = "https://mempool.space"
	} else if network == "signet" {
		mempool_rpc_url = "https://mempool.space/signet"
	} else if network == "testnet" {
		mempool_rpc_url = "https://mempool.space/testnet"
	}

	if mempool_rpc_url != "" {
		mempoolapisClient = mempoolapis.New(mempool_rpc_url, timeout)
	}

	client := &Client{
		addressClient: address.New(bitcoindrpcClient, mempoolapisClient),
		txClient:      tx.New(bitcoindrpcClient, mempoolapisClient),
		chainClient:   chain.New(bitcoindrpcClient, mempoolapisClient),
	}

	testClient := &TestClient{
		bitcoindrpcClient: bitcoindrpcClient,
		mempoolapisClient: mempoolapisClient,
	}

	return client, testClient
}
