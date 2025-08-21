// 统一接口：读链、广播、手续费、UTXO、区块、mempool 视图
package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// ChainReader 提供区块相关的查询, 如区块哈希、区块头、区块等.
type ChainReader interface {
	// GetRawTransaction 返回交易原始数据.
	GetRawTransaction(ctx context.Context, txid string) ([]byte, error)
	// GetBlockHash 返回区块哈希.
	GetBlockHash(ctx context.Context, height int64) (string, error)
	// GetBlockHeader 返回区块头.
	GetBlockHeader(ctx context.Context, hash string) ([]byte, error)
	// GetBlock 返回区块.
	GetBlock(ctx context.Context, hash string) ([]byte, error)
	// GetUTXO 返回UTXO.
	GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error)
}

// Broadcaster 提供交易广播功能.
type Broadcaster interface {
	// Broadcast 广播交易.
	Broadcast(ctx context.Context, rawtx []byte) (txid string, err error)
}

// FeeEstimator 提供手续费估计功能.
type FeeEstimator interface {
	// EstimateFeeRate 估计手续费.
	EstimateFeeRate(ctx context.Context, targetBlocks int) (satsPerVByte float64, err error)
}

// MempoolView 提供mempool相关的查询, 如mempool中的交易.
type MempoolView interface {
	// GetRawMempool 返回mempool中的交易.
	GetRawMempool(ctx context.Context) ([]string, error) // txids
	// TxInMempool 返回交易是否在mempool中.
	TxInMempool(ctx context.Context, txid string) (bool, error)
}

// Backend 提供链相关的查询, 如区块、交易、手续费、mempool等.
type Backend interface {
	ChainReader
	Broadcaster
	FeeEstimator
	MempoolView
	TxReader
	Capabilities(ctx context.Context) (Capabilities, error)
}
