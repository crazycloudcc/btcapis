// Package httpx 提供增强的HTTP客户端功能
package httpx

import (
	"context"
	"net/http"
	"time"
)

// Client 增强的HTTP客户端
type Client struct {
	httpClient *http.Client
	options    *Options
}

// Options HTTP客户端选项
type Options struct {
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
	RateLimit  int
	UserAgent  string
	Proxy      string
	Transport  http.RoundTripper
}

// Option 定义配置选项函数
type Option func(*Options)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithRetry 设置重试配置
func WithRetry(maxRetries int, retryDelay time.Duration) Option {
	return func(o *Options) {
		o.MaxRetries = maxRetries
		o.RetryDelay = retryDelay
	}
}

// WithRateLimit 设置速率限制
func WithRateLimit(rateLimit int) Option {
	return func(o *Options) {
		o.RateLimit = rateLimit
	}
}

// WithUserAgent 设置用户代理
func WithUserAgent(userAgent string) Option {
	return func(o *Options) {
		o.UserAgent = userAgent
	}
}

// WithProxy 设置代理
func WithProxy(proxy string) Option {
	return func(o *Options) {
		o.Proxy = proxy
	}
}

// WithTransport 设置传输层
func WithTransport(transport http.RoundTripper) Option {
	return func(o *Options) {
		o.Transport = transport
	}
}

// NewClient 创建新的HTTP客户端
func NewClient(opts ...Option) *Client {
	options := &Options{
		Timeout:    30 * time.Second,
		MaxRetries: 0,
		RetryDelay: 1 * time.Second,
		RateLimit:  0,
		UserAgent:  "btcapis/1.0",
	}

	for _, opt := range opts {
		opt(options)
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout: options.Timeout,
	}

	if options.Transport != nil {
		httpClient.Transport = options.Transport
	}

	return &Client{
		httpClient: httpClient,
		options:    options,
	}
}

// Do 执行HTTP请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// TODO: 实现重试、速率限制等逻辑
	return c.httpClient.Do(req)
}

// Get 执行GET请求
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post 执行POST请求
func (c *Client) Post(ctx context.Context, url, contentType string, body interface{}) (*http.Response, error) {
	// TODO: 实现POST请求
	return nil, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	// TODO: 实现资源清理
	return nil
}
