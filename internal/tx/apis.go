package tx

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// 转账交易-使用PSBTv0版本
func (c *Client) CreateTxUsePSBTv0(ctx context.Context, inputParams *types.TxInputParams) (string, error) {
	tx, utxos, err := c.createNormalTx(ctx, inputParams)
	if err != nil {
		return "", err
	}

	unsignedPsbt, err := c.MsgTxToPSBTV0(ctx, tx, inputParams, utxos)
	if err != nil {
		return "", err
	}

	return unsignedPsbt.PSBTBase64, nil
}

// 广播交易
func (c *Client) BroadcastRawTx(ctx context.Context, rawTx []byte) (string, error) {
	return c.bitcoindrpcClient.TxBroadcast(ctx, rawTx)
}

// 查询交易元数据: 优先使用bitcoindrpc, 如果没有再使用mempoolapis
func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	ret, err := c.bitcoindrpcClient.TxGetRaw(ctx, txid, false)
	if err == nil {
		return ret, nil
	}

	if c.mempoolapisClient != nil {
		ret, err = c.mempoolapisClient.TxGetRaw(ctx, txid)
		if err == nil {
			return ret, nil
		}
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
