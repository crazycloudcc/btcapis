package utils

// SatsToBTC 将 satoshi 转换为 BTC
func SatsToBTC(sats int64) float64 {
	return float64(sats) / 1e8
}

// BTCToSats 将 BTC 转换为 satoshi
func BTCToSats(btc float64) int64 {
	return int64(btc * 1e8)
}
