package mempoolspace

import (
	"context"
	"encoding/hex"
	"path"

	"github.com/crazycloudcc/btcapis/types"
)

// AddressBalance implements chain.AddressReader#AddressBalance.
func (c *Client) AddressBalance(ctx context.Context, addr string) (int64, int64, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/address/", addr)
	var dto struct {
		ChainStats struct {
			Funded int64 `json:"funded_txo_sum"`
			Spent  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
		MempoolStats struct {
			Funded int64 `json:"funded_txo_sum"`
			Spent  int64 `json:"spent_txo_sum"`
		} `json:"mempool_stats"`
	}
	if err := c.getJSON(ctx, u.String(), &dto); err != nil {
		return 0, 0, err
	}
	confirmed := dto.ChainStats.Funded - dto.ChainStats.Spent
	mempool := dto.MempoolStats.Funded - dto.MempoolStats.Spent
	return confirmed, mempool, nil
}

// AddressUTXOs implements chain.AddressReader#AddressUTXOs.
func (c *Client) AddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/address/", addr, "/utxo")
	var dtos []struct {
		Txid   string `json:"txid"`
		Vout   uint32 `json:"vout"`
		Value  int64  `json:"value"`
		Status struct {
			Confirmed   bool  `json:"confirmed"`
			BlockHeight int64 `json:"block_height"`
		} `json:"status"`
	}
	if err := c.getJSON(ctx, u.String(), &dtos); err != nil {
		return nil, err
	}
	utxos := make([]types.UTXO, 0, len(dtos))
	for _, d := range dtos {
		txidBytes, _ := hex.DecodeString(d.Txid)
		u := types.UTXO{
			OutPoint: types.OutPoint{Hash: types.Hash32(txidBytes), Index: d.Vout},
			Value:    d.Value,
		}
		if d.Status.Confirmed {
			u.Height = uint32(d.Status.BlockHeight)
		}
		utxos = append(utxos, u)
	}
	return utxos, nil
}
