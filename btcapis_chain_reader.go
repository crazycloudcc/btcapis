package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
)

// // ChainReader 提供区块相关的查询, 如区块哈希、区块头、区块等.
// type ChainReader interface {
// 	// GetBlockHash 返回区块哈希.
// 	GetBlockHash(ctx context.Context, height int64) (string, error)
// 	// GetBlockHeader 返回区块头.
// 	GetBlockHeader(ctx context.Context, hash string) ([]byte, error)
// 	// GetBlock 返回区块.
// 	GetBlock(ctx context.Context, hash string) ([]byte, error)
// 	// GetUTXO 返回UTXO详细信息.(因为这个UTXO不是通过关联钱包查询, 所以暂时放在ChainReader分类)
// 	GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error)
// }

// GetBlockHash 返回区块哈希.
func GetBlockHash(ctx context.Context, height int64) (string, error) {
	return bitcoindrpc.GetBlockHash(ctx, height)
}

// GetBlockHeader 返回区块头.
func GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
	return bitcoindrpc.GetBlockHeader(ctx, hash)
}

// GetBlock 返回区块.
func GetBlock(ctx context.Context, hash string) ([]byte, error) {
	return bitcoindrpc.GetBlock(ctx, hash)
}

// GetUTXO 返回UTXO详细信息.(因为这个UTXO不是通过关联钱包查询, 所以暂时放在ChainReader分类)
// func GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error) {
// 	return bitcoindrpc.GetUTXO(ctx, outpoint)
// }
