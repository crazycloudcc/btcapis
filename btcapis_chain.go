package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// // GetBlockHash 返回区块哈希.
// func GetBlockHash(ctx context.Context, height int64) (string, error) {
// 	return bitcoindrpc.GetBlockHash(ctx, height)
// }

// // GetBlockHeader 返回区块头.
// func GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
// 	return bitcoindrpc.GetBlockHeader(ctx, hash)
// }

// // GetBlock 返回区块.
// func GetBlock(ctx context.Context, hash string) ([]byte, error) {
// 	return bitcoindrpc.GetBlock(ctx, hash)
// }

// // GetUTXO 返回UTXO详细信息.(因为这个UTXO不是通过关联钱包查询, 所以暂时放在ChainReader分类)
// // func GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error) {
// // 	return bitcoindrpc.GetUTXO(ctx, outpoint)
// // }

// EstimateFeeRate 估计手续费.
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, float64, error) {
	return c.chainClient.EstimateFeeRate(ctx, targetBlocks)
}

// 查询 UTXO
func (c *Client) GetUTXO(ctx context.Context, hash [32]byte, index uint32) ([]byte, int64, error) {
	return c.chainClient.GetUTXO(ctx, hash, index)
}

// 获取节点区块数量
func (c *Client) GetBlockCount(ctx context.Context) (int, error) {
	return c.chainClient.GetBlockCount(ctx)
}

// 获取最新区块的hash
func (c *Client) GetBestBlockHash(ctx context.Context) (string, error) {
	return c.chainClient.GetBestBlockHash(ctx)
}

// 使用区块高度 查询区块哈希
func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	return c.chainClient.GetBlockHash(ctx, height)
}

// 使用区块block hash 查询区块头
func (c *Client) GetBlockHeader(ctx context.Context, hash string) (*types.ChainBlock, error) {
	return c.chainClient.GetBlockHeader(ctx, hash)
}

// 使用区块block hash 查询区块
func (c *Client) GetBlock(ctx context.Context, hash string) (*types.ChainBlock, error) {
	return c.chainClient.GetBlock(ctx, hash)
}
