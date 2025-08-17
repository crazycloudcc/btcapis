// Package mempoolspace 数据映射器
package mempoolspace

import (
	"encoding/hex"

	"github.com/crazycloudcc/btcapis/types"
)

func mapTxDTO(d TxDTO) *types.Tx {
	t := &types.Tx{
		TxID:     d.Txid,
		Version:  d.Version,
		LockTime: d.Locktime,
		Weight:   d.Weight,
		Vsize:    (d.Weight + 3) / 4,
		Vin:      make([]types.TxIn, len(d.Vin)),
		Vout:     make([]types.TxOut, len(d.Vout)),
	}
	for i, in := range d.Vin {
		w := make([][]byte, len(in.Witness))
		for j, wh := range in.Witness {
			b, _ := hex.DecodeString(wh)
			w[j] = b
		}
		var ss []byte
		if in.Scriptsig != "" {
			ss, _ = hex.DecodeString(in.Scriptsig)
		}
		t.Vin[i] = types.TxIn{
			TxID:      in.Txid,
			Vout:      in.Vout,
			Sequence:  in.Sequence,
			ScriptSig: ss,
			Witness:   w,
		}
	}
	for i, out := range d.Vout {
		spk, _ := hex.DecodeString(out.ScriptPubKey)
		addrs := []string(nil)
		if out.Address != "" {
			addrs = []string{out.Address}
		}
		t.Vout[i] = types.TxOut{
			Value:        out.Value,
			ScriptPubKey: spk,
			Type:         out.ScriptType,
			Addresses:    addrs,
		}
	}
	return t
}
