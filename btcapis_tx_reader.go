package btcapis

import (
	"context"
	"errors"

	"github.com/crazycloudcc/btcapis/src/types"
)

// TxReader 提供交易相关的查询, 如交易原始数据、交易.
type TxReader interface {
	// GetRawTx 返回交易原始数据.
	GetRawTx(ctx context.Context, txid string) ([]byte, error)
	// GetTx 返回交易.
	GetTx(ctx context.Context, txid string) (*types.Tx, error)
}

func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	if c.bitcoindrpcClient != nil {
		return c.bitcoindrpcClient.GetRawTx(ctx, txid)
	}
	if c.mempoolspaceClient != nil {
		return c.mempoolspaceClient.GetRawTx(ctx, txid)
	}
	return nil, errors.New("btcapis: no client available")
}

// GetTx 返回交易.(优先使用bitcoindrpcClient, 其次使用mempoolspaceClient, 两边的数据格式不一致, 所以需要兼容)
func (c *Client) GetTx(ctx context.Context, txid string) (*types.Tx, error) {
	if c.bitcoindrpcClient != nil {
		return c.bitcoindrpcClient.GetTx(ctx, txid)
	}
	if c.mempoolspaceClient != nil {
		return c.mempoolspaceClient.GetTx(ctx, txid)
	}
	return nil, errors.New("btcapis: no client available")
}
