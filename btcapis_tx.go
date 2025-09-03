package btcapis

import (
	"context"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/psbt"
	"github.com/crazycloudcc/btcapis/types"
)

// 查询交易元数据
func (c *Client) GetTxRaw(ctx context.Context, txid string) ([]byte, error) {
	return c.txClient.GetRawTx(ctx, txid)
}

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

// 上传钱包+publickey, 用于后续组装PSBT等数据, 后续需要在postgres创建映射;
func (c *Client) ImportAddressAndPublickey(ctx context.Context, address string, publickey string) error {
	fmt.Printf("import address: %s, publickey: %s\n", address, publickey)
	return nil
}

// 创建PSBT预览交易数据(钱包未签名状态)
func (c *Client) CreatePSBT(ctx context.Context, inputParams *types.TxInputParams) (*psbt.BuildResult, error) {
	fmt.Printf("create psbt: %+v\n", inputParams)
	return c.txClient.SendBTCByPSBTPreview(ctx, inputParams)
}

// 上传经过钱包签名的PSBT数据, 用于后续广播交易;
func (c *Client) SendBTCByPSBT(ctx context.Context, psbt string) (string, error) {
	fmt.Printf("send psbt: %s\n", psbt)
	return c.txClient.SendBTCByPSBT(ctx, psbt)
}

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
