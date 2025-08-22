package btcapis

// // GetRawTx 返回交易原始数据.
// func GetRawTx(ctx context.Context, txid string) ([]byte, error) {
// 	if bitcoindrpc.IsInited() {
// 		return bitcoindrpc.GetRawTx(ctx, txid)
// 	}

// 	if mempoolapis.IsInited() {
// 		return mempoolapis.GetRawTx(ctx, txid)
// 	}

// 	return nil, errors.New("btcapis: no client available")
// }

// // GetTx 返回交易.(优先使用bitcoindrpcClient, 其次使用mempoolspaceClient, 两边的数据格式不一致, 所以需要兼容)
// func GetTx(ctx context.Context, txid string) (*types.Tx, error) {
// 	if bitcoindrpc.IsInited() {
// 		raw, err := bitcoindrpc.GetRawTx(ctx, txid)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return decoders.DecodeRawTx(raw)
// 	}

// 	if mempoolapis.IsInited() {
// 		raw, err := mempoolapis.GetRawTx(ctx, txid)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return decoders.DecodeRawTx(raw)
// 	}

// 	return nil, errors.New("btcapis: no client available")
// }
