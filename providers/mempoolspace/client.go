// Package mempoolspace 提供mempool.space REST API客户端
package mempoolspace

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/btcapis/internal/httpx"
	"github.com/yourusername/btcapis/types"
)

// Client mempool.space REST客户端
type Client struct {
	httpClient   *httpx.Client
	baseURL      string
	config       *Config
	capabilities *types.Capabilities
}

// Config mempool.space配置
type Config struct {
	BaseURL    string        `json:"base_url"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	RateLimit  int           `json:"rate_limit"`
	Network    types.Network `json:"network"`
}

// NewClient 创建新的mempool.space客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// 创建HTTP客户端
	httpClient := httpx.NewClient(
		httpx.WithTimeout(config.Timeout),
		httpx.WithRetry(config.MaxRetries, 1*time.Second),
		httpx.WithRateLimit(config.RateLimit),
	)

	client := &Client{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		config:     config,
	}

	// 探测能力
	if err := client.detectCapabilities(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to detect capabilities: %w", err)
	}

	return client, nil
}

// detectCapabilities 探测后端能力
func (c *Client) detectCapabilities(ctx context.Context) error {
	// TODO: 实现能力探测
	c.capabilities = &types.Capabilities{
		HasChainReader:        true,
		HasBroadcaster:        false, // mempool.space不支持广播
		HasFeeEstimator:       true,
		HasMempoolView:        true,
		Network:               c.config.Network,
		SupportsSegWit:        true,
		SupportsTaproot:       true,
		MaxConcurrentRequests: 20,
		RequestTimeout:        int(c.config.Timeout.Seconds()),
		RateLimit:             200,
		ProvidesConfirmedData: true,
		ProvidesMempoolData:   true,
		DataFreshness:         10, // 10秒延迟
	}
	return nil
}

// Name 获取后端名称
func (c *Client) Name() string {
	return "mempool-space"
}

// IsHealthy 检查后端健康状态
func (c *Client) IsHealthy(ctx context.Context) bool {
	// TODO: 实现健康检查
	return true
}

// Capabilities 获取后端能力
func (c *Client) Capabilities(ctx context.Context) (*types.Capabilities, error) {
	return c.capabilities, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c.httpClient != nil {
		return c.httpClient.Close()
	}
	return nil
}
