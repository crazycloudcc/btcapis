// Package btcapis 提供比特币区块链的完整API功能
// 通过一次导入即可访问所有BTC相关功能
package btcapis

// 重新导出所有BTC相关功能
import (
	"github.com/yourusername/btcapis/pkg/api/btc"
	"github.com/yourusername/btcapis/pkg/api/common"
	"github.com/yourusername/btcapis/pkg/api/config"
	"github.com/yourusername/btcapis/pkg/api/utils"
)

// BTC 相关类型和函数
type (
	// BTCAddress 比特币地址结构
	BTCAddress = btc.BTCAddress
	// BTCTransaction 比特币交易结构
	BTCTransaction = btc.BTCTransaction
	// BTCBlock 比特币区块结构
	BTCBlock = btc.BTCBlock
)

// BTC 相关函数
var (
	// GenerateAddress 生成新的比特币地址
	GenerateAddress = btc.GenerateAddress
	// GenerateAddressWithType 根据指定类型生成比特币地址
	GenerateAddressWithType = btc.GenerateAddressWithType
	// ValidateAddress 验证比特币地址格式
	ValidateAddress = btc.ValidateAddress
	// GetAddressType 获取比特币地址类型
	GetAddressType = btc.GetAddressType
	// CalculateTransactionFee 计算交易费用（估算）
	CalculateTransactionFee = btc.CalculateTransactionFee
	// ValidatePrivateKey 验证私钥格式
	ValidatePrivateKey = btc.ValidatePrivateKey
)

// Common 相关类型和函数
type (
	// NetworkType 区块链网络类型
	NetworkType = common.NetworkType
	// TransactionStatus 交易状态
	TransactionStatus = common.TransactionStatus
	// Transaction 通用交易结构
	Transaction = common.Transaction
	// Block 通用区块结构
	Block = common.Block
	// APIResponse 通用API响应结构
	APIResponse = common.APIResponse
)

// Common 相关常量
const (
	// Mainnet 主网
	Mainnet = common.Mainnet
	// Testnet 测试网
	Testnet = common.Testnet
	// Regtest 回归测试网
	Regtest = common.Regtest
	// Pending 待确认
	Pending = common.Pending
	// Confirmed 已确认
	Confirmed = common.Confirmed
	// Failed 失败
	Failed = common.Failed
)

// Common 相关函数
var (
	// NewSuccessResponse 创建成功响应
	NewSuccessResponse = common.NewSuccessResponse
	// NewErrorResponse 创建错误响应
	NewErrorResponse = common.NewErrorResponse
)

// Utils 相关类型和函数
type (
	// HashType 哈希算法类型
	HashType = utils.HashType
)

// Utils 相关常量
const (
	// SHA256 SHA-256哈希算法
	SHA256 = utils.SHA256
	// SHA512 SHA-512哈希算法
	SHA512 = utils.SHA512
)

// Utils 相关函数
var (
	// GenerateRandomBytes 生成指定长度的随机字节
	GenerateRandomBytes = utils.GenerateRandomBytes
	// GenerateRandomHex 生成指定长度的随机十六进制字符串
	GenerateRandomHex = utils.GenerateRandomHex
	// CalculateHash 计算数据的哈希值
	CalculateHash = utils.CalculateHash
	// ValidateHexString 验证字符串是否为有效的十六进制格式
	ValidateHexString = utils.ValidateHexString
	// HexToBytes 将十六进制字符串转换为字节切片
	HexToBytes = utils.HexToBytes
	// BytesToHex 将字节切片转换为十六进制字符串
	BytesToHex = utils.BytesToHex
)

// Config 相关类型和函数
type (
	// Config 应用配置结构
	Config = config.Config
)

// Config 相关函数
var (
	// LoadConfig 从文件加载配置
	LoadConfig = config.LoadConfig
	// LoadDefaultConfig 加载默认配置
	LoadDefaultConfig = config.LoadDefaultConfig
)
