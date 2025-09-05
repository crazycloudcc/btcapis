// 链相关接口
package bitcoindrpc

import (
	"context"
	"encoding/hex"
	"fmt"
)

// 估算交易费率
func (c *Client) ChainEstimateSmartFeeRate(ctx context.Context, targetBlocks int) (*FeeRateSmartDTO, error) {
	var resp *FeeRateSmartDTO
	if err := c.rpcCall(ctx, "estimatesmartfee", []any{targetBlocks}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// 查询 UTXO
func (c *Client) ChainGetUTXO(ctx context.Context, hash [32]byte, index uint32) ([]byte, int64, error) {
	// 用 gettxout 查询未花费输出
	var dto struct {
		Value        float64 `json:"value"` // BTC
		ScriptPubKey struct {
			Hex string `json:"hex"`
		} `json:"scriptPubKey"`
	}
	if err := c.rpcCall(ctx, "gettxout", []any{hash, index, true}, &dto); err != nil {
		return nil, 0, err
	}

	// 如果 scriptPubKey 为空，则返回错误
	if dto.ScriptPubKey.Hex == "" {
		return nil, 0, fmt.Errorf("bitcoind: utxo not found")
	}
	spk, _ := hex.DecodeString(dto.ScriptPubKey.Hex)
	value := int64(dto.Value * 1e8)
	return spk, value, nil
}

// 获取节点区块数量
func (c *Client) ChainGetBlockCount(ctx context.Context) (int, error) {
	var res int
	if err := c.rpcCall(ctx, "getblockcount", []any{}, &res); err != nil {
		return 0, fmt.Errorf("getblockcount: %w", err)
	}

	return res, nil
}

// 获取最新区块的hash
func (c *Client) ChainGetBestBlockHash(ctx context.Context) (string, error) {
	var res string
	if err := c.rpcCall(ctx, "getbestblockhash", []any{}, &res); err != nil {
		return "", fmt.Errorf("getbestblockhash: %w", err)
	}

	return res, nil
}

// 使用区块高度 查询区块哈希
func (c *Client) ChainGetBlockHash(ctx context.Context, height int64) (string, error) {
	var hash string
	if err := c.rpcCall(ctx, "getblockhash", []any{height}, &hash); err != nil {
		return "", err
	}
	return hash, nil
}

// 使用区块block hash 查询区块头
func (c *Client) ChainGetBlockHeader(ctx context.Context, hash string) (*BlockDTO, error) {
	var dto *BlockDTO
	var verbose bool = true // false-返回hex字符串; true-返回json;
	if err := c.rpcCall(ctx, "getblockheader", []any{hash, verbose}, &dto); err != nil {
		return nil, err
	}
	return dto, nil
}

// 使用区块block hash 查询区块
func (c *Client) ChainGetBlock(ctx context.Context, hash string) (*BlockDTO, error) {
	var dto *BlockDTO
	var verbosity int = 1 // 0-返回hex字符串; 1-返回json(其中tx数组返回txid字符串); 2-返回json(其中tx数组返回完整tx数据); 3-返回json(其中tx数组返回完整tx数据);
	if err := c.rpcCall(ctx, "getblock", []any{hash, verbosity}, &dto); err != nil {
		return nil, err
	}
	return dto, nil
}
