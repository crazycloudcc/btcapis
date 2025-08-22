package chain

import (
	"context"
)

func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, float64, error) {
	var fee1, fee2 float64 = 0, 0
	if c.bitcoindrpcClient != nil {
		dto, err := c.bitcoindrpcClient.ChainEstimateSmartFeeRate(ctx, targetBlocks)
		if err != nil {
			return 0, 0, err
		}
		fee1 = dto.Feerate
	}

	if c.mempoolapisClient != nil {
		dto, err := c.mempoolapisClient.EstimateFeeRate(ctx, targetBlocks)
		if err != nil {
			return 0, 0, err
		}
		fee2 = dto.FastestFee
	}

	return fee1, fee2, nil
}
