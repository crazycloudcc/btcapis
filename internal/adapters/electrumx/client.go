// Package electrumx 提供ElectrumX JSON-RPC客户端
package electrumx

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
	http   *http.Client
	idSeed int
}

// New 创建ElectrumX客户端实例
// baseURL: ElectrumX服务器地址，例如 "http://localhost:50001"
// timeout: 请求超时时间（秒）
func New(baseURL string, timeout int) *Client {
	return &Client{
		url:  baseURL,
		http: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

// rpcCall 执行ElectrumX JSON-RPC调用
func (c *Client) rpcCall(ctx context.Context, method string, params []interface{}, out interface{}) error {
	c.idSeed++

	// 构建JSON-RPC请求
	req := struct {
		JSONRPC string        `json:"jsonrpc"`
		ID      int           `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}{
		JSONRPC: "2.0",
		ID:      c.idSeed,
		Method:  method,
		Params:  params,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		logger.Error("[ERROR] ElectrumX JSON 编码失败: %v", err)
		return err
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, &buf)
	if err != nil {
		logger.Error("[ERROR] 创建 HTTP 请求失败: %v", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.http.Do(httpReq)
	if err != nil {
		logger.Error("[ERROR] HTTP 请求执行失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		ID int `json:"id"`
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("[ERROR] 读取响应体失败: %v", err)
		return err
	}

	// 解码JSON响应
	reader := bytes.NewReader(respBody)
	if err := json.NewDecoder(reader).Decode(&rpcResp); err != nil {
		logger.Error("[ERROR] JSON 响应解码失败: %v", err)
		return err
	}

	// 检查RPC错误
	if rpcResp.Error != nil {
		logger.Error("[ERROR] ElectrumX RPC 错误 - Code: %d, Message: %s",
			rpcResp.Error.Code, rpcResp.Error.Message)
		return fmt.Errorf("electrumx rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	// 解析结果
	if out != nil {
		if err := json.Unmarshal(rpcResp.Result, out); err != nil {
			logger.Error("[ERROR] 结果反序列化失败: %v", err)
			return err
		}
	}

	return nil
}
