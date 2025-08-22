// 查询 bitcoind core 节点状态
package bitcoindrpc

import (
	"context"
	"fmt"
)

func (c *Client) GetNetworkInfo(ctx context.Context) (NetworkInfoDTO, error) {
	var res NetworkInfoDTO
	if err := c.rpcCall(ctx, "getnetworkinfo", []any{}, &res); err != nil {
		return NetworkInfoDTO{}, fmt.Errorf("getnetworkinfo: %w", err)
	}

	return res, nil
}
