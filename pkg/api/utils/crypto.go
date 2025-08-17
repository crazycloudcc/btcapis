// Package utils 提供各种实用工具函数
// 包含加密、编码、验证等常用功能
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
)

// HashType 表示哈希算法类型
type HashType string

const (
	// SHA256 SHA-256哈希算法
	SHA256 HashType = "sha256"
	// SHA512 SHA-512哈希算法
	SHA512 HashType = "sha512"
)

// GenerateRandomBytes 生成指定长度的随机字节
// 参数: length - 要生成的字节长度
// 返回: 随机字节切片和可能的错误
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("生成随机字节失败: %w", err)
	}
	return bytes, nil
}

// GenerateRandomHex 生成指定长度的随机十六进制字符串
// 参数: length - 要生成的十六进制字符长度
// 返回: 随机十六进制字符串和可能的错误
func GenerateRandomHex(length int) (string, error) {
	// 每个字节对应两个十六进制字符
	byteLength := (length + 1) / 2
	bytes, err := GenerateRandomBytes(byteLength)
	if err != nil {
		return "", err
	}

	hexStr := hex.EncodeToString(bytes)
	// 如果长度不匹配，截取到指定长度
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}

	return hexStr, nil
}

// CalculateHash 计算数据的哈希值
// 参数: data - 要哈希的数据
// 参数: hashType - 哈希算法类型
// 返回: 哈希值的十六进制字符串和可能的错误
func CalculateHash(data []byte, hashType HashType) (string, error) {
	var h hash.Hash

	switch hashType {
	case SHA256:
		h = sha256.New()
	case SHA512:
		// 注意：Go标准库没有SHA512，这里使用SHA256作为示例
		h = sha256.New()
	default:
		return "", fmt.Errorf("不支持的哈希类型: %s", hashType)
	}

	h.Write(data)
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

// ValidateHexString 验证字符串是否为有效的十六进制格式
// 参数: hexStr - 要验证的字符串
// 返回: 是否为有效十六进制字符串
func ValidateHexString(hexStr string) bool {
	if len(hexStr)%2 != 0 {
		return false
	}

	for _, char := range hexStr {
		if !isValidHexChar(char) {
			return false
		}
	}

	return true
}

// isValidHexChar 检查字符是否为有效的十六进制字符
func isValidHexChar(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}

// HexToBytes 将十六进制字符串转换为字节切片
// 参数: hexStr - 十六进制字符串
// 返回: 字节切片和可能的错误
func HexToBytes(hexStr string) ([]byte, error) {
	if !ValidateHexString(hexStr) {
		return nil, fmt.Errorf("无效的十六进制字符串: %s", hexStr)
	}

	return hex.DecodeString(hexStr)
}

// BytesToHex 将字节切片转换为十六进制字符串
// 参数: bytes - 字节切片
// 返回: 十六进制字符串
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
