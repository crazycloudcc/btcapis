package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/types"
)

// GetAddressBalance 返回地址的确认余额和未确认余额.
func (c *Client) GetAddressBalanceSats(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	return c.addressClient.GetAddressBalance(ctx, addr)
}

// GetAddressUTXOs 返回地址拥有的UTXO.
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	return c.addressClient.GetAddressUTXOs(ctx, addr)
}

// GetAddressScriptInfo 返回地址的锁定脚本信息.
func (c *Client) GetAddressScriptInfo(ctx context.Context, addr string) (*types.AddressScriptInfo, error) {
	return c.addressClient.GetAddressScriptInfo(ctx, addr)
}

func (c *Client) GetAddressInfo(ctx context.Context, pkScript []byte) (*types.AddressInfo, error) {
	return c.addressClient.GetAddressInfo(ctx, pkScript)
}
