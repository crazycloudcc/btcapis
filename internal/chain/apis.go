package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/utils"
)

func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, float64, error) {
	var fee1, fee2 float64 = 0, 0
	if c.bitcoindrpcClient != nil {
		dto, err := c.bitcoindrpcClient.ChainEstimateSmartFeeRate(ctx, targetBlocks)
		if err != nil {
			fee1 = 0
		} else {
			fee1 = utils.BTCToSats(dto.Feerate)
		}
	}

	if c.mempoolapisClient != nil {
		dto, err := c.mempoolapisClient.EstimateFeeRate(ctx, targetBlocks)
		if err != nil {
			fee2 = 0
		} else {
			fee2 = dto.FastestFee
		}
	}

	return fee1, fee2, nil
}

// // 查询 UTXO
// func (c *Client) GetUTXO(ctx context.Context, hash [32]byte, index uint32) ([]byte, int64, error) {
// 	return c.bitcoindrpcClient.ChainGetUTXO(ctx, hash, index)
// }

// // 获取节点区块数量
// func (c *Client) GetBlockCount(ctx context.Context) (int, error) {
// 	return c.bitcoindrpcClient.ChainGetBlockCount(ctx)
// }

// // 获取最新区块的hash
// func (c *Client) GetBestBlockHash(ctx context.Context) (string, error) {
// 	return c.bitcoindrpcClient.ChainGetBestBlockHash(ctx)
// }

// // 使用区块高度 查询区块哈希
// func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
// 	return c.bitcoindrpcClient.ChainGetBlockHash(ctx, height)
// }

// // 使用区块block hash 查询区块头
// func (c *Client) GetBlockHeader(ctx context.Context, hash string) (*types.ChainBlock, error) {
// 	blk, err := c.bitcoindrpcClient.ChainGetBlockHeader(ctx, hash)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &types.ChainBlock{
// 		Hash:              blk.Hash,
// 		Confirmations:     blk.Confirmations,
// 		Height:            blk.Height,
// 		Version:           blk.Version,
// 		VersionHex:        blk.VersionHex,
// 		MerkleRoot:        blk.MerkleRoot,
// 		Time:              blk.Time,
// 		MedianTime:        blk.MedianTime,
// 		Nonce:             blk.Nonce,
// 		Bits:              blk.Bits,
// 		Difficulty:        blk.Difficulty,
// 		Chainwork:         blk.Chainwork,
// 		NTx:               blk.NTx,
// 		PreviousBlockHash: blk.PreviousBlockHash,
// 		NextBlockHash:     blk.NextBlockHash,
// 	}, nil
// }

// // 使用区块block hash 查询区块
// func (c *Client) GetBlock(ctx context.Context, hash string) (*types.ChainBlock, error) {
// 	blk, err := c.bitcoindrpcClient.ChainGetBlock(ctx, hash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &types.ChainBlock{
// 		Hash:              blk.Hash,
// 		Confirmations:     blk.Confirmations,
// 		Height:            blk.Height,
// 		Version:           blk.Version,
// 		VersionHex:        blk.VersionHex,
// 		MerkleRoot:        blk.MerkleRoot,
// 		Time:              blk.Time,
// 		MedianTime:        blk.MedianTime,
// 		Nonce:             blk.Nonce,
// 		Bits:              blk.Bits,
// 		Difficulty:        blk.Difficulty,
// 		Chainwork:         blk.Chainwork,
// 		NTx:               blk.NTx,
// 		PreviousBlockHash: blk.PreviousBlockHash,
// 		NextBlockHash:     blk.NextBlockHash,
// 		StrippedSize:      blk.StrippedSize,
// 		Size:              blk.Size,
// 		Weight:            blk.Weight,
// 		Tx:                blk.Tx,
// 	}, nil
// }
