// Package types 网络相关类型定义
// 包含比特币网络节点信息、连接状态、网络统计等网络层相关的数据结构
package types

type Network string

const (
	Mainnet Network = "mainnet"
	Testnet Network = "testnet"
	Signet  Network = "signet"
	Regtest Network = "regtest"
)
