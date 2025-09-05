package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// GetAddressBalance 返回地址的确认余额和未确认余额.
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	return c.addressClient.GetAddressBalance(ctx, addr)
}

// GetAddressUTXOs 返回地址拥有的UTXO.
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.TxUTXO, error) {
	return c.addressClient.GetAddressUTXOs(ctx, addr)
}
