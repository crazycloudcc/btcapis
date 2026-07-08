package chain

import (
	"context"
)

// EstimateFeeRate 兼容旧接口，内部转调统一费率估算
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, float64, error) {
	fees, err := c.EstimateChainFeesDefault(ctx, int64(targetBlocks), nil)
	if err != nil {
		return 0, 0, err
	}
	return fees.High, fees.Low, nil
}
