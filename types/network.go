// Package types 网络相关类型定义
// 包含比特币网络节点信息、连接状态、网络统计等网络层相关的数据结构
package types

import "github.com/btcsuite/btcd/chaincfg"

// Network 网络类型
type Network string

const (
	Mainnet Network = "mainnet" // 主网
	Testnet Network = "testnet" // 测试网
	Signet  Network = "signet"  // 签名网
	Regtest Network = "regtest" // 回归测试网
)

var CurrentNetwork Network = Mainnet                                // 当前网络
var CurrentNetworkParams *chaincfg.Params = &chaincfg.MainNetParams // 当前网络参数

// SetCurrentNetwork 设置当前网络
func SetCurrentNetwork(net string) {
	CurrentNetwork = Network(net)
	CurrentNetworkParams = CurrentNetwork.ToParams()
}

// ToParams 将网络转换为网络参数
func (n Network) ToParams() *chaincfg.Params {
	switch n {
	case Mainnet:
		return &chaincfg.MainNetParams
	case Testnet:
		return &chaincfg.TestNet3Params
	case Signet:
		return &chaincfg.SigNetParams
	case Regtest:
		return &chaincfg.RegressionNetParams
	default:
		return nil
	}
}
