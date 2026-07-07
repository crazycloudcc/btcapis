package utils

import (
	"math"
	"strconv"
	"strings"
)

// SatPerVBFromBTCPerKB 将 BTC/kB 费率转换为 sat/vB
func SatPerVBFromBTCPerKB(btcPerKB float64) float64 {
	if btcPerKB < 0 {
		return 0
	}
	return btcPerKB * 1e5
}

// RoundFeeSatVB 将 sat/vB 费率四舍五入保留 2 位小数
func RoundFeeSatVB(value float64) float64 {
	if value < 0 {
		return 0
	}
	return math.Round(value*100) / 100
}

// ClampChainFeeTargetBlocks 限制目标确认区块数范围（1-1008）
func ClampChainFeeTargetBlocks(targetBlocks int64) int64 {
	if targetBlocks < 1 {
		return 1
	}
	if targetBlocks > 1008 {
		return 1008
	}
	return targetBlocks
}

// ParseFloatFee 解析字符串费率字段
func ParseFloatFee(raw string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0
	}
	return v
}
