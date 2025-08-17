// Package common 提供通用的类型定义和工具函数
// 这些类型和函数被其他API包共享使用
package common

import (
	"time"
)

// NetworkType 表示区块链网络类型
type NetworkType string

const (
	// Mainnet 主网
	Mainnet NetworkType = "mainnet"
	// Testnet 测试网
	Testnet NetworkType = "testnet"
	// Regtest 回归测试网
	Regtest NetworkType = "regtest"
)

// TransactionStatus 表示交易状态
type TransactionStatus string

const (
	// Pending 待确认
	Pending TransactionStatus = "pending"
	// Confirmed 已确认
	Confirmed TransactionStatus = "confirmed"
	// Failed 失败
	Failed TransactionStatus = "failed"
)

// Transaction 表示通用交易结构
type Transaction struct {
	TxHash      string            `json:"txHash"`      // 交易哈希
	From        string            `json:"from"`        // 发送方地址
	To          string            `json:"to"`          // 接收方地址
	Amount      string            `json:"amount"`      // 交易金额
	Fee         string            `json:"fee"`         // 交易费用
	Status      TransactionStatus `json:"status"`      // 交易状态
	BlockNumber uint64            `json:"blockNumber"` // 区块号
	Timestamp   time.Time         `json:"timestamp"`   // 时间戳
	Data        []byte            `json:"data"`        // 交易数据
}

// Block 表示通用区块结构
type Block struct {
	BlockHash    string    `json:"blockHash"`    // 区块哈希
	BlockNumber  uint64    `json:"blockNumber"`  // 区块号
	ParentHash   string    `json:"parentHash"`   // 父区块哈希
	Timestamp    time.Time `json:"timestamp"`    // 时间戳
	Transactions []string  `json:"transactions"` // 交易哈希列表
	Miner        string    `json:"miner"`        // 矿工地址
	Difficulty   string    `json:"difficulty"`   // 难度值
	GasLimit     uint64    `json:"gasLimit"`     // Gas限制
	GasUsed      uint64    `json:"gasUsed"`      // 已使用Gas
	ExtraData    []byte    `json:"extraData"`    // 额外数据
}

// APIResponse 表示通用API响应结构
type APIResponse struct {
	Success bool        `json:"success"` // 请求是否成功
	Data    interface{} `json:"data"`    // 响应数据
	Error   string      `json:"error"`   // 错误信息
	Code    int         `json:"code"`    // 响应状态码
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Error:   "",
		Code:    200,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(error string, code int) *APIResponse {
	return &APIResponse{
		Success: false,
		Data:    nil,
		Error:   error,
		Code:    code,
	}
}
