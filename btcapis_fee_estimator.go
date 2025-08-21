package btcapis

import (
	"context"
	"errors"
)

// FeeEstimator 提供手续费估计功能.
type FeeEstimator interface {
	// EstimateFeeRate 估计手续费.
	EstimateFeeRate(ctx context.Context, targetBlocks int) (satsPerVByte float64, err error)
}

func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error) {
	if c.bitcoindrpcClient != nil {
		return c.bitcoindrpcClient.EstimateFeeRate(ctx, targetBlocks)
	}
	return 0, errors.New("btcapis: no client available")
}
