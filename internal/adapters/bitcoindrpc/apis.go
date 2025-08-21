package bitcoindrpc

import (
	"context"
	"encoding/hex"
	"fmt"
)

// 获取交易
func GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	var hexStr string
	if err := rpcCall(ctx, "getrawtransaction", []any{txid, false}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

// 广播交易
func Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	hexRaw := hex.EncodeToString(rawtx)
	var txid string
	if err := rpcCall(ctx, "sendrawtransaction", []any{hexRaw}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// // 查询钱包余额
// func GetAddressBalance(ctx context.Context, addr string) (int64, int64, error) {
// }

// // 查询钱包UTXO集
// func GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
// }

// 查询 UTXO
func GetUTXO(ctx context.Context, hash [32]byte, index uint32) ([]byte, int64, error) {
	// 用 gettxout 查询未花费输出
	var dto struct {
		Value        float64 `json:"value"` // BTC
		ScriptPubKey struct {
			Hex string `json:"hex"`
		} `json:"scriptPubKey"`
	}
	if err := rpcCall(ctx, "gettxout", []any{hash, index, true}, &dto); err != nil {
		return nil, 0, err
	}

	// 如果 scriptPubKey 为空，则返回错误
	if dto.ScriptPubKey.Hex == "" {
		return nil, 0, fmt.Errorf("bitcoind: utxo not found")
	}
	spk, _ := hex.DecodeString(dto.ScriptPubKey.Hex)
	value := int64(dto.Value * 1e8)
	return spk, value, nil

	// 外部组装, 保持最小化封装
	// return &types.UTXO{
	// 	OutPoint: op,
	// 	Value:    int64(dto.Value * 1e8),
	// 	PkScript: spk,
	// }, nil
}

// 估算交易费率
func EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, error) {
	// estimatesmartfee 返回 BTC/kB（可能为 null）
	var resp struct {
		Feerate *float64 `json:"feerate"` // BTC/KB
		Errors  []string `json:"errors"`
	}
	if err := rpcCall(ctx, "estimatesmartfee", []any{targetBlocks}, &resp); err != nil {
		return 0, err
	}
	if resp.Feerate == nil {
		return 0, fmt.Errorf("bitcoind: estimatesmartfee no data")
	}
	// BTC/kB -> sats/vB
	satsPerVB := (*resp.Feerate) * 1e8 / 1000.0
	return satsPerVB, nil
}

// 查询区块哈希
func GetBlockHash(ctx context.Context, height int64) (string, error) {
	var hash string
	if err := rpcCall(ctx, "getblockhash", []any{height}, &hash); err != nil {
		return "", err
	}
	return hash, nil
}

// 查询区块头
func GetBlockHeader(ctx context.Context, hash string) ([]byte, error) {
	var hexStr string
	if err := rpcCall(ctx, "getblockheader", []any{hash, false}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

// 查询区块
func GetBlock(ctx context.Context, hash string) ([]byte, error) {
	var hexStr string
	if err := rpcCall(ctx, "getblock", []any{hash, 0}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

// // 查询 mempool
// func GetRawMempool(ctx context.Context) ([]string, error) {
// 	var ids []string
// 	if err := rpcCall(ctx, "getrawmempool", []any{}, &ids); err != nil {
// 		return nil, err
// 	}
// 	return ids, nil
// }

// // TxInMempool 实现：简单地遍历 getrawmempool 结果
// func TxInMempool(ctx context.Context, txid string) (bool, error) {
// 	ids, err := GetRawMempool(ctx)
// 	if err != nil {
// 		return false, err
// 	}
// 	for _, id := range ids {
// 		if id == txid {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }
