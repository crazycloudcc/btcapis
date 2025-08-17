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
	weight, vsize := computeWeightVSize(&m)

	t := &types.Tx{
		TxID:     m.TxHash().String(),
		Version:  m.Version,
		LockTime: m.LockTime,
		Weight:   weight,
		Vsize:    vsize,
		Vin:      make([]types.TxIn, len(m.TxIn)),
		Vout:     make([]types.TxOut, len(m.TxOut)),
	}

	for i, in := range m.TxIn {
		w := make([][]byte, len(in.Witness))
		copy(w, in.Witness)
		t.Vin[i] = types.TxIn{
			TxID:      in.PreviousOutPoint.Hash.String(),
			Vout:      in.PreviousOutPoint.Index,
			Sequence:  in.Sequence,
			ScriptSig: append([]byte(nil), in.SignatureScript...),
			Witness:   w,
		}
	}

	for i, o := range m.TxOut {
		spk := append([]byte(nil), o.PkScript...)
		typ, addrs := script.Classify(spk)
		t.Vout[i] = types.TxOut{
			Value:        o.Value,
			ScriptPubKey: spk,
			Type:         typ,
			Addresses:    addrs,
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
