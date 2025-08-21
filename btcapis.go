package btcapis

import (
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/types"
)

func Init(network string, rpc_url, rpc_user, rpc_pass string, timeout int) {
	types.SetCurrentNetwork(network)

	if rpc_url != "" {
		bitcoindrpc.Init(rpc_url, rpc_user, rpc_pass, timeout)
	}

	mempool_rpc_url := ""
	if network == "mainnet" {
		mempool_rpc_url = "https://mempool.space"
	} else if network == "testnet" {
		mempool_rpc_url = "https://mempool.space/signet"
	}

	if mempool_rpc_url != "" {
		mempoolapis.Init(mempool_rpc_url, timeout)
	}
}
