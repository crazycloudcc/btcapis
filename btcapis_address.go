package btcapis

import "github.com/crazycloudcc/btcapis/types"

// Parse: 先返回最小结构（占位实现）。后续可以接入 bech32/base58 真正解析。
func (AddressModule) Parse(addr string, net types.Network) (types.AddressInfo, error) {
	return types.AddressInfo{
		Address: addr,
		Network: net,
		Type:    types.AddrUnknown,
	}, nil
}
