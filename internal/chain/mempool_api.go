package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// GetMempoolStats 获取 mempool.space 内存池统计
func (c *Client) GetMempoolStats(ctx context.Context) (*types.MempoolStats, error) {
	if c.mempoolapisClient == nil {
		return nil, ErrProviderUnavailable
	}
	dto, err := c.mempoolapisClient.GetMempoolStats(ctx)
	if err != nil {
		return nil, err
	}
	return &types.MempoolStats{
		Count:        dto.Count,
		Vsize:        dto.Vsize,
		TotalFee:     dto.TotalFee,
		FeeHistogram: dto.FeeHistogram,
	}, nil
}

// GetMempoolTxStatus 获取交易确认状态
func (c *Client) GetMempoolTxStatus(ctx context.Context, txid string) (*types.MempoolTxStatus, error) {
	if c.mempoolapisClient == nil {
		return nil, ErrProviderUnavailable
	}
	status, err := c.mempoolapisClient.GetTxStatus(ctx, txid)
	if err != nil {
		return nil, err
	}
	result := &types.MempoolTxStatus{
		TxID:          txid,
		Confirmed:     status.Confirmed,
		BlockHeight:   status.BlockHeight,
		BlockHash:     status.BlockHash,
		BlockTime:     status.BlockTime,
		InMempool:     !status.Confirmed,
		Confirmations: 0,
	}
	if status.Confirmed {
		result.InMempool = false
		result.Confirmations = 1
	}
	tx, err := c.mempoolapisClient.TxGetVerbose(ctx, txid)
	if err == nil && tx != nil {
		result.Confirmations = tx.Confirmations
		result.Replaceable = txHasReplaceableSequence(tx.Vin)
		if tx.Fee > 0 && tx.Vsize > 0 {
			result.Fee = int64(tx.Fee)
			result.FeeRateSatVB = tx.Fee / float64(tx.Vsize)
		}
	}
	return result, nil
}

func txHasReplaceableSequence(vin []types.TxVerboseInput) bool {
	for _, input := range vin {
		if input.Sequence < 0xfffffffe {
			return true
		}
	}
	return false
}

// GetMempoolFeesRecommend 从 mempool.space 获取多档推荐费率
func (c *Client) GetMempoolFeesRecommend(ctx context.Context) (*types.FeesRecommend, error) {
	if c.mempoolapisClient == nil {
		return nil, ErrProviderUnavailable
	}
	fees, err := c.mempoolapisClient.GetRecommendedFees(ctx)
	if err != nil {
		return nil, err
	}
	return &types.FeesRecommend{
		Blocks1: *feesFromMempoolMap(fees, 1),
		Blocks3: *feesFromMempoolMap(fees, 3),
		Blocks6: *feesFromMempoolMap(fees, 6),
		Source:  "mempoolSpace",
	}, nil
}

func feesFromMempoolMap(fees map[string]float64, blocks int64) *types.ChainFees {
	get := func(key string) float64 {
		if v, ok := fees[key]; ok {
			return v
		}
		return 0
	}
	return ChainFeesFromRecommended(blocks, types.RecommendedFees{
		FastestFee:  get("fastestFee"),
		HalfHourFee: get("halfHourFee"),
		HourFee:     get("hourFee"),
		EconomyFee:  get("economyFee"),
		MinimumFee:  get("minimumFee"),
	})
}
