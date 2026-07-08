package address

import (
	"context"
	"errors"

	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/pkg/logger"
	"github.com/crazycloudcc/btcapis/types"
)

var ErrAddressProviderUnavailable = errors.New("address provider unavailable")

// GetAddressBalanceSatsWithOptions electrumX → mempool.space
func (c *Client) GetAddressBalanceSatsWithOptions(ctx context.Context, addr string, opts types.AddressProviderOptions) (*types.AddressBalanceSats, error) {
	const method = "GetAddressBalanceSats"
	if opts.ElectrumX != nil {
		balance, err := opts.ElectrumX.GetBalanceSats(ctx, addr)
		if err == nil {
			logger.Debug("[%s] source=%s", method, opts.ElectrumX.Name())
			return balance, nil
		}
		logger.Warn("[%s] %s 获取失败: %v，回退 mempool.space", method, opts.ElectrumX.Name(), err)
	}
	if c.mempoolapisClient != nil {
		confirmed, unconfirmed, err := c.mempoolapisClient.AddressGetBalance(ctx, addr)
		if err == nil {
			total := confirmed + unconfirmed
			logger.Debug("[%s] source=mempoolSpace", method)
			return &types.AddressBalanceSats{
				Confirmed:   confirmed,
				Unconfirmed: unconfirmed,
				Total:       total,
				Address:     addr,
			}, nil
		}
		logger.Warn("[%s] mempool.space 获取失败: %v", method, err)
		return nil, err
	}
	return nil, ErrAddressProviderUnavailable
}

// GetAddressUTXOsWithOptions electrumX → mempool.space
func (c *Client) GetAddressUTXOsWithOptions(ctx context.Context, addr string, opts types.AddressProviderOptions) ([]types.AddressUTXO, error) {
	const method = "GetAddressUTXOs"
	if opts.ElectrumX != nil {
		utxos, err := opts.ElectrumX.GetUTXOs(ctx, addr)
		if err == nil {
			logger.Debug("[%s] source=%s", method, opts.ElectrumX.Name())
			return utxos, nil
		}
		logger.Warn("[%s] %s 获取失败: %v，回退 mempool.space", method, opts.ElectrumX.Name(), err)
	}
	utxos, err := utxosFromMempool(ctx, c.mempoolapisClient, addr)
	if err == nil {
		logger.Debug("[%s] source=mempoolSpace", method)
	}
	return utxos, err
}

func utxosFromMempool(ctx context.Context, mempool *mempoolapis.Client, addr string) ([]types.AddressUTXO, error) {
	if mempool == nil {
		return nil, ErrAddressProviderUnavailable
	}
	raw, err := mempool.AddressGetUTXOs(ctx, addr)
	if err != nil {
		return nil, err
	}
	out := make([]types.AddressUTXO, 0, len(raw))
	for _, u := range raw {
		out = append(out, types.AddressUTXO{
			Value:       u.Value,
			TxId:        u.Txid,
			Vout:        u.Vout,
			Height:      u.Status.BlockHeight,
			Confirmed:   u.Status.Confirmed,
			BlockHeight: u.Status.BlockHeight,
		})
	}
	return out, nil
}

// GetAddressHistoryTxs 地址交易历史（仅 electrumX）
func GetAddressHistoryTxs(ctx context.Context, addr string, electrum types.ElectrumXAddressProvider) ([]types.AddressHistoryTx, error) {
	if electrum == nil {
		return nil, ErrAddressProviderUnavailable
	}
	return electrum.GetHistoryTxs(ctx, addr)
}

// GetAddressUnconfirmedTxs 地址内存池交易（仅 electrumX）
func GetAddressUnconfirmedTxs(ctx context.Context, addr string, electrum types.ElectrumXAddressProvider) ([]types.AddressUnconfirmedTx, error) {
	if electrum == nil {
		return nil, ErrAddressProviderUnavailable
	}
	return electrum.GetUnconfirmedTxs(ctx, addr)
}

// BatchGetAddressBalanceSats 批量地址余额（仅 electrumX）
func BatchGetAddressBalanceSats(ctx context.Context, addrs []string, electrum types.ElectrumXAddressProvider) ([]types.AddressBalanceSats, error) {
	if electrum == nil {
		return nil, ErrAddressProviderUnavailable
	}
	return electrum.BatchGetBalanceSats(ctx, addrs)
}
