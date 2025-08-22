// rpc原生数据结构定义
package bitcoindrpc

// UTXO数据结构
type UTXODTO struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	Desc         string  `json:"desc"`
	AmountBTC    float64 `json:"amount"`
	Height       int64   `json:"height"`
}

// 扫描UTXO结果数据结构
type scanResult struct {
	Success       bool      `json:"success"`
	Bestblock     string    `json:"bestblock"`
	Height        int64     `json:"height"`
	SearchedItems int64     `json:"searched_items"`
	Unspents      []UTXODTO `json:"unspents"`
}

// 网络信息数据结构
type NetworkInfoDTO struct {
	Version         int    `json:"version"`
	Subversion      string `json:"subversion"`
	ProtocolVersion int    `json:"protocolversion"`
	LocalServices   string `json:"localservices"`
	LocalRelay      bool   `json:"localrelay"`
	TimeOffset      int    `json:"timeoffset"`
	Connections     int    `json:"connections"`
}

// 区块高度数据结构
type BlockCountDTO struct {
	Result int64 `json:"result"`
}

// 估算交易费率数据结构
type FeeRateDTO struct {
	Feerate float64  `json:"feerate"` // BTC/KB
	Errors  []string `json:"errors"`
}
