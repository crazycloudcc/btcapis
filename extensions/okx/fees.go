// Package okx 提供 OKX Web3 API 费率扩展
package okx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/crazycloudcc/btcapis/internal/utils"
	"github.com/crazycloudcc/btcapis/types"
)

const defaultBaseURL = "https://www.okx.com"

// Config OKX API 配置
type Config struct {
	BaseURL    string
	APIKey     string
	SecretKey  string
	Passphrase string
	Timeout    int // 秒
}

// Client OKX 费率客户端
type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	passphrase string
	http       *http.Client
}

// New 创建 OKX 客户端
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" || cfg.SecretKey == "" || cfg.Passphrase == "" {
		return nil, fmt.Errorf("okx api credentials are required")
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
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     cfg.APIKey,
		secretKey:  cfg.SecretKey,
		passphrase: cfg.Passphrase,
		http:       &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}, nil
}

func (c *Client) Name() string { return "okx" }

type okxResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type okxFeeRate struct {
	FastestFee  string `json:"fastestFee"`
	HalfHourFee string `json:"halfHourFee"`
	HourFee     string `json:"hourFee"`
	EconomyFee  string `json:"economyFee"`
	MinimumFee  string `json:"minimumFee"`
}

func (c *Client) sign(timestamp, method, requestPath, body string) string {
	preHash := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(preHash))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// GetRecommendedFees 获取推荐费率（sat/vB）
func (c *Client) GetRecommendedFees(ctx context.Context) (types.RecommendedFees, error) {
	requestPath := "/api/v5/wallet/pre-transaction/gas-price?chainIndex=0"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+requestPath, nil)
	if err != nil {
		return types.RecommendedFees{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("OK-ACCESS-KEY", c.apiKey)
	req.Header.Set("OK-ACCESS-SIGN", c.sign(timestamp, "GET", requestPath, ""))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", c.passphrase)

	resp, err := c.http.Do(req)
	if err != nil {
		return types.RecommendedFees{}, fmt.Errorf("okx request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.RecommendedFees{}, err
	}
	if resp.StatusCode >= 300 {
		return types.RecommendedFees{}, fmt.Errorf("okx request failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var wrapper okxResponse
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return types.RecommendedFees{}, fmt.Errorf("okx decode response: %w", err)
	}
	if wrapper.Code != "0" {
		return types.RecommendedFees{}, fmt.Errorf("okx api error: %s - %s", wrapper.Code, wrapper.Msg)
	}
	var feesList []okxFeeRate
	if err := json.Unmarshal(wrapper.Data, &feesList); err != nil {
		return types.RecommendedFees{}, fmt.Errorf("okx decode fees: %w", err)
	}
	if len(feesList) == 0 {
		return types.RecommendedFees{}, fmt.Errorf("okx api returned empty fees")
	}
	fees := feesList[0]
	return types.RecommendedFees{
		FastestFee:  utils.ParseFloatFee(fees.FastestFee),
		HalfHourFee: utils.ParseFloatFee(fees.HalfHourFee),
		HourFee:     utils.ParseFloatFee(fees.HourFee),
		EconomyFee:  utils.ParseFloatFee(fees.EconomyFee),
		MinimumFee:  utils.ParseFloatFee(fees.MinimumFee),
	}, nil
}
