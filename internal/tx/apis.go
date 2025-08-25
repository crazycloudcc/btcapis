package tx

import (
	"context"
	"encoding/hex"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// 普通转账交易预览 =>
func (c *Client) SendBTCByPSBTPreview(ctx context.Context, inputParams *TxInputParams) ([]byte, error) {

	dto := bitcoindrpc.TxCreateRawDTO{
		Inputs: []bitcoindrpc.TxInputCreateRawDTO{
			{
				TxID: inputParams.FromAddress[0],
				Vout: 0,
			},
		},
		Outputs: []bitcoindrpc.TxOutputCreateRawDTO{},
	}

	for i, to := range inputParams.ToAddress {
		dto.Outputs = append(dto.Outputs, bitcoindrpc.TxOutputCreateRawDTO{
			Address: to,
			Amount:  inputParams.AmountBTC[i],
		})
	}

	if inputParams.Data != "" {
		dto.Outputs = append(dto.Outputs, bitcoindrpc.TxOutputCreateRawDTO{
			DataHex: inputParams.Data,
		})
	}

	// step.1 构建交易
	raw, err := c.createTx(ctx, dto)
	if err != nil {
		return nil, err
	}

	// step.2 填充交易费用
	funded, err := c.fundTx(ctx, string(raw), bitcoindrpc.TxFundOptionsDTO{})
	if err != nil {
		return nil, err
	}

	// step.3 签名交易
	signed, err := c.signTx(ctx, funded.Hex)
	if err != nil {
		return nil, err
	}

	// step.4 广播交易
	broadcasted, err := c.broadcastTx(ctx, signed)
	if err != nil {
		return nil, err
	}

	return broadcasted, nil
}

// 接收外部钱包签名后的PSBT交易数据并广播
func (c *Client) SendBTCByPSBT(ctx context.Context, psbt string) (string, error) {
	return "", nil
}

// 查询交易元数据: 优先使用bitcoindrpc, 如果没有再使用mempoolapis
func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	ret, err := c.bitcoindrpcClient.TxGetRaw(ctx, txid, false)
	if err == nil {
		return ret, nil
	}

	ret, err = c.mempoolapisClient.TxGetRaw(ctx, txid)
	if err == nil {
		return ret, nil
	}

	return nil, err
}

// 按照types.Tx格式返回交易数据
func (c *Client) GetTx(ctx context.Context, txid string) (*types.Tx, error) {
	raw, err := c.GetRawTx(ctx, txid)
	if err != nil {
		return nil, err
	}

	ret, err := decoders.DecodeRawTx(raw)
	if err != nil {
		return nil, err
	}

	return ret, err
}

// 构建交易 createrawtransaction
func (c *Client) createTx(ctx context.Context, dto bitcoindrpc.TxCreateRawDTO) ([]byte, error) {
	return c.bitcoindrpcClient.TxCreateRaw(ctx, dto)
}

// 填充交易费用
func (c *Client) fundTx(ctx context.Context, rawtx string, options bitcoindrpc.TxFundOptionsDTO) (bitcoindrpc.TxFundRawResultDTO, error) {
	return c.bitcoindrpcClient.TxFundRaw(ctx, rawtx, options)
}

// 签名交易
func (c *Client) signTx(ctx context.Context, rawtx string) (string, error) {
	return c.bitcoindrpcClient.TxSignRawWithKey(ctx, rawtx)
}

// 广播交易
func (c *Client) broadcastTx(ctx context.Context, rawtx string) (string, error) {
	bin, err := hex.DecodeString(rawtx)
	if err != nil {
		return "", err
	}
	return c.bitcoindrpcClient.TxBroadcast(ctx, bin)
}
