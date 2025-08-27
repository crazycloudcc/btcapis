package decoders

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/types"
)

func DecodeRawTx(raw []byte) (*types.Tx, error) {
	var m wire.MsgTx
	if err := m.Deserialize(bytes.NewReader(raw)); err != nil {
		return nil, err
	}

	t := &types.Tx{
		Version:  m.Version,
		LockTime: m.LockTime,
		TxIn:     make([]types.TxIn, len(m.TxIn)),
		TxOut:    make([]types.TxOut, len(m.TxOut)),
	}

	for i, in := range m.TxIn {
		w := make([][]byte, len(in.Witness))
		copy(w, in.Witness)
		t.TxIn[i] = types.TxIn{
			PreviousOutPoint: types.TxOutPoint{
				Hash:  types.Hash32(in.PreviousOutPoint.Hash), // ccflag: 小端序
				Index: in.PreviousOutPoint.Index,
			},
			Sequence:  in.Sequence,
			ScriptSig: append([]byte(nil), in.SignatureScript...),
			Witness:   w,
		}
	}

	for i, o := range m.TxOut {
		spk := append([]byte(nil), o.PkScript...)
		addrInfo, err := DecodePkScript(spk)
		if err != nil {
			return nil, err
		}
		t.TxOut[i] = types.TxOut{
			Value:      o.Value,
			PkScript:   spk,
			ScriptType: string(addrInfo.Typ),
			Address:    addrInfo.Addresses[0],
		}
	}

	return t, nil
}
