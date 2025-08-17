package btcapis

import (
	"context"
	"fmt"
	"strings"

	"github.com/crazycloudcc/btcapis/chain"
	"github.com/crazycloudcc/btcapis/types"
)

// GetAddressBalance returns the balance string for an address in format "confirmed(mempool)" in BTC.
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (string, error) {
	confirmed, mempool, err := c.addressBalance(ctx, addr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", satsToBTC(confirmed), satsToBTC(mempool)), nil
}

// GetAddressUTXOs returns UTXOs belonging to the address.
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	for _, b := range append(c.primaries, c.fallbacks...) {
		if ar, ok := b.(chain.AddressReader); ok {
			if u, err := ar.AddressUTXOs(ctx, addr); err == nil {
				return u, nil
			}
		}
	}
	return nil, chain.ErrBackendUnavailable
}

func (c *Client) addressBalance(ctx context.Context, addr string) (int64, int64, error) {
	for _, b := range append(c.primaries, c.fallbacks...) {
		if ar, ok := b.(chain.AddressReader); ok {
			if confirmed, mempool, err := ar.AddressBalance(ctx, addr); err == nil {
				return confirmed, mempool, nil
			}
		}
	}
	return 0, 0, chain.ErrBackendUnavailable
}

func satsToBTC(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	f := float64(v) / 1e8
	s := fmt.Sprintf("%.8f", f)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return sign + s
}
