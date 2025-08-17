package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// AddressReader provides address related queries such as balance and utxos.
type AddressReader interface {
	// AddressBalance returns confirmed balance and mempool (unconfirmed) delta in sats.
	AddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error)
	// AddressUTXOs returns UTXOs owned by the address.
	AddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error)
}
