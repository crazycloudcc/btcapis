package chain

import (
	"github.com/crazycloudcc/btcapis/internal/utils"
	"github.com/crazycloudcc/btcapis/types"
)

// PickFeeByTarget 按目标确认区块数选取费率档位（sat/vB）
func PickFeeByTarget(targetBlocks int64, fastest, halfHour, hour, economy, minimum float64) float64 {
	switch {
	case targetBlocks <= 1:
		return fastest
	case targetBlocks <= 3:
		return halfHour
	case targetBlocks <= 6:
		return hour
	case targetBlocks <= 12:
		return economy
	default:
		if minimum > 0 {
			return minimum
		}
		return economy
	}
}

// NormalizeChainFees 对齐费率字段并统一为 sat/vB（保留 2 位小数）
func NormalizeChainFees(fees *types.ChainFees) *types.ChainFees {
	if fees == nil {
		return nil
	}
	blocks := utils.ClampChainFeeTargetBlocks(fees.Blocks)
	normalized := &types.ChainFees{
		High:    utils.RoundFeeSatVB(fees.High),
		Medium:  utils.RoundFeeSatVB(fees.Medium),
		Low:     utils.RoundFeeSatVB(fees.Low),
		FeeRate: utils.RoundFeeSatVB(fees.FeeRate),
		Blocks:  blocks,
	}
	if normalized.FeeRate == 0 {
		normalized.FeeRate = PickFeeByTarget(
			blocks,
			normalized.High,
			normalized.Medium,
			normalized.Medium,
			normalized.Low,
			normalized.Low,
		)
		normalized.FeeRate = utils.RoundFeeSatVB(normalized.FeeRate)
	}
	return normalized
}

// ChainFeesFromRecommended 将多档推荐费率映射为统一结构（sat/vB）
func ChainFeesFromRecommended(targetBlocks int64, recommended types.RecommendedFees) *types.ChainFees {
	blocks := utils.ClampChainFeeTargetBlocks(targetBlocks)
	return NormalizeChainFees(&types.ChainFees{
		High:    recommended.FastestFee,
		Medium:  recommended.HalfHourFee,
		Low:     recommended.EconomyFee,
		FeeRate: PickFeeByTarget(blocks, recommended.FastestFee, recommended.HalfHourFee, recommended.HourFee, recommended.EconomyFee, recommended.MinimumFee),
		Blocks:  blocks,
	})
}
