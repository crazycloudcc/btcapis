// Package bitcoindrpc 提供Bitcoin Core JSON-RPC客户端
package bitcoindrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/btcapis/internal/httpx"
	"github.com/yourusername/btcapis/internal/jsonrpc"
	"github.com/yourusername/btcapis/types"
)

// Client Bitcoin Core RPC客户端
type Client struct {
	httpClient   *httpx.Client
	rpcClient    *jsonrpc.Client
	config       *Config
	capabilities *types.Capabilities
}

// Config Bitcoin Core RPC配置
type Config struct {
	URL      string        `json:"url"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Timeout  time.Duration `json:"timeout"`
	Network  types.Network `json:"network"`

	// 高级配置
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
	RateLimit  int           `json:"rate_limit"`
	BatchSize  int           `json:"batch_size"`
}

// NewClient 创建新的Bitcoin Core RPC客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// 创建HTTP客户端
	httpClient := httpx.NewClient(
		httpx.WithTimeout(config.Timeout),
		httpx.WithRetry(config.MaxRetries, config.RetryDelay),
		httpx.WithRateLimit(config.RateLimit),
	)

	// 创建JSON-RPC客户端
	rpcClient := jsonrpc.NewClient(httpClient, config.URL)

	client := &Client{
		httpClient: httpClient,
		rpcClient:  rpcClient,
		config:     config,
	}

	// 设置认证
	client.rpcClient.SetAuth(config.Username, config.Password)

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
		HasBroadcaster:        true,
		HasFeeEstimator:       true,
		HasMempoolView:        true,
		Network:               c.config.Network,
		SupportsSegWit:        true,
		SupportsTaproot:       true,
		MaxConcurrentRequests: 10,
		RequestTimeout:        int(c.config.Timeout.Seconds()),
		RateLimit:             100,
		ProvidesConfirmedData: true,
		ProvidesMempoolData:   true,
		DataFreshness:         0,
	}
	return nil
}

// Name 获取后端名称
func (c *Client) Name() string {
	return "bitcoind-rpc"
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
