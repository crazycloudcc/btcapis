// tx/decoder.go
package tx

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/script"
	"github.com/crazycloudcc/btcapis/types"
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
			PreviousOutPoint: types.OutPoint{
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
		typ, addrs := script.Classify(spk)
		t.TxOut[i] = types.TxOut{
			Value:      o.Value,
			PkScript:   spk,
			ScriptType: typ,
			Address:    addrs[0],
		}
	}

	return t, nil
}

func computeWeightVSize(m *wire.MsgTx) (int64, int64) {
	var wbuf bytes.Buffer
	_ = m.Serialize(&wbuf) // 含 witness
	total := wbuf.Len()

	var sbuf bytes.Buffer
	_ = m.SerializeNoWitness(&sbuf) // 去除 witness
	stripped := sbuf.Len()

	witness := total - stripped
	weight := int64(stripped*4 + witness)
	vsize := (weight + 3) / 4 // 向上取整
	return weight, vsize
}
