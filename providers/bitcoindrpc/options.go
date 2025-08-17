// Package bitcoindrpc 配置选项
package bitcoindrpc

import (
	"time"

	"github.com/crazycloudcc/btcapis/types"
)

// Option 定义配置选项函数
type Option func(*Config)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithNetwork 设置网络类型
func WithNetwork(network types.Network) Option {
	return func(c *Config) {
		c.Network = network
	}
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
	}
}

// WithRetryDelay 设置重试延迟
func WithRetryDelay(retryDelay time.Duration) Option {
	return func(c *Config) {
		c.RetryDelay = retryDelay
	}
}

// WithRateLimit 设置速率限制
func WithRateLimit(rateLimit int) Option {
	return func(c *Config) {
		c.RateLimit = rateLimit
	}
}

// WithBatchSize 设置批处理大小
func WithBatchSize(batchSize int) Option {
	return func(c *Config) {
		c.BatchSize = batchSize
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Timeout:    30 * time.Second,
		Network:    types.NetworkMainnet,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
		RateLimit:  100,
		BatchSize:  100,
	}
}

// NewConfig 创建新配置
func NewConfig(url, username, password string, opts ...Option) *Config {
	config := DefaultConfig()
	config.URL = url
	config.Username = username
	config.Password = password

	for _, opt := range opts {
		opt(config)
	}

	return config
}
