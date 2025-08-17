// Package btcapis 链上数据门面
package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/chain"
	"github.com/crazycloudcc/btcapis/types"
)

// chainFacade 提供链上数据查询的功能接口
type chainFacade struct {
	router *chain.Router
}

// GetBlockHeight 获取当前区块高度
func (c *chainFacade) GetBlockHeight(ctx context.Context) (int64, error) {
	// TODO: 实现区块高度查询
	return 0, nil
}

// GetBlockHash 根据高度获取区块哈希
func (c *chainFacade) GetBlockHash(ctx context.Context, height int64) (string, error) {
	// TODO: 实现区块哈希查询
	return "", nil
}

// GetBlockHeader 获取区块头
func (c *chainFacade) GetBlockHeader(ctx context.Context, hash string) (*types.BlockHeader, error) {
	// TODO: 实现区块头查询
	return nil, nil
}

// GetBlock 获取完整区块
func (c *chainFacade) GetBlock(ctx context.Context, hash string) (*types.Block, error) {
	// TODO: 实现区块查询
	return nil, nil
}

// GetTransaction 获取交易详情
func (c *chainFacade) GetTransaction(ctx context.Context, txid string) (*types.Transaction, error) {
	// TODO: 实现交易查询
	return nil, nil
}

// GetUTXO 获取UTXO信息
func (c *chainFacade) GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error) {
	// TODO: 实现UTXO查询
	return nil, nil
}

// GetMempool 获取内存池信息
func (c *chainFacade) GetMempool(ctx context.Context) (*types.MempoolInfo, error) {
	// TODO: 实现内存池查询
	return nil, nil
}

// EstimateFee 估算手续费
func (c *chainFacade) EstimateFee(ctx context.Context, targetBlocks int) (*types.FeeEstimate, error) {
	// TODO: 实现手续费估算
	return nil, nil
}
