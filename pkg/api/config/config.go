// Package config 提供配置管理功能
// 包含环境变量读取、配置文件解析等配置相关功能
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config 表示应用配置结构
type Config struct {
	// 区块链网络配置
	Network struct {
		// BTC网络配置
		BTC struct {
			MainnetRPC  string `json:"mainnetRPC"`  // 主网RPC地址
			TestnetRPC  string `json:"testnetRPC"`  // 测试网RPC地址
			RegtestRPC  string `json:"regtestRPC"`  // 回归测试网RPC地址
			NetworkType string `json:"networkType"` // 网络类型（mainnet, testnet, regtest）
		} `json:"btc"`
	} `json:"network"`

	// API配置
	API struct {
		Port         int    `json:"port"`         // 服务端口
		Host         string `json:"host"`         // 服务主机
		ReadTimeout  int    `json:"readTimeout"`  // 读取超时（秒）
		WriteTimeout int    `json:"writeTimeout"` // 写入超时（秒）
	} `json:"api"`

	// 日志配置
	Log struct {
		Level      string `json:"level"`      // 日志级别
		OutputPath string `json:"outputPath"` // 日志输出路径
		MaxSize    int    `json:"maxSize"`    // 日志文件最大大小（MB）
		MaxBackups int    `json:"maxBackups"` // 最大备份文件数
	} `json:"log"`

	// 安全配置
	Security struct {
		EnableHTTPS bool   `json:"enableHTTPS"` // 是否启用HTTPS
		CertFile    string `json:"certFile"`    // 证书文件路径
		KeyFile     string `json:"keyFile"`     // 私钥文件路径
	} `json:"security"`

	// BTC特定配置
	BTC struct {
		DefaultFeeRate float64 `json:"defaultFeeRate"` // 默认费率（sat/byte）
		Confirmations  int     `json:"confirmations"`  // 默认确认数
		AddressType    string  `json:"addressType"`    // 默认地址类型
	} `json:"btc"`
}

// LoadConfig 从文件加载配置
// 参数: configPath - 配置文件路径
// 返回: 配置结构和可能的错误
func LoadConfig(configPath string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析JSON配置
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 从环境变量覆盖配置
	config.overrideFromEnv()

	return &config, nil
}

// LoadDefaultConfig 加载默认配置
// 返回: 默认配置结构
func LoadDefaultConfig() *Config {
	config := &Config{}

	// 设置BTC网络默认值
	config.Network.BTC.MainnetRPC = "https://btc.getblock.io/mainnet"
	config.Network.BTC.TestnetRPC = "https://btc.getblock.io/testnet"
	config.Network.BTC.RegtestRPC = "http://localhost:18443"
	config.Network.BTC.NetworkType = "mainnet"

	// API配置
	config.API.Port = 8080
	config.API.Host = "0.0.0.0"
	config.API.ReadTimeout = 30
	config.API.WriteTimeout = 30

	// 日志配置
	config.Log.Level = "info"
	config.Log.OutputPath = "./logs"
	config.Log.MaxSize = 100
	config.Log.MaxBackups = 3

	// 安全配置
	config.Security.EnableHTTPS = false
	config.Security.CertFile = ""
	config.Security.KeyFile = ""

	// BTC特定配置
	config.BTC.DefaultFeeRate = 10.0  // 10 sat/byte
	config.BTC.Confirmations = 6      // 6个确认
	config.BTC.AddressType = "legacy" // legacy地址类型

	return config
}

// overrideFromEnv 从环境变量覆盖配置
func (c *Config) overrideFromEnv() {
	// BTC网络配置
	if rpc := os.Getenv("BTC_MAINNET_RPC"); rpc != "" {
		c.Network.BTC.MainnetRPC = rpc
	}
	if rpc := os.Getenv("BTC_TESTNET_RPC"); rpc != "" {
		c.Network.BTC.TestnetRPC = rpc
	}
	if rpc := os.Getenv("BTC_REGTEST_RPC"); rpc != "" {
		c.Network.BTC.RegtestRPC = rpc
	}
	if networkType := os.Getenv("BTC_NETWORK_TYPE"); networkType != "" {
		c.Network.BTC.NetworkType = networkType
	}

	// API配置
	if port := os.Getenv("API_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.API.Port = p
		}
	}
	if host := os.Getenv("API_HOST"); host != "" {
		c.API.Host = host
	}

	// 日志配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		c.Log.Level = strings.ToLower(level)
	}

	// BTC特定配置
	if feeRate := os.Getenv("BTC_DEFAULT_FEE_RATE"); feeRate != "" {
		if f, err := strconv.ParseFloat(feeRate, 64); err == nil {
			c.BTC.DefaultFeeRate = f
		}
	}
	if confirmations := os.Getenv("BTC_CONFIRMATIONS"); confirmations != "" {
		if conf, err := strconv.Atoi(confirmations); err == nil {
			c.BTC.Confirmations = conf
		}
	}
	if addressType := os.Getenv("BTC_ADDRESS_TYPE"); addressType != "" {
		c.BTC.AddressType = addressType
	}
}

// GetBTCNetworkRPC 根据网络类型获取BTC RPC地址
// 参数: networkType - 网络类型
// 返回: RPC地址
func (c *Config) GetBTCNetworkRPC(networkType string) string {
	switch strings.ToLower(networkType) {
	case "mainnet":
		return c.Network.BTC.MainnetRPC
	case "testnet":
		return c.Network.BTC.TestnetRPC
	case "regtest":
		return c.Network.BTC.RegtestRPC
	default:
		return c.Network.BTC.MainnetRPC
	}
}

// GetCurrentNetworkRPC 获取当前配置的网络RPC地址
// 返回: 当前网络的RPC地址
func (c *Config) GetCurrentNetworkRPC() string {
	return c.GetBTCNetworkRPC(c.Network.BTC.NetworkType)
}
