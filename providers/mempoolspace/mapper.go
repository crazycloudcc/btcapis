// Package mempoolspace 数据映射器
package mempoolspace

import (
	"encoding/hex"

	"github.com/crazycloudcc/btcapis/types"
)

func mapTxDTO(d TxDTO) *types.Tx {
	t := &types.Tx{
		Version:  d.Version,
		LockTime: d.Locktime,
		TxIn:     make([]types.TxIn, len(d.Vin)),
		TxOut:    make([]types.TxOut, len(d.Vout)),
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

		txidBytes, _ := hex.DecodeString(in.Txid)
		t.TxIn[i] = types.TxIn{
			PreviousOutPoint: types.OutPoint{Hash: types.Hash32(txidBytes), Index: in.Vout},
			Sequence:         in.Sequence,
			ScriptSig:        ss,
			Witness:          w,
		}
	}
	for i, out := range d.Vout {
		spk, _ := hex.DecodeString(out.ScriptPubKey)
		addrs := []string(nil)
		if out.Address != "" {
			addrs = []string{out.Address}
		}
		t.TxOut[i] = types.TxOut{
			Value:      out.Value,
			PkScript:   spk,
			ScriptType: out.ScriptType,
			Address:    addrs[0],
		}
	}
	return t
}
