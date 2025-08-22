// 地址相关接口
package bitcoindrpc

import (
	"context"
	"fmt"
)

// // 查询钱包余额
// func GetAddressBalance(ctx context.Context, addr string) (int64, int64, error) {
// }

// 查询钱包UTXO集（基于描述符/地址扫描全链 UTXO） - 会导致节点进行全量查询, 慎用
func (c *Client) AddressGetUTXOs(ctx context.Context, addr string) ([]UTXODTO, error) {
	// scantxoutset "start" [ scanobjects ] ; 直接用 addr() 描述符
	params := []interface{}{"start", []interface{}{fmt.Sprintf("addr(%s)", addr)}}
	var res scanResult
	if err := c.rpcCall(ctx, "scantxoutset", params, &res); err != nil {
		return nil, err
	}
	if !res.Success {
		return nil, fmt.Errorf("scantxoutset failed")
	}
	return res.Unspents, nil
}

// 查询钱包详细信息: 根据是否导入到本地节点, 返回数据不同
func (c *Client) AddressGetInfo(ctx context.Context, addr string) (*AddressInfoDTO, error) {
	var res AddressInfoDTO
	if err := c.rpcCall(ctx, "getaddressinfo", []any{addr}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// 校验钱包
func (c *Client) AddressValidate(ctx context.Context, addr string) (*ValidateAddressDTO, error) {
	var res ValidateAddressDTO
	if err := c.rpcCall(ctx, "validateaddress", []any{addr}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
