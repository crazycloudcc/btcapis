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
)

// Client 是btcapis的主要客户端，聚合多个后端服务
type Client struct {
	router *chain.Router
}

// Option 定义客户端配置选项
type Option func(*Client)

// New 创建新的btcapis客户端
func New(opts ...Option) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithBitcoindRPC 配置Bitcoin Core RPC后端
func WithBitcoindRPC(url, user, pass string, opts ...BitcoindOption) Option {
	return func(c *Client) {
		// TODO: 实现Bitcoin Core RPC后端配置
	}
}

// WithMempoolSpace 配置mempool.space REST后端
func WithMempoolSpace(baseURL string, opts ...MempoolOption) Option {
	return func(c *Client) {
		// TODO: 实现mempool.space后端配置
	}
}

// WithElectrum 配置Electrum TCP后端
func WithElectrum(addr string, opts ...ElectrumOption) Option {
	return func(c *Client) {
		// TODO: 实现Electrum后端配置
	}
}

// 网络类型常量
const (
	Mainnet = types.Network("mainnet")
	Testnet = types.Network("testnet")
	Signet  = types.Network("signet")
	Regtest = types.Network("regtest")
)

// 地址模块门面
var Address = &addressFacade{}

// 脚本模块门面
var Script = &scriptFacade{}

// 交易模块门面
var Tx = &txFacade{}

// PSBT模块门面
var PSBT = &psbtFacade{}

// 链上数据查询方法
func (c *Client) GetRawTransaction(ctx context.Context, txid string) ([]byte, error) {
	// TODO: 实现交易查询
	return nil, nil
}

func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	// TODO: 实现区块哈希查询
	return "", nil
}

func (c *Client) GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
	// TODO: 实现区块头查询
	return nil, nil
}

func (c *Client) GetBlock(ctx context.Context, hash string) ([]byte, error) {
	// TODO: 实现区块查询
	return nil, nil
}

func (c *Client) GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error) {
	// TODO: 实现UTXO查询
	return nil, nil
}

// 交易广播
func (c *Client) Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	// TODO: 实现交易广播
	return "", nil
}

// 费率估算
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error) {
	// TODO: 实现费率估算
	return 0, nil
}

// Mempool查询
func (c *Client) GetRawMempool(ctx context.Context) ([]string, error) {
	// TODO: 实现mempool查询
	return nil, nil
}

func (c *Client) TxInMempool(ctx context.Context, txid string) (bool, error) {
	// TODO: 实现交易mempool状态查询
	return false, nil
}

// 后端能力探测
func (c *Client) Capabilities(ctx context.Context) (types.Capabilities, error) {
	// TODO: 实现能力探测
	return types.Capabilities{}, nil
}

// 配置选项类型定义
type BitcoindOption func(interface{})
type MempoolOption func(interface{})
type ElectrumOption func(interface{})
