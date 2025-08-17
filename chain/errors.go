// Package chain 错误定义
package chain

import "errors"

var (
	ErrTxNotFound         = errors.New("btcapis/chain: transaction not found")
	ErrBackendUnavailable = errors.New("btcapis/chain: backend unavailable")
)
