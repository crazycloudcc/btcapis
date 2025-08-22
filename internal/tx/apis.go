package tx

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
)

// 构建交易
func (c *Client) BuildTx(ctx context.Context, dto bitcoindrpc.TxCreateRawDTO) ([]byte, error) {
	return c.bitcoindrpcClient.TxCreateRaw(ctx, dto)
}
