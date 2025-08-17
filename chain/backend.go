// 统一接口：读链、广播、手续费、UTXO、区块、mempool 视图
package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

type ChainReader interface {
	GetRawTransaction(ctx context.Context, txid string) ([]byte, error)
	GetBlockHash(ctx context.Context, height int64) (string, error)
	GetBlockHeader(ctx context.Context, hash string) ([]byte, error)
	GetBlock(ctx context.Context, hash string) ([]byte, error)
	GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error)
}

type Broadcaster interface {
	Broadcast(ctx context.Context, rawtx []byte) (txid string, err error)
}

type FeeEstimator interface {
	EstimateFeeRate(ctx context.Context, targetBlocks int) (satsPerVByte float64, err error)
}

type MempoolView interface {
	GetRawMempool(ctx context.Context) ([]string, error) // txids
	TxInMempool(ctx context.Context, txid string) (bool, error)
}

type TxProvider interface {
	// 优先实现 raw，便于统一解析；若后端无法给 raw，可直接返回解析好的 Tx。
	GetRawTransaction(ctx context.Context, txid string) ([]byte, error) // 原有若已存在，保留
	GetTx(ctx context.Context, txid string) (*types.Tx, error)          // 新增：可选快速通道（例如 mempool JSON）
}

type Backend interface {
	ChainReader
	Broadcaster
	FeeEstimator
	MempoolView
	TxProvider
	Capabilities(ctx context.Context) (Capabilities, error)
}

type Capabilities struct {
	HasMempool     bool
	HasFeeEstimate bool
	Network        types.Network // mainnet/testnet/signet/regtest
}
