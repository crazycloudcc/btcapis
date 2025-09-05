package btcapis

import (
	"context"
	"fmt"

	"github.com/crazycloudcc/btcapis/types"
)

// 上传钱包+publickey, 用于后续组装PSBT等数据, 后续需要在postgres创建映射;
func (c *Client) ImportAddressAndPublickey(ctx context.Context, address string, publickey string) error {
	fmt.Printf("!!! (unsupport) import address: %s, publickey: %s\n", address, publickey)
	return nil
}

// GetAddressBalance 返回地址的确认余额和未确认余额.
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	return c.addressClient.GetAddressBalance(ctx, addr)
}

// GetAddressUTXOs 返回地址拥有的UTXO.
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.TxUTXO, error) {
	return c.addressClient.GetAddressUTXOs(ctx, addr)
}
