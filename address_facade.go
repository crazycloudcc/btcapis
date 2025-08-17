// Package btcapis 地址模块门面
package btcapis

import (
	"github.com/crazycloudcc/btcapis/types"
)

// addressFacade 提供地址相关的功能接口
type addressFacade struct{}

// Parse 解析比特币地址
func (a *addressFacade) Parse(addr string, network types.Network) (*types.AddressInfo, error) {
	// TODO: 实现地址解析
	return nil, nil
}

// Validate 验证地址格式
func (a *addressFacade) Validate(addr string, network types.Network) error {
	// TODO: 实现地址验证
	return nil
}

// Classify 分类地址类型
func (a *addressFacade) Classify(addr string) (types.AddressType, error) {
	// TODO: 实现地址分类
	return types.AddressTypeUnknown, nil
}

// ToScriptPubKey 转换为脚本公钥
func (a *addressFacade) ToScriptPubKey(addr string) ([]byte, error) {
	// TODO: 实现地址到脚本公钥转换
	return nil, nil
}
