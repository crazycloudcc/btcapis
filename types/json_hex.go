package types

import (
	"encoding/hex"
	"encoding/json"
)

// —— 小工具
func hexOf(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return hex.EncodeToString(b)
}
func hexList(bb [][]byte) []string {
	if len(bb) == 0 {
		return []string{}
	}
	out := make([]string, len(bb))
	for i, b := range bb {
		out[i] = hex.EncodeToString(b)
	}
	return out
}

// —— JSON 视图结构（仅用于 Marshal，不影响内部字段类型）

type txInJSON struct {
	TxID      string   `json:"TxID"`
	Vout      uint32   `json:"Vout"`
	Sequence  uint32   `json:"Sequence"`
	ScriptSig string   `json:"ScriptSig"` // hex
	Witness   []string `json:"Witness"`   // hex
}

type txOutJSON struct {
	Value        int64    `json:"Value"`
	ScriptPubKey string   `json:"ScriptPubKey"` // hex
	Type         string   `json:"Type"`
	Addresses    []string `json:"Addresses,omitempty"`
}

// —— 自定义 Marshal（只改输出，不改结构体定义）

func (in TxIn) MarshalJSON() ([]byte, error) {
	view := txInJSON{
		TxID:      in.TxID,
		Vout:      in.Vout,
		Sequence:  in.Sequence,
		ScriptSig: hexOf(in.ScriptSig),
		Witness:   hexList(in.Witness),
	}
	return json.Marshal(view)
}

func (out TxOut) MarshalJSON() ([]byte, error) {
	view := txOutJSON{
		Value:        out.Value,
		ScriptPubKey: hexOf(out.ScriptPubKey),
		Type:         out.Type,
		Addresses:    out.Addresses,
	}
	return json.Marshal(view)
}
