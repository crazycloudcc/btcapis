package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// TxReader 提供交易相关的查询, 如交易原始数据、交易.
type TxReader interface {
	// GetRawTransaction 返回交易原始数据.
	GetRawTransaction(ctx context.Context, txid string) ([]byte, error)
	// GetTx 返回交易.
	GetTx(ctx context.Context, txid string) (*types.Tx, error)
}
