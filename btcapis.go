package btcapis

import (
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/electrumx"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/address"
	"github.com/crazycloudcc/btcapis/internal/chain"
	"github.com/crazycloudcc/btcapis/internal/tx"
	"github.com/crazycloudcc/btcapis/types"
)

type Client struct {
	// bitcoindrpcClient *bitcoindrpc.Client // bitcoindrpc接口调用集合.
	// mempoolapisClient *mempoolapis.Client // mempool.space接口调用集合.
	// electrumxClient   *electrumx.Client   // electrumx接口调用集合.
	addressClient *address.Client // 钱包地址操作
	txClient      *tx.Client      // 交易操作
	chainClient   *chain.Client   // 链操作 - 无法归类到钱包和交易类的其他链上操作
}

var bitcoindrpcClient *bitcoindrpc.Client
var mempoolapisClient *mempoolapis.Client
var electrumxClient *electrumx.Client

func New(network string, rpc_url, rpc_user, rpc_pass string, timeout int) *Client {
	types.SetCurrentNetwork(network)

	client := &Client{}

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
	mempoolapisClient = mempoolapis.New(mempool_rpc_url, timeout)

	ex_rpc_url := ""
	if network == "mainnet" {
		ex_rpc_url = "http://localhost:50001"
	} else if network == "signet" {
		ex_rpc_url = "https://blockstream.info/electrum"
	} else if network == "testnet" {
		ex_rpc_url = "https://blockstream.info/electrum"
	}
	electrumxClient = electrumx.New(ex_rpc_url, timeout)

	client.addressClient = address.New(bitcoindrpcClient, mempoolapisClient, electrumxClient)
	client.txClient = tx.New(bitcoindrpcClient, mempoolapisClient, client.addressClient)
	client.chainClient = chain.New(bitcoindrpcClient, mempoolapisClient)

	return client
}

// NewWithElectrumX 创建包含ElectrumX支持的客户端
func NewWithElectrumX(network string, rpc_url, rpc_user, rpc_pass string, electrumx_url string, timeout int) *Client {
	types.SetCurrentNetwork(network)

	client := &Client{}

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

	if electrumx_url != "" {
		electrumxClient = electrumx.New(electrumx_url, timeout)
	}

	client.addressClient = address.New(bitcoindrpcClient, mempoolapisClient, electrumxClient)
	client.txClient = tx.New(bitcoindrpcClient, mempoolapisClient, client.addressClient)
	client.chainClient = chain.New(bitcoindrpcClient, mempoolapisClient)

	return client
}

func NewTestClient(client *Client) *TestClient {
	return &TestClient{
		bitcoindrpcClient: bitcoindrpcClient,
		mempoolapisClient: mempoolapisClient,
		electrumxClient:   electrumxClient,
	}
}
