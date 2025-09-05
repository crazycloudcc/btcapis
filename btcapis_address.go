package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/decoders"
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

// GetAddressScriptInfo 返回地址的锁定脚本信息.
func (c *Client) GetAddressScriptInfo(ctx context.Context, addr string) (*types.AddressScriptInfo, error) {
	return c.addressClient.GetAddressScriptInfo(ctx, addr)
}

// GetAddressInfo 返回地址的详细信息.
func (c *Client) GetAddressInfo(ctx context.Context, pkScript []byte) (*types.AddressInfo, error) {
	return c.addressClient.GetAddressInfo(ctx, pkScript)
}

// 通过地址转换为锁定脚本;
func (c *Client) AddressToPkScript(ctx context.Context, addr string) ([]byte, error) {
	return decoders.AddressToPkScript(addr)
}

// 通过地址转类型;
func (c *Client) AddressToType(ctx context.Context, addr string) (types.AddressType, error) {
	return decoders.AddressToType(addr)
}

// 通过脚本转类型;
func (c *Client) PKScriptToType(ctx context.Context, pkScript []byte) types.AddressType {
	return decoders.PKScriptToType(pkScript)
}
