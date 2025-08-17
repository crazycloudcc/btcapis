// Package btc 提供比特币相关的API功能
// 包含地址生成、交易处理、余额查询等核心功能
package btc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// BTCAddress 表示比特币地址结构
type BTCAddress struct {
	Address     string `json:"address"`     // 地址字符串
	PrivateKey  string `json:"privateKey"`  // 私钥（十六进制）
	PublicKey   string `json:"publicKey"`   // 公钥（十六进制）
	AddressType string `json:"addressType"` // 地址类型（legacy, p2sh, bech32）
}

// BTCTransaction 表示比特币交易结构
type BTCTransaction struct {
	TxHash        string  `json:"txHash"`        // 交易哈希
	From          string  `json:"from"`          // 发送方地址
	To            string  `json:"to"`            // 接收方地址
	Amount        float64 `json:"amount"`        // 交易金额（BTC）
	Fee           float64 `json:"fee"`           // 交易费用（BTC）
	Confirmations int     `json:"confirmations"` // 确认数
	BlockHeight   int64   `json:"blockHeight"`   // 区块高度
	Timestamp     int64   `json:"timestamp"`     // 时间戳
	Size          int     `json:"size"`          // 交易大小（字节）
	Weight        int     `json:"weight"`        // 交易权重
}

// BTCBlock 表示比特币区块结构
type BTCBlock struct {
	BlockHash        string  `json:"blockHash"`        // 区块哈希
	BlockHeight      int64   `json:"blockHeight"`      // 区块高度
	PreviousHash     string  `json:"previousHash"`     // 前一个区块哈希
	MerkleRoot       string  `json:"merkleRoot"`       // Merkle根
	Timestamp        int64   `json:"timestamp"`        // 时间戳
	Difficulty       float64 `json:"difficulty"`       // 难度
	Nonce            uint32  `json:"nonce"`            // 随机数
	TransactionCount int     `json:"transactionCount"` // 交易数量
	Size             int     `json:"size"`             // 区块大小
	Weight           int     `json:"weight"`           // 区块权重
	Version          int32   `json:"version"`          // 版本
}

// GenerateAddress 生成新的比特币地址
// 返回包含地址、私钥和公钥的BTCAddress结构
func GenerateAddress() (*BTCAddress, error) {
	// 生成随机私钥（这里使用简化版本，实际应用中应使用更安全的密钥生成方法）
	privateKey := make([]byte, 32)
	if _, err := rand.Read(privateKey); err != nil {
		return nil, fmt.Errorf("生成私钥失败: %w", err)
	}

	// 计算公钥哈希（简化版本）
	publicKeyHash := sha256.Sum256(privateKey)

	// 生成地址（简化版本，实际应使用Base58Check编码）
	address := hex.EncodeToString(publicKeyHash[:20])

	return &BTCAddress{
		Address:     address,
		PrivateKey:  hex.EncodeToString(privateKey),
		PublicKey:   hex.EncodeToString(publicKeyHash[:]),
		AddressType: "legacy", // 默认使用legacy地址类型
	}, nil
}

// GenerateAddressWithType 根据指定类型生成比特币地址
// 参数: addressType - 地址类型（legacy, p2sh, bech32）
// 返回: BTCAddress结构
func GenerateAddressWithType(addressType string) (*BTCAddress, error) {
	address, err := GenerateAddress()
	if err != nil {
		return nil, err
	}

	// 设置地址类型
	switch strings.ToLower(addressType) {
	case "legacy", "p2sh", "bech32":
		address.AddressType = strings.ToLower(addressType)
	default:
		address.AddressType = "legacy"
	}

	return address, nil
}

// ValidateAddress 验证比特币地址格式
// 参数: address - 要验证的地址字符串
// 返回: 是否为有效地址
func ValidateAddress(address string) bool {
	// 这里应该实现完整的地址验证逻辑
	// 包括长度检查、字符集验证、校验和验证等
	if len(address) < 26 || len(address) > 35 {
		return false
	}

	// 检查是否只包含有效的Base58字符
	for _, char := range address {
		if !isValidBase58Char(char) {
			return false
		}
	}

	return true
}

// GetAddressType 获取比特币地址类型
// 参数: address - 地址字符串
// 返回: 地址类型字符串
func GetAddressType(address string) string {
	if !ValidateAddress(address) {
		return "invalid"
	}

	// 根据地址前缀判断类型（简化版本）
	if strings.HasPrefix(address, "1") {
		return "legacy"
	} else if strings.HasPrefix(address, "3") {
		return "p2sh"
	} else if strings.HasPrefix(address, "bc1") {
		return "bech32"
	}

	return "unknown"
}

// CalculateTransactionFee 计算交易费用（估算）
// 参数: inputCount - 输入数量
// 参数: outputCount - 输出数量
// 参数: feeRate - 费率（sat/byte）
// 返回: 估算的交易费用（BTC）
func CalculateTransactionFee(inputCount, outputCount int, feeRate float64) float64 {
	// 基础交易大小：版本(4) + 输入计数(1) + 输出计数(1) + 锁定时间(4) = 10字节
	baseSize := 10

	// 每个输入：outpoint(36) + 脚本长度(1) + 签名脚本 + 序列号(4)
	inputSize := inputCount * (36 + 1 + 73 + 4) // 73是典型的签名脚本大小

	// 每个输出：值(8) + 脚本长度(1) + 公钥脚本
	outputSize := outputCount * (8 + 1 + 25) // 25是典型的公钥脚本大小

	totalSize := baseSize + inputSize + outputSize

	// 计算费用（sat转换为BTC）
	feeSat := float64(totalSize) * feeRate
	return feeSat / 100000000 // 1 BTC = 100,000,000 sat
}

// ValidatePrivateKey 验证私钥格式
// 参数: privateKey - 私钥字符串（十六进制）
// 返回: 是否为有效私钥
func ValidatePrivateKey(privateKey string) bool {
	// 检查长度（32字节 = 64个十六进制字符）
	if len(privateKey) != 64 {
		return false
	}

	// 检查是否为有效十六进制
	for _, char := range privateKey {
		if !isValidHexChar(char) {
			return false
		}
	}

	return true
}

// isValidBase58Char 检查字符是否为有效的Base58字符
func isValidBase58Char(char rune) bool {
	// Base58字符集: 123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz
	validChars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, validChar := range validChars {
		if char == validChar {
			return true
		}
	}
	return false
}

// isValidHexChar 检查字符是否为有效的十六进制字符
func isValidHexChar(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}
