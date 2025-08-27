package btcapis

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// // 构建交易 (调试已经完成)
// func (c *Client) BuildTx(ctx context.Context) ([]byte, error) {

// 	inputs := []bitcoindrpc.TxInputCreateRawDTO{
// 		bitcoindrpc.NewTxInput("0d31e59675c85f17d942f4510bb4760d9ed4b661df22af3b7cd5ef3c2116626b", 0),
// 	}
// 	outputs := []bitcoindrpc.TxOutputCreateRawDTO{
// 		bitcoindrpc.NewPayToAddress("tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw", 0.001),
// 		bitcoindrpc.NewOpReturn("0100000000000000000000000000000000000000000000000000000000000000"),
// 	}

// 	// locktime := int64(0)
// 	// replaceable := false

// 	dto := bitcoindrpc.TxCreateRawDTO{
// 		Inputs:  inputs,
// 		Outputs: outputs,
// 		// Locktime:    &locktime,
// 		// Replaceable: &replaceable,
// 	}

// 	raw, err := dto.MarshalJSON()
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Printf("dto: %+v\n", string(raw))
// 	fmt.Println("--------------------------------")

// 	return c.txClient.BuildTx(ctx, dto)
// }

// // 填充交易费用 (调试已经完成)
// func (c *Client) FundTx(ctx context.Context, rawtx string) (bitcoindrpc.TxFundRawResultDTO, error) {

// 	options := bitcoindrpc.TxFundOptionsDTO{
// 		AddInputs:     true,
// 		FeeRateSats:   10,
// 		Replaceable:   true,
// 		ChangeAddress: "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw",
// 		// ChangeType: "bech32",
// 	}
// 	fmt.Printf("options: %+v\n", options)
// 	fmt.Println("--------------------------------")

// 	return c.txClient.FundTx(ctx, rawtx, options)
// }

// // 签名交易
// func (c *Client) SignTxWithKey(ctx context.Context, rawtx string) (string, error) {
// 	return c.txClient.TxSignRawWithKey(ctx, rawtx)
// }

// 查询交易信息 => 适用于通过txid查询详细交易信息
func (c *Client) GetTx(ctx context.Context, txid string) (*types.Tx, error) {
	rawtx, err := c.txClient.GetRawTx(ctx, txid)
	if err != nil {
		return nil, err
	}
	ret, err := decoders.DecodeRawTx(rawtx)
	if err != nil {
		return nil, err
	}

	return ret, err
}

// 解析一笔交易元数据 => 适用于外部直接输入交易元数据解析结构
func (c *Client) DecodeRawTx(ctx context.Context, rawtx []byte) (*types.Tx, error) {
	ret, err := decoders.DecodeRawTx(rawtx)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// // 普通转账交易
// func (c *Client) SendBTC(ctx context.Context, inputParams *TxInputParams) {
// 	c.txClient.SendBTC(ctx, *inputParams)
// }
