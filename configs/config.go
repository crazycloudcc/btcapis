package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/crazycloudcc/btcapis"
)

type Config struct {
	TxID       string        `json:"txid"`
	TimeoutSec int           `json:"timeout_sec"`
	Bitcoind   *BitcoindConf `json:"bitcoind,omitempty"`
	Mempool    *MempoolConf  `json:"mempool,omitempty"`
}

type BitcoindConf struct {
	URL  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type MempoolConf struct {
	BaseURL string `json:"base_url"`
}

// LoadConfig 读取 JSON 文件，做最小校验与默认值填充
func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.TimeoutSec <= 0 {
		c.TimeoutSec = 8
	}
	if c.Bitcoind == nil && (c.Mempool == nil || c.Mempool.BaseURL == "") {
		return nil, fmt.Errorf("no backend configured: set bitcoind or mempool in %s", path)
	}
	return &c, nil
}

// BuildClient 根据配置构建 btcapis.Client
func (c *Config) BuildClient() *btcapis.Client {
	opts := make([]btcapis.Option, 0, 2)
	if c.Bitcoind != nil && c.Bitcoind.URL != "" {
		opts = append(opts, btcapis.WithBitcoindRPC(c.Bitcoind.URL, c.Bitcoind.User, c.Bitcoind.Pass))
	}
	if c.Mempool != nil && c.Mempool.BaseURL != "" {
		opts = append(opts, btcapis.WithMempoolSpace(c.Mempool.BaseURL))
	}
	return btcapis.New(opts...)
}

// Ctx 返回带超时的 context
func (c *Config) Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(c.TimeoutSec)*time.Second)
}
