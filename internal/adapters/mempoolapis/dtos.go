// mempool.space 原生数据结构定义
package mempoolapis

// 从mempool.space获取到的交易数据结构
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

// 从mempool.space获取到的UTXO数据结构
type UTXODTO struct {
	Txid   string `json:"txid"`
	Vout   uint32 `json:"vout"`
	Value  int64  `json:"value"`
	Status struct {
		Confirmed   bool  `json:"confirmed"`
		BlockHeight int64 `json:"block_height"`
	} `json:"status"`
}
