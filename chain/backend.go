// Package chain 定义后端接口和路由策略
package chain

import (
	"context"

	"github.com/yourusername/btcapis/types"
)

// ChainReader 定义链上数据读取接口
type ChainReader interface {
	// GetRawTransaction 获取原始交易数据
	GetRawTransaction(ctx context.Context, txid string) ([]byte, error)

	// GetBlockHash 根据高度获取区块哈希
	GetBlockHash(ctx context.Context, height int64) (string, error)

	// GetBlockHeader 获取区块头
	GetBlockHeader(ctx context.Context, hash string) ([]byte, error)

	// GetBlock 获取完整区块
	GetBlock(ctx context.Context, hash string) ([]byte, error)

	// GetUTXO 获取UTXO信息
	GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error)

	// GetBlockHeight 获取当前区块高度
	GetBlockHeight(ctx context.Context) (int64, error)
}

// Broadcaster 定义交易广播接口
type Broadcaster interface {
	// Broadcast 广播交易
	Broadcast(ctx context.Context, rawtx []byte) (string, error)
}

// FeeEstimator 定义手续费估算接口
type FeeEstimator interface {
	// EstimateFeeRate 估算手续费率
	EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error)

	// GetFeeHistory 获取手续费历史
	GetFeeHistory(ctx context.Context, targetBlocks []int) (*types.FeeHistory, error)
}

// MempoolView 定义内存池视图接口
type MempoolView interface {
	// GetRawMempool 获取原始内存池
	GetRawMempool(ctx context.Context) ([]string, error)

	// TxInMempool 检查交易是否在内存池中
	TxInMempool(ctx context.Context, txid string) (bool, error)

	// GetMempoolInfo 获取内存池信息
	GetMempoolInfo(ctx context.Context) (*types.MempoolInfo, error)

	// GetMempoolEntry 获取内存池条目
	GetMempoolEntry(ctx context.Context, txid string) (*types.MempoolEntry, error)
}

// Backend 定义完整的后端接口
type Backend interface {
	ChainReader
	Broadcaster
	FeeEstimator
	MempoolView

	// Capabilities 获取后端能力
	Capabilities(ctx context.Context) (types.Capabilities, error)

	// Name 获取后端名称
	Name() string

	// IsHealthy 检查后端健康状态
	IsHealthy(ctx context.Context) bool
}

// BackendFactory 定义后端工厂接口
type BackendFactory interface {
	// Create 创建后端实例
	Create(config interface{}) (Backend, error)

	// Name 获取工厂名称
	Name() string
}
