package btcapis

import (
	"context"
	"errors"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
)

// // FeeEstimator 提供手续费估计功能.
// type FeeEstimator interface {
// 	// EstimateFeeRate 估计手续费.
// 	EstimateFeeRate(ctx context.Context, targetBlocks int) (satsPerVByte float64, err error)
// }

// EstimateFeeRate 估计手续费.
func EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error) {
	if bitcoindrpc.IsInited() {
		return bitcoindrpc.EstimateFeeRate(ctx, targetBlocks)
	}

	// if mempoolapis.IsInited() {
	// 	return mempoolapis.EstimateFeeRate(ctx, targetBlocks)
	// }

	return 0, errors.New("btcapis: no client available")
}
