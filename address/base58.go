// Package address Base58地址处理
package address

import (
	"github.com/btcsuite/btcd/btcutil/base58"
)

// Base58Encode 对字节数据进行 Base58Check 编码。
// 注意：此处仅包装 btcd 的 base58.Encode（不含校验和计算）。
func Base58Encode(data []byte) string {
	return base58.Encode(data)
}

// Base58Decode 解码 Base58 字符串为原始字节。
// 注意：此处仅包装 btcd 的 base58.Decode（不校验校验和）。
func Base58Decode(s string) []byte {
	return base58.Decode(s)
}
