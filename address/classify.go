// Package address 地址分类
package address

import (
	"strings"

	"github.com/crazycloudcc/btcapis/types"
)

// ClassifyAddress 分类地址类型
func ClassifyAddress(addr string) (types.AddressType, error) {
	if addr == "" {
		return types.AddressTypeUnknown, nil
	}

	// 根据地址前缀和格式进行分类
	if strings.HasPrefix(addr, "1") || strings.HasPrefix(addr, "3") {
		return types.AddressTypeP2PKH, nil
	}

	if strings.HasPrefix(addr, "bc1") || strings.HasPrefix(addr, "tb1") {
		// 进一步分类Bech32地址
		return classifyBech32(addr)
	}

	return types.AddressTypeUnknown, nil
}

// classifyBech32 分类Bech32地址
func classifyBech32(addr string) (types.AddressType, error) {
	// TODO: 实现Bech32地址分类
	// 根据地址长度和格式判断类型

	if len(addr) == 42 {
		return types.AddressTypeP2WPKH, nil
	}

	if len(addr) == 62 {
		return types.AddressTypeP2WSH, nil
	}

	if len(addr) == 62 && strings.HasPrefix(addr, "bc1p") {
		return types.AddressTypeP2TR, nil
	}

	return types.AddressTypeUnknown, nil
}

// IsValidAddress 检查地址是否有效
func IsValidAddress(addr string, network types.Network) bool {
	if addr == "" {
		return false
	}

	// TODO: 实现地址有效性检查
	// 根据网络类型和地址格式进行验证

	return true
}

// GetNetworkFromAddress 从地址推断网络类型
func GetNetworkFromAddress(addr string) types.Network {
	if addr == "" {
		return types.NetworkMainnet
	}

	// TODO: 实现网络类型推断
	// 根据地址前缀判断网络

	if strings.HasPrefix(addr, "tb1") || strings.HasPrefix(addr, "2") {
		return types.NetworkTestnet
	}

	if strings.HasPrefix(addr, "bcrt1") || strings.HasPrefix(addr, "m") || strings.HasPrefix(addr, "n") {
		return types.NetworkRegtest
	}

	return types.NetworkMainnet
}
