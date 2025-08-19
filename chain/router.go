// 多后端聚合/故障转移/负载策略（供根包 Client 使用）
package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/script"
	"github.com/crazycloudcc/btcapis/types"
)

type Router struct {
	primaries []Backend // 建议 bitcoind 放这里
	fallbacks []Backend // mempool.space / electrum 等
}

func NewRouter(primaries, fallbacks []Backend) *Router {
	return &Router{primaries: primaries, fallbacks: fallbacks}
}

func (r *Router) GetTransaction(ctx context.Context, txid string) (*types.Tx, error) {
	// 1) 首选 raw + 本地解析
	for _, b := range r.primaries {
		if raw, err := b.GetRawTransaction(ctx, txid); err == nil && len(raw) > 0 {
			if t, err := script.DecodeRawTx(raw); err == nil {
				return t, nil
			}
		}
	}
	// 2) 降级：直接返回解析后的 Tx
	for _, b := range r.fallbacks {
		if t, err := b.GetTx(ctx, txid); err == nil && t != nil {
			return t, nil
		}
	}
	// 3) 再降级：fallback raw + 本地解析
	for _, b := range r.fallbacks {
		if raw, err := b.GetRawTransaction(ctx, txid); err == nil && len(raw) > 0 {
			if t, err := script.DecodeRawTx(raw); err == nil {
				return t, nil
			}
		}
	}
	return nil, ErrTxNotFound
}

func (r *Router) GetRawTransaction(ctx context.Context, txid string) ([]byte, error) {
	for _, b := range r.primaries {
		if raw, err := b.GetRawTransaction(ctx, txid); err == nil && len(raw) > 0 {
			return raw, nil
		}
	}
	for _, b := range r.fallbacks {
		if raw, err := b.GetRawTransaction(ctx, txid); err == nil && len(raw) > 0 {
			return raw, nil
		}
	}
	return nil, ErrTxNotFound
}

func (r *Router) EstimateFeeRate(ctx context.Context, target int) (float64, error) {
	for _, b := range r.primaries {
		if v, err := b.EstimateFeeRate(ctx, target); err == nil {
			return v, nil
		}
	}
	for _, b := range r.fallbacks {
		if v, err := b.EstimateFeeRate(ctx, target); err == nil {
			return v, nil
		}
	}
	return 0, ErrBackendUnavailable
}

func (r *Router) Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	for _, b := range r.primaries {
		if txid, err := b.Broadcast(ctx, rawtx); err == nil && txid != "" {
			return txid, nil
		}
	}
	for _, b := range r.fallbacks {
		if txid, err := b.Broadcast(ctx, rawtx); err == nil && txid != "" {
			return txid, nil
		}
	}
	return "", ErrBackendUnavailable
}
