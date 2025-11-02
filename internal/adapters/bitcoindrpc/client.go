// Package bitcoindrpc 提供Bitcoin Core JSON-RPC客户端
package bitcoindrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/crazycloudcc/btcapis/pkg/logger"
)

type Client struct {
	url    string
	user   string
	pass   string
	http   *http.Client
	idSeed int
}

func New(url, user, pass string, timeout int) *Client {
	return &Client{
		url:  url,
		user: user,
		pass: pass,
		http: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func (c *Client) rpcCall(ctx context.Context, method string, params []any, out any) error {
	// startTime := time.Now()
	c.idSeed++

	// 记录请求开始
	// logger.Debug("[DEBUG] Bitcoin RPC 请求开始 - Method: %s, ID: %d, Params: %+v", method, c.idSeed, params)

	req := struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Method  string `json:"method"`
		Params  []any  `json:"params"`
	}{
		JSONRPC: "2.0",
		ID:      c.idSeed,
		Method:  method,
		Params:  params,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		logger.Error("[ERROR] JSON 编码失败: %v", err)
		return err
	}

	// 记录请求体内容
	// logger.Debug("[DEBUG] RPC 请求体: %s", buf.String())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, &buf)
	if err != nil {
		logger.Error("[ERROR] 创建 HTTP 请求失败: %v", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		httpReq.SetBasicAuth(c.user, c.pass)
		// logger.Debug("[DEBUG] 使用认证 - User: %s", c.user)
	}

	// logger.Debug("[DEBUG] 发送 HTTP 请求到: %s", c.url)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		logger.Error("[ERROR] HTTP 请求执行失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 记录响应状态
	// logger.Debug("[DEBUG] HTTP 响应状态: %s, StatusCode: %d", resp.Status, resp.StatusCode)

	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		ID int `json:"id"`
	}

	// 读取响应体用于日志记录
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("[ERROR] 读取响应体失败: %v", err)
		return err
	}

	// logger.Debug("[DEBUG] RPC 响应体: %s", string(respBody))

	// 重新创建 reader 用于 JSON 解码
	reader := bytes.NewReader(respBody)
	if err := json.NewDecoder(reader).Decode(&rpcResp); err != nil {
		logger.Error("[ERROR] JSON 响应解码失败: %v", err)
		return err
	}

	// 记录 RPC 响应详情
	// logger.Debug("[DEBUG] RPC 响应 - ID: %d, HasError: %v, HasResult: %v",
	// 	rpcResp.ID, rpcResp.Error != nil, len(rpcResp.Result) > 0)

	if rpcResp.Error != nil {
		logger.Error("[ERROR] Bitcoin RPC 错误 - Code: %d, Message: %s",
			rpcResp.Error.Code, rpcResp.Error.Message)
		return fmt.Errorf("bitcoind rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if out != nil {
		if err := json.Unmarshal(rpcResp.Result, out); err != nil {
			logger.Error("[ERROR] 结果反序列化失败: %v", err)
			return err
		}
		// logger.Debug("[DEBUG] 结果反序列化成功")
	}

	// duration := time.Since(startTime)
	// logger.Debug("[DEBUG] Bitcoin RPC 请求完成 - Method: %s, 耗时: %v", method, duration)

	return nil
}

// createrawtransaction接口必须使用any, 不能使用[]any.
func (c *Client) rpcCallWithAny(ctx context.Context, method string, params any, out any) error {
	// startTime := time.Now()
	c.idSeed++

	// 记录请求开始
	// logger.Debug("[DEBUG] Bitcoin RPC 请求开始 - Method: %s, ID: %d, Params: %+v", method, c.idSeed, params)

	req := struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Method  string `json:"method"`
		Params  any    `json:"params"`
	}{
		JSONRPC: "1.0",
		ID:      c.idSeed,
		Method:  method,
		Params:  params,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		logger.Error("[ERROR] JSON 编码失败: %v", err)
		return err
	}

	// 记录请求体内容
	// logger.Debug("[DEBUG] RPC 请求体: %s", buf.String())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, &buf)
	if err != nil {
		logger.Error("[ERROR] 创建 HTTP 请求失败: %v", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		httpReq.SetBasicAuth(c.user, c.pass)
		// logger.Debug("[DEBUG] 使用认证 - User: %s", c.user)
	}

	// logger.Debug("[DEBUG] 发送 HTTP 请求到: %s", c.url)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		logger.Error("[ERROR] HTTP 请求执行失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 记录响应状态
	// logger.Debug("[DEBUG] HTTP 响应状态: %s, StatusCode: %d", resp.Status, resp.StatusCode)

	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		ID int `json:"id"`
	}

	// 读取响应体用于日志记录
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("[ERROR] 读取响应体失败: %v", err)
		return err
	}

	// logger.Debug("[DEBUG] RPC 响应体: %s", string(respBody))

	// 重新创建 reader 用于 JSON 解码
	reader := bytes.NewReader(respBody)
	if err := json.NewDecoder(reader).Decode(&rpcResp); err != nil {
		logger.Error("[ERROR] JSON 响应解码失败: %v", err)
		return err
	}

	// 记录 RPC 响应详情
	// logger.Debug("[DEBUG] RPC 响应 - ID: %d, HasError: %v, HasResult: %v",
	// 	rpcResp.ID, rpcResp.Error != nil, len(rpcResp.Result) > 0)

	if rpcResp.Error != nil {
		logger.Error("[ERROR] Bitcoin RPC 错误 - Code: %d, Message: %s",
			rpcResp.Error.Code, rpcResp.Error.Message)
		return fmt.Errorf("bitcoind rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if out != nil {
		if err := json.Unmarshal(rpcResp.Result, out); err != nil {
			logger.Error("[ERROR] 结果反序列化失败: %v", err)
			return err
		}
		// logger.Debug("[DEBUG] 结果反序列化成功")
	}

	// duration := time.Since(startTime)
	// logger.Debug("[DEBUG] Bitcoin RPC 请求完成 - Method: %s, 耗时: %v", method, duration)

	return nil
}
