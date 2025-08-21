package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/types"
)

// ChainReader 提供区块相关的查询, 如区块哈希、区块头、区块等.
type ChainReader interface {
	// GetBlockHash 返回区块哈希.
	GetBlockHash(ctx context.Context, height int64) (string, error)
	// GetBlockHeader 返回区块头.
	GetBlockHeader(ctx context.Context, hash string) ([]byte, error)
	// GetBlock 返回区块.
	GetBlock(ctx context.Context, hash string) ([]byte, error)
	// GetUTXO 返回UTXO详细信息.(因为这个UTXO不是通过关联钱包查询, 所以暂时放在ChainReader分类)
	GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error)
}

func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	return c.bitcoindrpcClient.GetBlockHash(ctx, height)
}

func (c *Client) GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
	return c.bitcoindrpcClient.GetBlockHeader(ctx, hash)
}

func (c *Client) GetBlock(ctx context.Context, hash string) ([]byte, error) {
	return c.bitcoindrpcClient.GetBlock(ctx, hash)
}

func (c *Client) GetUTXO(ctx context.Context, outpoint types.OutPoint) (*types.UTXO, error) {
	return c.bitcoindrpcClient.GetUTXO(ctx, outpoint)
}
