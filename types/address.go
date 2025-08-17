// Package types 提供跨模块共享的轻量类型定义
package types

import (
	"time"
)

// Network 表示比特币网络类型
type Network string

// AddressType 表示地址类型
type AddressType string

const (
	// 网络类型常量
	NetworkMainnet Network = "mainnet"
	NetworkTestnet Network = "testnet"
	NetworkSignet  Network = "signet"
	NetworkRegtest Network = "regtest"

	// 地址类型常量
	AddressTypeP2PKH   AddressType = "p2pkh"
	AddressTypeP2SH    AddressType = "p2sh"
	AddressTypeP2WPKH  AddressType = "p2wpkh"
	AddressTypeP2WSH   AddressType = "p2wsh"
	AddressTypeP2TR    AddressType = "p2tr"
	AddressTypeUnknown AddressType = "unknown"
)

// AddressInfo 表示地址信息
type AddressInfo struct {
	Address      string      `json:"address"`
	Network      Network     `json:"network"`
	Type         AddressType `json:"type"`
	ScriptPubKey []byte      `json:"script_pub_key"`
	IsValid      bool        `json:"is_valid"`
	CreatedAt    time.Time   `json:"created_at"`
}

// OutPoint 表示UTXO输出点
type OutPoint struct {
	TxID string `json:"txid"`
	Vout uint32 `json:"vout"`
}

// UTXO 表示未花费交易输出
type UTXO struct {
	OutPoint     OutPoint `json:"outpoint"`
	ScriptPubKey []byte   `json:"script_pub_key"`
	Amount       int64    `json:"amount"`
	Height       int64    `json:"height"`
	Spent        bool     `json:"spent"`
}
