package utils

import (
	"math"
	"testing"
)

func TestSatPerVBFromBTCPerKB(t *testing.T) {
	got := SatPerVBFromBTCPerKB(0.00015)
	if math.Abs(got-15) > 1e-6 {
		t.Fatalf("expected ~15, got %v", got)
	}
}

func TestRoundFeeSatVB(t *testing.T) {
	if got := RoundFeeSatVB(12.345); got != 12.35 {
		t.Fatalf("expected 12.35, got %v", got)
	}
}

func TestParseFloatFee(t *testing.T) {
	if got := ParseFloatFee(" 12.5 "); got != 12.5 {
		t.Fatalf("expected 12.5, got %v", got)
	}
}
