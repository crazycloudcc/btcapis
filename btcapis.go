// Package btcapis 提供比特币区块链API的统一接口
//
// 该包采用端口/适配器/门面架构，支持多种后端服务：
// - Bitcoin Core RPC
// - mempool.space REST API
// - Electrum TCP
//
// 主要特性：
// - 多后端聚合与故障转移
// - 统一的错误处理模型
// - 智能路由和负载均衡
// - 与btcd生态的兼容性
package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/chain"
	"github.com/crazycloudcc/btcapis/types"

	// 若这些 provider 还没实现，可先注释掉对应 WithXxx
	"github.com/crazycloudcc/btcapis/providers/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/providers/mempoolspace"
)

type (
	AddressModule struct{}
	ScriptModule  struct{}
	TxModule      struct{}
	PSBTModule    struct{}
)

var (
	Address AddressModule
	Script  ScriptModule
	Tx      TxModule
	PSBT    PSBTModule
)

type Client struct {
	router    *chain.Router
	primaries []chain.Backend
	fallbacks []chain.Backend
}

type option func(*Client)

func newClient(opts ...option) *Client {
	c := &Client{}
	for _, o := range opts {
		o(c)
	}
	c.router = chain.NewRouter(c.primaries, c.fallbacks)
	return c
}

// BuildClient 根据配置构建 Client
func BuildClient(network string, bitcoindUrl, bitcoindUser, bitcoindPass, mempoolBaseUrl string) *Client {
	types.SetCurrentNetwork(network)
	opts := make([]option, 0, 2)
	if bitcoindUrl != "" {
		opts = append(opts, WithBitcoindRPC(bitcoindUrl, bitcoindUser, bitcoindPass))
	}
	if mempoolBaseUrl != "" {
		opts = append(opts, WithMempoolSpace(mempoolBaseUrl))
	}
	return newClient(opts...)
}

func WithBitcoindRPC(url, user, pass string, opts ...bitcoindrpc.Option) option {
	return func(c *Client) {
		b := bitcoindrpc.New(url, user, pass, opts...)
		c.primaries = append(c.primaries, b)
	}
}

func WithMempoolSpace(baseURL string, opts ...mempoolspace.Option) option {
	return func(c *Client) {
		m := mempoolspace.New(baseURL, opts...)
		c.fallbacks = append(c.fallbacks, m)
	}
}

// ==== 高层转发 ====

func (c *Client) GetTransaction(ctx context.Context, txid string) (*types.Tx, error) {
	return c.router.GetTransaction(ctx, txid)
}

func (c *Client) GetRawTransaction(ctx context.Context, txid string) ([]byte, error) {
	return c.router.GetRawTransaction(ctx, txid)
}

func (c *Client) EstimateFeeRate(ctx context.Context, target int) (float64, error) {
	return c.router.EstimateFeeRate(ctx, target)
}

func (c *Client) Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	return c.router.Broadcast(ctx, rawtx)
}

// // Ctx 返回带超时的 context
// func Ctx() (context.Context, context.CancelFunc) {
// 	return context.WithTimeout(context.Background(), 30*time.Second)
// }
