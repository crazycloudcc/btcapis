// Package address Bech32地址处理
package address

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil/bech32"
)

// Bech32Encode 对 hrp 与数据进行 Bech32 编码。
// data 应为 5-bit 分组（0..31）。若是 8-bit 原始数据，请自行先做转换。
func Bech32Encode(hrp string, data []byte) (string, error) {
	return bech32.Encode(hrp, data)
}

// Bech32Decode 解析 Bech32 字符串，返回 hrp 与 5-bit 数据。
func Bech32Decode(addr string) (string, []byte, error) {
	hrp, data, err := bech32.Decode(addr)
	if err != nil {
		return "", nil, err
	}
	if len(hrp) == 0 {
		return "", nil, errors.New("invalid bech32 hrp")
	}
	return hrp, data, nil
}
