package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// // AddressReader 提供地址相关的查询, 如余额和UTXO.
// type AddressReader interface {
// 	AddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error)
// 	AddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error)
// }

// GetAddressBalance 返回地址的确认余额和未确认余额.
func GetAddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	return mempoolapis.GetAddressBalance(ctx, addr)
}

// GetAddressUTXOs 返回地址拥有的UTXO.
func GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	return mempoolapis.GetAddressUTXOs(ctx, addr)
}
