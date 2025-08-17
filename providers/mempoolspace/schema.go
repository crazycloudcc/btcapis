// Package mempoolspace 响应模式定义
package mempoolspace

import (
	"time"
)

// FeeResponse 手续费响应
type FeeResponse struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	EconomyFee  int `json:"economyFee"`
	MinimumFee  int `json:"minimumFee"`
}

// MempoolResponse 内存池响应
type MempoolResponse struct {
	Count        int                 `json:"count"`
	VSize        int64               `json:"vsize"`
	TotalFee     int64               `json:"total_fee"`
	FeeHistogram []FeeHistogramEntry `json:"fee_histogram"`
}

// FeeHistogramEntry 手续费直方图条目
type FeeHistogramEntry struct {
	FeeRate int   `json:"feeRate"`
	VSize   int64 `json:"vsize"`
}

// BlockResponse 区块响应
type BlockResponse struct {
	ID                string    `json:"id"`
	Height            int64     `json:"height"`
	Version           int       `json:"version"`
	Timestamp         time.Time `json:"timestamp"`
	TxCount           int       `json:"tx_count"`
	Size              int       `json:"size"`
	Weight            int       `json:"weight"`
	MerkleRoot        string    `json:"merkle_root"`
	PreviousBlockHash string    `json:"previous_block_hash"`
	MedianFee         int       `json:"median_fee"`
	FeeRange          []int     `json:"fee_range"`
	Reward            int64     `json:"reward"`
	FeeTotal          int64     `json:"fee_total"`
}

// TransactionResponse 交易响应
type TransactionResponse struct {
	TxID     string     `json:"txid"`
	Version  int        `json:"version"`
	LockTime int64      `json:"locktime"`
	Size     int        `json:"size"`
	Weight   int        `json:"weight"`
	Fee      int64      `json:"fee"`
	Status   TxStatus   `json:"status"`
	Inputs   []TxInput  `json:"vin"`
	Outputs  []TxOutput `json:"vout"`
}

// TxStatus 交易状态
type TxStatus struct {
	Confirmed   bool      `json:"confirmed"`
	BlockHeight int64     `json:"block_height"`
	BlockHash   string    `json:"block_hash"`
	BlockTime   time.Time `json:"block_time"`
}

// TxInput 交易输入
type TxInput struct {
	TxID       string    `json:"txid"`
	Vout       int       `json:"vout"`
	Prevout    *TxOutput `json:"prevout"`
	ScriptSig  string    `json:"scriptsig"`
	Witness    []string  `json:"witness"`
	IsCoinbase bool      `json:"is_coinbase"`
	Sequence   int64     `json:"sequence"`
}

// TxOutput 交易输出
type TxOutput struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyASM     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}
