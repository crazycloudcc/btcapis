package mempoolapis

import (
	"context"
	"path"
)

// 获取地址余额
func (c *Client) AddressGetBalance(ctx context.Context, addr string) (float64, float64, error) {
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
	return float64(confirmed), float64(mempool), nil
}

// 获取地址 UTXO
func (c *Client) AddressGetUTXOs(ctx context.Context, addr string) ([]UTXODTO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/address/", addr, "/utxo")
	var dtos []UTXODTO
	if err := c.getJSON(ctx, u.String(), &dtos); err != nil {
		return nil, err
	}
	return dtos, nil
}
