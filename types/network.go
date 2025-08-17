// Package types 网络相关类型定义
// 包含比特币网络节点信息、连接状态、网络统计等网络层相关的数据结构
package types

import (
	"time"
)

// NetworkInfo 表示网络信息
// 包含比特币节点的基本网络配置和状态信息
type NetworkInfo struct {
	// Name 网络名称
	// 如 "mainnet", "testnet", "regtest" 等
	Name string `json:"name"`
	// Version 节点版本号
	// 比特币核心软件的版本号
	Version int `json:"version"`
	// Subversion 子版本字符串
	// 完整的版本标识符，如 "/Satoshi:0.21.0/"
	Subversion string `json:"subversion"`
	// ProtocolVersion 协议版本号
	// 比特币网络协议版本，如 70016
	ProtocolVersion int `json:"protocol_version"`
	// Connections 当前连接数
	// 节点当前维护的网络连接总数
	Connections int `json:"connections"`
	// LastUpdated 最后更新时间
	// 网络信息最后更新的时间戳
	LastUpdated time.Time `json:"last_updated"`
}

// PeerInfo 表示节点信息
// 包含与特定对等节点的连接详情和状态信息
type PeerInfo struct {
	// ID 节点ID
	// 节点的唯一标识符
	ID int `json:"id"`
	// Addr 节点地址
	// 节点的IP地址和端口，如 "192.168.1.1:8333"
	Addr string `json:"addr"`
	// AddrBind 绑定地址
	// 节点绑定的本地地址
	AddrBind string `json:"addr_bind"`
	// AddrLocal 本地地址
	// 节点的本地网络地址
	AddrLocal string `json:"addr_local"`
	// Services 服务标识
	// 节点提供的服务标识，如 "0000000000000009"
	Services string `json:"services"`
	// RelayTxes 是否中继交易
	// 节点是否转发交易到其他节点
	RelayTxes bool `json:"relay_txes"`
	// LastSend 最后发送时间
	// 最后向该节点发送数据的时间
	LastSend time.Time `json:"last_send"`
	// LastRecv 最后接收时间
	// 最后从该节点接收数据的时间
	LastRecv time.Time `json:"last_recv"`
	// BytesSent 发送字节数
	// 向该节点发送的总字节数
	BytesSent int64 `json:"bytes_sent"`
	// BytesRecv 接收字节数
	// 从该节点接收的总字节数
	BytesRecv int64 `json:"bytes_recv"`
	// ConnTime 连接时间
	// 与该节点建立连接的时间
	ConnTime time.Time `json:"conn_time"`
	// TimeOffset 时间偏移
	// 与节点的时间差（秒）
	TimeOffset int64 `json:"time_offset"`
	// PingTime 当前ping时间
	// 最近一次ping的往返时间（毫秒）
	PingTime float64 `json:"ping_time"`
	// MinPing 最小ping时间
	// 记录的最小ping时间（毫秒）
	MinPing float64 `json:"min_ping"`
	// Version 节点版本
	// 对等节点的软件版本号
	Version int `json:"version"`
	// SubVer 子版本
	// 对等节点的完整版本字符串
	SubVer string `json:"sub_ver"`
	// Inbound 是否为入站连接
	// true表示该节点主动连接到我们，false表示我们主动连接
	Inbound bool `json:"inbound"`
	// StartingHeight 起始区块高度
	// 对等节点开始同步的区块高度
	StartingHeight int64 `json:"starting_height"`
	// BanScore 禁止分数
	// 节点的违规行为评分，达到阈值会被禁止
	BanScore int `json:"ban_score"`
	// SyncedHeaders 已同步区块头数
	// 已同步的区块头数量
	SyncedHeaders int64 `json:"synced_headers"`
	// SyncedBlocks 已同步区块数
	// 已同步的完整区块数量
	SyncedBlocks int64 `json:"synced_blocks"`
	// Inflight 传输中的区块
	// 正在从该节点下载的区块高度列表
	Inflight []int `json:"inflight"`
	// Whitelisted 是否在白名单中
	// 节点是否被加入白名单，白名单节点不会被禁止
	Whitelisted bool `json:"whitelisted"`
	// MinRelayFee 最小中继手续费
	// 节点接受的最小交易手续费率
	MinRelayFee float64 `json:"min_relay_fee"`
	// LastPingNonce 最后ping随机数
	// 最近一次ping使用的随机数
	LastPingNonce uint64 `json:"last_ping_nonce"`
	// LastPingTime 最后ping时间
	// 最近一次ping的时间戳
	LastPingTime time.Time `json:"last_ping_time"`
	// LastPingMicros 最后ping微秒数
	// 最近一次ping的往返时间（微秒）
	LastPingMicros int64 `json:"last_ping_micros"`
}

// NetworkStats 表示网络统计信息
// 包含整个网络的连接和流量统计汇总
type NetworkStats struct {
	// TotalConnections 总连接数
	// 当前所有网络连接的总数
	TotalConnections int `json:"total_connections"`
	// InboundConnections 入站连接数
	// 其他节点主动连接到我们的数量
	InboundConnections int `json:"inbound_connections"`
	// OutboundConnections 出站连接数
	// 我们主动连接到其他节点的数量
	OutboundConnections int `json:"outbound_connections"`
	// TotalBytesSent 总发送字节数
	// 向所有节点发送的总字节数
	TotalBytesSent int64 `json:"total_bytes_sent"`
	// TotalBytesRecv 总接收字节数
	// 从所有节点接收的总字节数
	TotalBytesRecv int64 `json:"total_bytes_recv"`
	// LastUpdated 最后更新时间
	// 网络统计信息最后更新的时间戳
	LastUpdated time.Time `json:"last_updated"`
}
