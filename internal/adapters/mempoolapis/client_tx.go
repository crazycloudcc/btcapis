package mempoolapis

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

// 获取交易的原始数据，返回二进制格式
func (c *Client) TxGetRaw(ctx context.Context, txid string) ([]byte, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/tx/", txid, "hex")
	b, err := c.getBytes(ctx, u.String())
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(b)))
}

// 广播交易，返回交易ID
func (c *Client) TxBroadcast(ctx context.Context, rawtx []byte) (string, error) {
	// mempool.space 支持 POST /api/tx，body 为 hex
	u := *c.base
	u.Path = path.Join(u.Path, "/api/tx")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(hex.EncodeToString(rawtx)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("mempool POST /api/tx status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	txid, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(txid)), nil
}

// 估算交易费率
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (*FeeRateDTO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/v1/fees/recommended")
	var dto FeeRateDTO
	if err := c.getJSON(ctx, u.String(), &dto); err != nil {
		return nil, err
	}

	fmt.Printf("mempooldto: %+v\n", dto)

	return &dto, nil
}
