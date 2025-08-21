package btcapis

import "context"

// Broadcaster 提供交易广播功能.
type Broadcaster interface {
	// Broadcast 广播交易.
	Broadcast(ctx context.Context, rawtx []byte) (txid string, err error)
}

func (c *Client) Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	return c.mempoolspaceClient.Broadcast(ctx, rawtx)
}
