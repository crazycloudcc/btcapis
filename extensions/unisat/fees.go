// Package unisat 提供 Unisat Open API 费率扩展
package unisat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/crazycloudcc/btcapis/types"
)

const defaultBaseURL = "https://open-api.unisat.io"

// Config Unisat API 配置
type Config struct {
	BaseURL string
	APIKey  string
	Timeout int // 秒
}

// Client Unisat 费率客户端
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// New 创建 Unisat 客户端
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("unisat api key is required")
	}
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}, nil
}

func (c *Client) Name() string { return "unisat" }

type unisatResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type unisatRecommendedFees struct {
	FastestFee  int64 `json:"fastestFee"`
	HalfHourFee int64 `json:"halfHourFee"`
	HourFee     int64 `json:"hourFee"`
	EconomyFee  int64 `json:"economyFee"`
	MinimumFee  int64 `json:"minimumFee"`
}

// GetRecommendedFees 获取推荐费率（sat/vB）
func (c *Client) GetRecommendedFees(ctx context.Context) (types.RecommendedFees, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/indexer/fees/recommended", nil)
	if err != nil {
		return types.RecommendedFees{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return types.RecommendedFees{}, fmt.Errorf("unisat request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.RecommendedFees{}, err
	}
	if resp.StatusCode >= 300 {
		return types.RecommendedFees{}, fmt.Errorf("unisat request failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var wrapper unisatResponse
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return types.RecommendedFees{}, fmt.Errorf("unisat decode response: %w", err)
	}
	if wrapper.Code != 0 {
		return types.RecommendedFees{}, fmt.Errorf("unisat api error: %s", wrapper.Msg)
	}
	var fees unisatRecommendedFees
	if err := json.Unmarshal(wrapper.Data, &fees); err != nil {
		return types.RecommendedFees{}, fmt.Errorf("unisat decode fees: %w", err)
	}
	return types.RecommendedFees{
		FastestFee:  float64(fees.FastestFee),
		HalfHourFee: float64(fees.HalfHourFee),
		HourFee:     float64(fees.HourFee),
		EconomyFee:  float64(fees.EconomyFee),
		MinimumFee:  float64(fees.MinimumFee),
	}, nil
}
