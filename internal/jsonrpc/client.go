// Package jsonrpc 提供JSON-RPC 2.0客户端
package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
)

// Client JSON-RPC 2.0客户端
type Client struct {
	httpClient *http.Client
	url        string
	username   string
	password   string
	requestID  int64
}

// Request JSON-RPC请求
type Request struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// Response JSON-RPC响应
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error JSON-RPC错误
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// NewClient 创建新的JSON-RPC客户端
func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{
		httpClient: httpClient,
		url:        url,
	}
}

// SetAuth 设置认证信息
func (c *Client) SetAuth(username, password string) {
	c.username = username
	c.password = password
}

// Call 执行JSON-RPC调用
func (c *Client) Call(ctx context.Context, method string, params []interface{}) (*Response, error) {
	request := &Request{
		JSONRPC: "2.0",
		ID:      c.nextID(),
		Method:  method,
		Params:  params,
	}

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 设置认证
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// 解析响应
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查JSON-RPC错误
	if response.Error != nil {
		return nil, response.Error
	}

	return &response, nil
}

// BatchCall 执行批量JSON-RPC调用
func (c *Client) BatchCall(ctx context.Context, calls []Request) ([]Response, error) {
	// TODO: 实现批量调用
	return nil, nil
}

// nextID 生成下一个请求ID
func (c *Client) nextID() interface{} {
	id := atomic.AddInt64(&c.requestID, 1)
	return strconv.FormatInt(id, 10)
}
