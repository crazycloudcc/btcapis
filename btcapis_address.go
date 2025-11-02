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

// 创建新钱包
func (c *Client) CreateNewWallet(ctx context.Context) (*types.WalletInfo, error) {
	return c.addressClient.GenerateNew()
}

// GetAddressBalance 返回地址的确认余额和未确认余额.
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (confirmed float64, mempool float64, err error) {
	return c.addressClient.GetAddressBalance(ctx, addr)
}

// GetAddressUTXOs 返回地址拥有的UTXO.
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.TxUTXO, error) {
	return c.addressClient.GetAddressUTXOs(ctx, addr)
}
