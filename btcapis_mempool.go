package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/address"
	"github.com/crazycloudcc/btcapis/internal/chain"
	"github.com/crazycloudcc/btcapis/types"
)

// GetMempoolStats 获取 mempool.space 内存池统计
func (c *Client) GetMempoolStats(ctx context.Context) (*types.MempoolStats, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetMempoolStats(ctx)
}

// GetMempoolTxStatus 获取交易确认状态
func (c *Client) GetMempoolTxStatus(ctx context.Context, txid string) (*types.MempoolTxStatus, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetMempoolTxStatus(ctx, txid)
}

// GetMempoolFeesRecommend 从 mempool.space 获取多档推荐费率
func (c *Client) GetMempoolFeesRecommend(ctx context.Context) (*types.FeesRecommend, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetMempoolFeesRecommend(ctx)
}

// GetAddressBalanceSatsWithOptions 地址余额：electrumX → mempool.space
func (c *Client) GetAddressBalanceSatsWithOptions(ctx context.Context, addr string, opts types.AddressProviderOptions) (*types.AddressBalanceSats, error) {
	if c == nil || c.addressClient == nil {
		return nil, address.ErrAddressProviderUnavailable
	}
	return c.addressClient.GetAddressBalanceSatsWithOptions(ctx, addr, opts)
}

// GetAddressUTXOsWithOptions 地址 UTXO：electrumX → mempool.space
func (c *Client) GetAddressUTXOsWithOptions(ctx context.Context, addr string, opts types.AddressProviderOptions) ([]types.AddressUTXO, error) {
	if c == nil || c.addressClient == nil {
		return nil, address.ErrAddressProviderUnavailable
	}
	return c.addressClient.GetAddressUTXOsWithOptions(ctx, addr, opts)
}

// GetAddressHistoryTxs 地址交易历史（仅 electrumX）
func (c *Client) GetAddressHistoryTxs(ctx context.Context, addr string, electrum types.ElectrumXAddressProvider) ([]types.AddressHistoryTx, error) {
	return address.GetAddressHistoryTxs(ctx, addr, electrum)
}

// GetAddressUnconfirmedTxs 地址内存池交易（仅 electrumX）
func (c *Client) GetAddressUnconfirmedTxs(ctx context.Context, addr string, electrum types.ElectrumXAddressProvider) ([]types.AddressUnconfirmedTx, error) {
	return address.GetAddressUnconfirmedTxs(ctx, addr, electrum)
}

// BatchGetAddressBalanceSats 批量地址余额（仅 electrumX）
func (c *Client) BatchGetAddressBalanceSats(ctx context.Context, addrs []string, electrum types.ElectrumXAddressProvider) ([]types.AddressBalanceSats, error) {
	return address.BatchGetAddressBalanceSats(ctx, addrs, electrum)
}
