// rpc原生数据结构定义
package bitcoindrpc

type UTXODTO struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	Desc         string  `json:"desc"`
	AmountBTC    float64 `json:"amount"`
	Height       int64   `json:"height"`
}

type scanResult struct {
	Success       bool      `json:"success"`
	Bestblock     string    `json:"bestblock"`
	Height        int64     `json:"height"`
	SearchedItems int64     `json:"searched_items"`
	Unspents      []UTXODTO `json:"unspents"`
}
