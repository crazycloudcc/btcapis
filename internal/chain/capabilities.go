// 能力探测/特性位（是否支持rawmempool、feeEstimates等）
package chain

import "github.com/crazycloudcc/btcapis/types"

// Capabilities 提供链相关的查询, 如区块、交易、手续费、mempool等.
type Capabilities struct {
	HasMempool     bool
	HasFeeEstimate bool
	Network        types.Network // mainnet/testnet/signet/regtest
}
