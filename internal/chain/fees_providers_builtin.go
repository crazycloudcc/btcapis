package chain

import (
	"context"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/electrumx"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/utils"
	"github.com/crazycloudcc/btcapis/types"
)

// MempoolFeesProvider mempool.space 推荐费率
type MempoolFeesProvider struct {
	Client *mempoolapis.Client
}

func (p MempoolFeesProvider) Name() string { return "mempoolSpace" }

func (p MempoolFeesProvider) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	if p.Client == nil {
		return nil, ErrProviderUnavailable
	}
	dto, err := p.Client.EstimateFeeRate(ctx, int(targetBlocks))
	if err != nil {
		return nil, err
	}
	blocks := utils.ClampChainFeeTargetBlocks(targetBlocks)
	return NormalizeChainFees(&types.ChainFees{
		High:    dto.FastestFee,
		Medium:  dto.HalfHourFee,
		Low:     dto.EconomyFee,
		FeeRate: PickFeeByTarget(blocks, dto.FastestFee, dto.HalfHourFee, dto.HourFee, dto.EconomyFee, dto.MinimumFee),
		Blocks:  blocks,
	}), nil
}

// ElectrumXFeesProvider electrumX blockchain.estimatefee
type ElectrumXFeesProvider struct {
	Client *electrumx.Client
}

func (p ElectrumXFeesProvider) Name() string { return "electrumX" }

func (p ElectrumXFeesProvider) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	if p.Client == nil {
		return nil, ErrProviderUnavailable
	}
	blocks := utils.ClampChainFeeTargetBlocks(targetBlocks)
	feeRateBTCPerKB, err := p.estimateFeeBTCPerKB(ctx, int(blocks))
	if err != nil {
		return nil, err
	}
	highBTCPerKB, errHigh := p.estimateFeeBTCPerKB(ctx, 1)
	mediumBTCPerKB, errMedium := p.estimateFeeBTCPerKB(ctx, 6)
	lowBTCPerKB, errLow := p.estimateFeeBTCPerKB(ctx, 25)
	fallback := feeRateBTCPerKB
	if errHigh != nil {
		highBTCPerKB = fallback
	}
	if errMedium != nil {
		mediumBTCPerKB = fallback
	}
	if errLow != nil {
		lowBTCPerKB = fallback
	}
	return NormalizeChainFees(&types.ChainFees{
		High:    utils.SatPerVBFromBTCPerKB(highBTCPerKB),
		Medium:  utils.SatPerVBFromBTCPerKB(mediumBTCPerKB),
		Low:     utils.SatPerVBFromBTCPerKB(lowBTCPerKB),
		FeeRate: utils.SatPerVBFromBTCPerKB(feeRateBTCPerKB),
		Blocks:  blocks,
	}), nil
}

func (p ElectrumXFeesProvider) estimateFeeBTCPerKB(ctx context.Context, blocks int) (float64, error) {
	if blocks < 1 {
		blocks = 1
	}
	fee, err := p.Client.EstimateFee(ctx, blocks)
	if err != nil {
		return 0, fmt.Errorf("electrumx estimatefee(%d) failed: %w", blocks, err)
	}
	if fee < 0 {
		return 0, fmt.Errorf("electrumx estimatefee(%d) returned insufficient data", blocks)
	}
	return fee, nil
}

// BitcoindFeesProvider Bitcoin Core estimatesmartfee
type BitcoindFeesProvider struct {
	Client *bitcoindrpc.Client
}

func (p BitcoindFeesProvider) Name() string { return "bitcoinCore" }

func (p BitcoindFeesProvider) EstimateChainFees(ctx context.Context, targetBlocks int64) (*types.ChainFees, error) {
	if p.Client == nil {
		return nil, ErrProviderUnavailable
	}
	blocks := utils.ClampChainFeeTargetBlocks(targetBlocks)
	feeRateSatVB, err := p.estimateSmartFeeSatVB(ctx, int(blocks))
	if err != nil {
		return nil, err
	}
	highSatVB, errHigh := p.estimateSmartFeeSatVB(ctx, 1)
	mediumSatVB, errMedium := p.estimateSmartFeeSatVB(ctx, 6)
	lowSatVB, errLow := p.estimateSmartFeeSatVB(ctx, 25)
	fallback := feeRateSatVB
	if errHigh != nil {
		highSatVB = fallback
	}
	if errMedium != nil {
		mediumSatVB = fallback
	}
	if errLow != nil {
		lowSatVB = fallback
	}
	return NormalizeChainFees(&types.ChainFees{
		High:    highSatVB,
		Medium:  mediumSatVB,
		Low:     lowSatVB,
		FeeRate: feeRateSatVB,
		Blocks:  blocks,
	}), nil
}

func (p BitcoindFeesProvider) estimateSmartFeeSatVB(ctx context.Context, blocks int) (float64, error) {
	blocks = int(utils.ClampChainFeeTargetBlocks(int64(blocks)))
	dto, err := p.Client.ChainEstimateSmartFeeRate(ctx, blocks)
	if err != nil {
		return 0, fmt.Errorf("estimatesmartfee(%d) failed: %w", blocks, err)
	}
	if dto == nil || dto.Feerate <= 0 {
		return 0, fmt.Errorf("estimatesmartfee(%d) returned insufficient data", blocks)
	}
	return utils.SatPerVBFromBTCPerKB(dto.Feerate), nil
}
