package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/address"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// GetAddressBalance 返回地址的确认余额和未确认余额.
func GetAddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	return mempoolapis.GetAddressBalance(ctx, addr)
}

// GetAddressUTXOs 返回地址拥有的UTXO.
func GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	return address.GetAddressUTXOs(ctx, addr)
}

// GetAddressScriptInfo 返回地址的锁定脚本信息.
func GetAddressScriptInfo(ctx context.Context, addr string) (*types.AddressScriptInfo, error) {
	return address.GetAddressScriptInfo(ctx, addr)
}

func GetAddressInfo(ctx context.Context, pkScript []byte) (*types.AddressInfo, error) {
	return address.GetAddressInfo(ctx, pkScript)
}
