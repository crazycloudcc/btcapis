// Package types 网络相关类型定义
package types

import (
	"time"
)

// NetworkInfo 表示网络信息
type NetworkInfo struct {
	Name            string    `json:"name"`
	Version         int       `json:"version"`
	Subversion      string    `json:"subversion"`
	ProtocolVersion int       `json:"protocol_version"`
	Connections     int       `json:"connections"`
	LastUpdated     time.Time `json:"last_updated"`
}

// PeerInfo 表示节点信息
type PeerInfo struct {
	ID             int       `json:"id"`
	Addr           string    `json:"addr"`
	AddrBind       string    `json:"addr_bind"`
	AddrLocal      string    `json:"addr_local"`
	Services       string    `json:"services"`
	RelayTxes      bool      `json:"relay_txes"`
	LastSend       time.Time `json:"last_send"`
	LastRecv       time.Time `json:"last_recv"`
	BytesSent      int64     `json:"bytes_sent"`
	BytesRecv      int64     `json:"bytes_recv"`
	ConnTime       time.Time `json:"conn_time"`
	TimeOffset     int64     `json:"time_offset"`
	PingTime       float64   `json:"ping_time"`
	MinPing        float64   `json:"min_ping"`
	Version        int       `json:"version"`
	SubVer         string    `json:"sub_ver"`
	Inbound        bool      `json:"inbound"`
	StartingHeight int64     `json:"starting_height"`
	BanScore       int       `json:"ban_score"`
	SyncedHeaders  int64     `json:"synced_headers"`
	SyncedBlocks   int64     `json:"synced_blocks"`
	Inflight       []int     `json:"inflight"`
	Whitelisted    bool      `json:"whitelisted"`
	MinRelayTxFee  float64   `json:"min_relay_tx_fee"`
	LastPingNonce  uint64    `json:"last_ping_nonce"`
	LastPingTime   time.Time `json:"last_ping_time"`
	LastPingMicros int64     `json:"last_ping_micros"`
}

// NetworkStats 表示网络统计信息
type NetworkStats struct {
	TotalConnections    int       `json:"total_connections"`
	InboundConnections  int       `json:"inbound_connections"`
	OutboundConnections int       `json:"outbound_connections"`
	TotalBytesSent      int64     `json:"total_bytes_sent"`
	TotalBytesRecv      int64     `json:"total_bytes_recv"`
	LastUpdated         time.Time `json:"last_updated"`
}
