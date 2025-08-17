// Package mempoolspace 提供mempool.space REST API客户端
package mempoolspace

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/crazycloudcc/btcapis/chain"
	"github.com/crazycloudcc/btcapis/types"
)

type Client struct {
	base *url.URL
	http *http.Client
}

type Option func(*Client)

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.http = h }
}

func New(baseURL string, opts ...Option) *Client {
	u, _ := url.Parse(baseURL)
	c := &Client{
		base: u,
		http: &http.Client{Timeout: 8 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) Capabilities(ctx context.Context) (chain.Capabilities, error) {
	return chain.Capabilities{
		HasMempool:     true,
		HasFeeEstimate: false,         // 如需可接 /api/v1/fees/recommended
		Network:        types.Mainnet, // 简化：默认主网
	}, nil
}

func (c *Client) GetTx(ctx context.Context, txid string) (*types.Tx, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/tx/", txid)
	var dto TxDTO
	if err := c.getJSON(ctx, u.String(), &dto); err != nil {
		return nil, err
	}
	return mapTxDTO(dto), nil
}

func (c *Client) GetRawTransaction(ctx context.Context, txid string) ([]byte, error) {
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

func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error) {
	// 如需可实现 /api/v1/fees/recommended；此处先不实现
	return 0, fmt.Errorf("mempool: fee estimate not implemented")
}

// 其它 ChainReader 方法（此 provider 非权威链数据，先不实现）
func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	return "", fmt.Errorf("mempool: not implemented")
}
func (c *Client) GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
	return nil, fmt.Errorf("mempool: not implemented")
}
func (c *Client) GetBlock(ctx context.Context, hash string) ([]byte, error) {
	return nil, fmt.Errorf("mempool: not implemented")
}
func (c *Client) GetUTXO(ctx context.Context, op types.OutPoint) (*types.UTXO, error) {
	return nil, fmt.Errorf("mempool: not implemented")
}

// GetRawMempool：mempool.space 没有公开“全量 txids 列表”的稳定接口，这里先占位返回未实现。
// 需要时可以后续改成分页/条件查询。
func (c *Client) GetRawMempool(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("mempool: GetRawMempool not implemented")
}

func (c *Client) TxInMempool(ctx context.Context, txid string) (bool, error) {
	// 若要实现，可调用 /api/tx/{txid}：存在且未确认则视为在 mempool 中。
	return false, fmt.Errorf("mempool: TxInMempool not implemented")
}

// ===== HTTP helpers =====
func (c *Client) getJSON(ctx context.Context, url string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s status=%d body=%s", url, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) getBytes(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET %s status=%d body=%s", url, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return io.ReadAll(resp.Body)
}
