package btcapis

import "context"

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

// // Broadcast 广播交易.
// func Broadcast(ctx context.Context, rawtx []byte) (string, error) {
// 	if bitcoindrpc.IsInited() {
// 		return bitcoindrpc.Broadcast(ctx, rawtx)
// 	}

// 	if mempoolapis.IsInited() {
// 		return mempoolapis.Broadcast(ctx, rawtx)
// 	}

// 	return "", errors.New("btcapis: no client available")
// }

// EstimateFeeRate 估计手续费.
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (float64, float64, error) {
	return c.chainClient.EstimateFeeRate(ctx, targetBlocks)
}
