package btcapis

import (
	"context"
	"errors"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
)

// // Broadcaster 提供交易广播功能.
// type Broadcaster interface {
// 	// Broadcast 广播交易.
// 	Broadcast(ctx context.Context, rawtx []byte) (txid string, err error)
// }

// Broadcast 广播交易.
func Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	if bitcoindrpc.IsInited() {
		return bitcoindrpc.Broadcast(ctx, rawtx)
	}

	if mempoolapis.IsInited() {
		return mempoolapis.Broadcast(ctx, rawtx)
	}

	return "", errors.New("btcapis: no client available")
}
