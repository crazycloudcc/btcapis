// Package mempoolspace 提供mempool.space REST API客户端
package mempoolapis

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

	"github.com/crazycloudcc/btcapis/src/types"
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

func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/tx/", txid, "hex")
	b, err := c.getBytes(ctx, u.String())
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(b)))
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
