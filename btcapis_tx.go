package btcapis

import (
	"context"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/chain"
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

// GetTxVerbose 查询 verbose 格式交易（bitcoin core → mempool.space 降级）
func (c *Client) GetTxVerbose(ctx context.Context, txid string, verbosity int) (*types.TxVerbose, error) {
	if c == nil || c.txClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.txClient.GetTxVerbose(ctx, txid, verbosity)
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

// 校验psbt base64串是否合法
func (c *Client) ValidateUnsignedPsbtBase64(ctx context.Context, psbtBase64 string) error {
	return c.txClient.ValidateUnsignedPsbtBase64(ctx, psbtBase64)
}

// 校验已签名psbt的base64串
func (c *Client) ValidateSignedPsbtBase64(ctx context.Context, psbtBase64 string) (string, error) {
	return c.txClient.ValidateSignedPsbtBase64(ctx, psbtBase64)
}

// transfer all to new address
func (c *Client) TransferAllToNewAddress(ctx context.Context, toAddress string, privateKeyWIF string, fromAddress string, feeRate float64) (string, error) {
	return c.txClient.TransferAllToNewAddress(ctx, toAddress, privateKeyWIF, fromAddress, feeRate)
}

// TestMempoolAccept 预检查交易是否可被内存池接受
func (c *Client) TestMempoolAccept(ctx context.Context, rawtx []byte) (string, error) {
	if c == nil || c.chainClient == nil {
		return "", chain.ErrBitcoindUnavailable
	}
	return c.chainClient.TestMempoolAccept(ctx, rawtx)
}

// GetMempoolTxIds 拉取内存池交易 txid 列表
func (c *Client) GetMempoolTxIds(ctx context.Context) ([]string, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetMempoolTxIds(ctx)
}

// GetMempoolTxEntry 获取内存池交易详情
func (c *Client) GetMempoolTxEntry(ctx context.Context, txid string) (*types.MempoolTxEntry, error) {
	if c == nil || c.chainClient == nil {
		return nil, chain.ErrBitcoindUnavailable
	}
	return c.chainClient.GetMempoolTxEntry(ctx, txid)
}
