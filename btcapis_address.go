package btcapis

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/crazycloudcc/btcapis/address"
	"github.com/crazycloudcc/btcapis/chain"
	"github.com/crazycloudcc/btcapis/types"
)

// 对外门面：使用地址 解析脚本信息
func GetAddress2ScriptInfo(addr string, params *chaincfg.Params) (*types.AddressScriptInfo, error) {
	return address.ParseAddress(addr, params)
}

// 对外门面：使用地址 查询余额
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (float64, float64, error) {
	confirmed, mempool, err := c.addressBalance(ctx, addr)
	if err != nil {
		return 0, 0, err
	}
	return satsToBTC(confirmed), satsToBTC(mempool), nil
}

// 对外门面：使用地址 查询UTXO
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	for _, b := range append(c.primaries, c.fallbacks...) {
		if ar, ok := b.(chain.AddressReader); ok {
			if u, err := ar.AddressUTXOs(ctx, addr); err == nil {
				return u, nil
			}
		}
	}
	return nil, chain.ErrBackendUnavailable
}

// 内部实现：使用地址 查询余额
func (c *Client) addressBalance(ctx context.Context, addr string) (int64, int64, error) {
	for _, b := range append(c.primaries, c.fallbacks...) {
		if ar, ok := b.(chain.AddressReader); ok {
			if confirmed, mempool, err := ar.AddressBalance(ctx, addr); err == nil {
				return confirmed, mempool, nil
			}
		}
	}
	return 0, 0, chain.ErrBackendUnavailable
}

func satsToBTC(v int64) float64 {
	f := float64(v) / 1e8
	return f
}
