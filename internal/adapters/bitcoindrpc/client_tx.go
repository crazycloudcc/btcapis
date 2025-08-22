// 交易相关接口
package bitcoindrpc

import (
	"context"
	"encoding/hex"
)

// 获取交易元数据
// [可以修改decodeFlag获取json格式数据, 也可以使用decoderawtransaction(hex)解析raw数据]
// 目前使用btcd库统一解析交易数据的hex.
func (c *Client) TxGetRaw(ctx context.Context, txid string) ([]byte, error) {
	var hexStr string
	var decodeFlag bool = false // false-返回hex字符串; true-返回json;
	if err := c.rpcCall(ctx, "getrawtransaction", []any{txid, decodeFlag}, &hexStr); err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

// 构建交易
func (c *Client) TxCreateRaw(ctx context.Context, dto TxCreateRawDTO) ([]byte, error) {
	var rawtx string
	params, err := dto.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if err := c.rpcCall(ctx, "createrawtransaction", []any{params}, &rawtx); err != nil {
		return nil, err
	}
	return hex.DecodeString(rawtx)
}

// // 获取交易费用
// func (c *Client) TxFundRaw(ctx context.Context, dto TxFundRawDTO) ([]byte, error) {
// 	var rawtx string
// 	if err := c.rpcCall(ctx, "fundrawtransaction", []any{dto}, &rawtx); err != nil {
// 		return nil, err
// 	}
// 	return hex.DecodeString(rawtx)
// }

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
