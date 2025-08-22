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

func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/tx/", txid, "hex")
	b, err := c.getBytes(ctx, u.String())
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(b)))
}

func (c *Client) Broadcast(ctx context.Context, rawtx []byte) (string, error) {
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

// 获取地址余额
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (int64, int64, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/address/", addr)
	var dto struct {
		ChainStats struct {
			Funded int64 `json:"funded_txo_sum"`
			Spent  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
		MempoolStats struct {
			Funded int64 `json:"funded_txo_sum"`
			Spent  int64 `json:"spent_txo_sum"`
		} `json:"mempool_stats"`
	}
	if err := c.getJSON(ctx, u.String(), &dto); err != nil {
		return 0, 0, err
	}
	confirmed := dto.ChainStats.Funded - dto.ChainStats.Spent
	mempool := dto.MempoolStats.Funded - dto.MempoolStats.Spent
	return confirmed, mempool, nil
}

// 获取地址 UTXO
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]UTXODTO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/address/", addr, "/utxo")
	var dtos []UTXODTO
	if err := c.getJSON(ctx, u.String(), &dtos); err != nil {
		return nil, err
	}
	return dtos, nil
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
