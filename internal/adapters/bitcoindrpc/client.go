// Package bitcoindrpc 提供Bitcoin Core JSON-RPC客户端
package bitcoindrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

// ===== 内部 JSON-RPC =====

func (c *Client) rpcCall(ctx context.Context, method string, params []any, out any) error {
	c.idSeed++
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
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, &buf)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		httpReq.SetBasicAuth(c.user, c.pass)
	}

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("bitcoind rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if out != nil {
		return json.Unmarshal(rpcResp.Result, out)
	}
	return nil
}
