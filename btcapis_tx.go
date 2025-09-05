package btcapis

import (
	"context"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/decoders"
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
	fmt.Printf("!!! (unsupport) import address: %s, publickey: %s\n", address, publickey)
	return nil
}

// 创建PSBT预览交易数据(钱包未签名状态)
func (c *Client) CreatePSBT(ctx context.Context, inputParams *types.TxInputParams) (string, error) {
	fmt.Printf("create psbt: %+v\n", inputParams)
	return c.txClient.CreateTxUsePSBTv0(ctx, inputParams)
}

// 上传经过钱包签名的PSBT数据并进行广播;
func (c *Client) FinalizePSBTAndBroadcast(ctx context.Context, psbt string) (string, error) {
	fmt.Printf("FinalizePSBTAndBroadcast: %s\n", psbt)
	rawTx, err := c.txClient.FinalizePSBT(ctx, psbt)
	if err != nil {
		return "", err
	}

	// 广播交易
	return c.BroadcastRawTx(ctx, rawTx)
}

// 广播签名
func (c *Client) BroadcastRawTx(ctx context.Context, rawtx []byte) (string, error) {
	return c.txClient.BroadcastRawTx(ctx, rawtx)
}
