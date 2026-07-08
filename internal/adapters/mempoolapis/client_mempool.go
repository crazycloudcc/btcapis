package mempoolapis

import (
	"context"
	"net/url"
	"path"
)

// GetMempoolStats 获取内存池统计
func (c *Client) GetMempoolStats(ctx context.Context) (*MempoolStatsDTO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/mempool")
	var stats MempoolStatsDTO
	if err := c.getJSON(ctx, u.String(), &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetTxStatus 获取交易确认状态
func (c *Client) GetTxStatus(ctx context.Context, txid string) (*TxStatusDTO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/tx/", url.PathEscape(txid), "/status")
	var status TxStatusDTO
	if err := c.getJSON(ctx, u.String(), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// GetRecommendedFees 获取推荐费率原始数据
func (c *Client) GetRecommendedFees(ctx context.Context) (map[string]float64, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/v1/fees/recommended")
	var fees map[string]float64
	if err := c.getJSON(ctx, u.String(), &fees); err != nil {
		return nil, err
	}
	return fees, nil
}

// MempoolStatsDTO mempool 统计
type MempoolStatsDTO struct {
	Count        int64       `json:"count"`
	Vsize        int64       `json:"vsize"`
	TotalFee     int64       `json:"total_fee"`
	FeeHistogram [][]float64 `json:"fee_histogram"`
}

// TxStatusDTO 交易状态
type TxStatusDTO struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int64  `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}
