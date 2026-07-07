package chain

import (
	"context"
	"testing"

	"github.com/crazycloudcc/btcapis/types"
)

func TestPickFeeByTarget(t *testing.T) {
	if got := PickFeeByTarget(1, 10, 8, 6, 4, 2); got != 10 {
		t.Fatalf("expected 10, got %v", got)
	}
	if got := PickFeeByTarget(3, 10, 8, 6, 4, 2); got != 8 {
		t.Fatalf("expected 8, got %v", got)
	}
	if got := PickFeeByTarget(20, 10, 8, 6, 4, 2); got != 2 {
		t.Fatalf("expected 2, got %v", got)
	}
}

func TestNormalizeChainFees(t *testing.T) {
	fees := NormalizeChainFees(&types.ChainFees{
		High:    21.456,
		Medium:  15.444,
		Low:     5.001,
		FeeRate: 0,
		Blocks:  3,
	})
	if fees.High != 21.46 || fees.Medium != 15.44 || fees.Low != 5 || fees.FeeRate != 15.44 || fees.Blocks != 3 {
		t.Fatalf("unexpected fees: %+v", fees)
	}
}

func TestChainFeesFromRecommended(t *testing.T) {
	fees := ChainFeesFromRecommended(3, types.RecommendedFees{
		FastestFee:  21,
		HalfHourFee: 15,
		HourFee:     10,
		EconomyFee:  5,
		MinimumFee:  1,
	})
	if fees.High != 21 || fees.Medium != 15 || fees.Low != 5 || fees.FeeRate != 15 || fees.Blocks != 3 {
		t.Fatalf("unexpected fees: %+v", fees)
	}
}

func TestRecommendedFeesAdapter(t *testing.T) {
	adapter := RecommendedFeesAdapter{Provider: stubRecommendedProvider{}}
	fees, err := adapter.EstimateChainFees(context.Background(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fees.FeeRate != 10 {
		t.Fatalf("expected feerate 10, got %v", fees.FeeRate)
	}
}

type stubRecommendedProvider struct{}

func (stubRecommendedProvider) Name() string { return "stub" }

func (stubRecommendedProvider) GetRecommendedFees(context.Context) (types.RecommendedFees, error) {
	return types.RecommendedFees{FastestFee: 21, HalfHourFee: 15, HourFee: 10, EconomyFee: 5, MinimumFee: 1}, nil
}
