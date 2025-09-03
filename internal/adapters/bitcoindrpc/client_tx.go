// 交易相关接口
package bitcoindrpc

import (
	"context"
	"encoding/hex"
)

// 获取交易元数据
// [可以修改decodeFlag获取json格式数据, 也可以使用decoderawtransaction(hex)解析raw数据]
// 目前使用btcd库统一解析交易数据的hex.
// decodeFlag: false-返回hex字符串; true-返回json;
func (c *Client) TxGetRaw(ctx context.Context, txid string, decodeFlag bool) ([]byte, error) {
	var hexStr string
	if err := c.rpcCall(ctx, "getrawtransaction", []any{txid, decodeFlag}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

// 构建交易(taproot需要使用psbt)
func (c *Client) TxCreateRaw(ctx context.Context, dto *TxCreateRawDTO) ([]byte, error) {
	var rawtx string
	if err := c.rpcCallWithAny(ctx, "createrawtransaction", dto, &rawtx); err != nil {
		return nil, err
	}
	return hex.DecodeString(rawtx)
}

// 填充交易费用(taproot需要使用psbt)
func (c *Client) TxFundRaw(ctx context.Context, rawtx string, options *TxFundOptionsDTO) (*TxFundRawResultDTO, error) {
	var result *TxFundRawResultDTO
	if err := c.rpcCall(ctx, "fundrawtransaction", []any{rawtx, options}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// 签名交易(taproot需要使用psbt)
func (c *Client) TxSignRawWithKey(ctx context.Context, rawtx string) (string, error) {
	var signedTx string
	if err := c.rpcCall(ctx, "signrawtransactionwithkey", []any{rawtx}, &signedTx); err != nil {
		return "", err
	}
	return signedTx, nil
}

// 完成psbt交易
func (c *Client) TxFinalizePsbt(ctx context.Context, psbt string) (*SignedTxDTO, error) {
	var signedTx *SignedTxDTO
	if err := c.rpcCall(ctx, "finalizepsbt", []any{psbt}, &signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}

// 广播交易
func (c *Client) TxBroadcast(ctx context.Context, rawtx []byte) (string, error) {
	hexRaw := hex.EncodeToString(rawtx)
	var txid string
	if err := c.rpcCall(ctx, "sendrawtransaction", []any{hexRaw}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// 预检查交易 testmempoolaccept: 需要组装交易数据后生成hex字符串再测试
func (c *Client) TxTestMempoolAccept(ctx context.Context, rawtx []byte) (string, error) {
	hexRaw := hex.EncodeToString(rawtx)
	var txid string
	if err := c.rpcCall(ctx, "testmempoolaccept", []any{hexRaw}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}
