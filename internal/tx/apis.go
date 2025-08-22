package tx

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
)

// 创建普通交易

// 构建交易
func (c *Client) BuildTx(ctx context.Context, dto bitcoindrpc.TxCreateRawDTO) ([]byte, error) {
	return c.bitcoindrpcClient.TxCreateRaw(ctx, dto)
}

// 填充交易费用
func (c *Client) FundTx(ctx context.Context, rawtx string, options bitcoindrpc.TxFundOptionsDTO) (bitcoindrpc.TxFundRawResultDTO, error) {
	return c.bitcoindrpcClient.TxFundRaw(ctx, rawtx, options)
}
