package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/src/types"
)

// AddressReader 提供地址相关的查询, 如余额和UTXO.
type AddressReader interface {
	// AddressBalance 返回地址的确认余额和未确认余额.
	AddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error)
	// AddressUTXOs 返回地址拥有的UTXO.
	AddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error)
}

func (c *Client) AddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	return c.mempoolspaceClient.AddressBalance(ctx, addr)
}

func (c *Client) AddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	return c.mempoolspaceClient.AddressUTXOs(ctx, addr)
}
