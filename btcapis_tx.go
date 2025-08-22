package btcapis

import (
	"context"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
)

// 构建交易
func (c *Client) BuildTx(ctx context.Context) ([]byte, error) {

	inputs := []bitcoindrpc.TxInputCreateRawDTO{
		{
			TxID: "0d31e59675c85f17d942f4510bb4760d9ed4b661df22af3b7cd5ef3c2116626b",
			Vout: 0,
		},
	}
	outputs := []bitcoindrpc.TxOutputCreateRawDTO{
		{
			Address: "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw",
			Amount:  10000,
		},
		{
			DataHex: "0100000000000000000000000000000000000000000000000000000000000000",
		},
	}

	// locktime := int64(0)
	// replaceable := false

	dto := bitcoindrpc.TxCreateRawDTO{
		Inputs:  inputs,
		Outputs: outputs,
		// Locktime:    &locktime,
		// Replaceable: &replaceable,
	}

	raw, err := dto.MarshalJSON()
	if err != nil {
		return nil, err
	}

	fmt.Printf("dto: %+v\n", raw)
	fmt.Println("--------------------------------")

	return c.txClient.BuildTx(ctx, dto)
}

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
