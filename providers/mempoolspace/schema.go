// 第三方响应定义
package mempoolspace

type TxDTO struct {
	Txid     string `json:"txid"`
	Version  int32  `json:"version"`
	Locktime uint32 `json:"locktime"`
	Weight   int64  `json:"weight"`
	Size     int64  `json:"size"`

	Vin []struct {
		Txid      string   `json:"txid"`
		Vout      uint32   `json:"vout"`
		Sequence  uint32   `json:"sequence"`
		Scriptsig string   `json:"scriptsig"`
		Witness   []string `json:"witness"`
	} `json:"vin"`

	Vout []struct {
		Value        int64  `json:"value"` // sats
		ScriptPubKey string `json:"scriptpubkey"`
		ScriptType   string `json:"scriptpubkey_type"`
		Address      string `json:"scriptpubkey_address"`
	} `json:"vout"`
}
