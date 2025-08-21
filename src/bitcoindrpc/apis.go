package bitcoindrpc

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/crazycloudcc/btcapis/src/decoders"
	"github.com/crazycloudcc/btcapis/src/types"
)

// —— 交易相关
func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	var hexStr string
	if err := c.rpcCall(ctx, "getrawtransaction", []any{txid, false}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

func (c *Client) GetTx(ctx context.Context, txid string) (*types.Tx, error) {
	raw, err := c.GetRawTx(ctx, txid)
	if err != nil {
		return nil, err
	}
	return decoders.DecodeRawTx(raw)
}

// 对外门面：广播交易
func (c *Client) Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	hexRaw := hex.EncodeToString(rawtx)
	var txid string
	if err := c.rpcCall(ctx, "sendrawtransaction", []any{hexRaw}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// 对外门面：估算交易费率
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error) {
	// estimatesmartfee 返回 BTC/kB（可能为 null）
	var resp struct {
		Feerate *float64 `json:"feerate"` // BTC/KB
		Errors  []string `json:"errors"`
	}
	if err := c.rpcCall(ctx, "estimatesmartfee", []any{targetBlocks}, &resp); err != nil {
		return 0, err
	}
	if resp.Feerate == nil {
		return 0, fmt.Errorf("bitcoind: estimatesmartfee no data")
	}
	// BTC/kB -> sats/vB
	satsPerVB := (*resp.Feerate) * 1e8 / 1000.0
	return satsPerVB, nil
}

// —— 其它 ChainReader 方法（先占位，避免接口不满足导致编译失败）
func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
	var hash string
	if err := c.rpcCall(ctx, "getblockhash", []any{height}, &hash); err != nil {
		return "", err
	}
	return hash, nil
}

func (c *Client) GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
	var hexStr string
	if err := c.rpcCall(ctx, "getblockheader", []any{hash, false}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

func (c *Client) GetBlock(ctx context.Context, hash string) ([]byte, error) {
	var hexStr string
	if err := c.rpcCall(ctx, "getblock", []any{hash, 0}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

// 对外门面：查询 UTXO
func (c *Client) GetUTXO(ctx context.Context, op types.OutPoint) (*types.UTXO, error) {
	// 用 gettxout 查询未花费输出
	var dto struct {
		Value        float64 `json:"value"` // BTC
		ScriptPubKey struct {
			Hex string `json:"hex"`
		} `json:"scriptPubKey"`
	}
	if err := c.rpcCall(ctx, "gettxout", []any{op.Hash, op.Index, true}, &dto); err != nil {
		return nil, err
	}
	if dto.ScriptPubKey.Hex == "" {
		return nil, fmt.Errorf("bitcoind: utxo not found")
	}
	spk, _ := hex.DecodeString(dto.ScriptPubKey.Hex)
	return &types.UTXO{
		OutPoint: op,
		Value:    int64(dto.Value * 1e8),
		PkScript: spk,
	}, nil
}

func (c *Client) GetRawMempool(ctx context.Context) ([]string, error) {
	var ids []string
	if err := c.rpcCall(ctx, "getrawmempool", []any{}, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// TxInMempool 实现：简单地遍历 getrawmempool 结果
func (c *Client) TxInMempool(ctx context.Context, txid string) (bool, error) {
	ids, err := c.GetRawMempool(ctx)
	if err != nil {
		return false, err
	}
	for _, id := range ids {
		if id == txid {
			return true, nil
		}
	}
	return false, nil
}
