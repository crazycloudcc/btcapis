// Package types 网络相关类型定义
// 包含比特币网络节点信息、连接状态、网络统计等网络层相关的数据结构
package types

import "github.com/btcsuite/btcd/chaincfg"

type Network string

const (
	Mainnet Network = "mainnet"
	Testnet Network = "testnet"
	Signet  Network = "signet"
	Regtest Network = "regtest"
)

var CurrentNetwork Network = Mainnet
var CurrentNetworkParams *chaincfg.Params = &chaincfg.MainNetParams

func SetCurrentNetwork(net string) {
	CurrentNetwork = Network(net)
	CurrentNetworkParams = CurrentNetwork.ToParams()
}

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
