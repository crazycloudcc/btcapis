package types

// TxVerbose Bitcoin Core / mempool.space verbose 交易 JSON（与 apis DataTx 字段对齐）
type TxVerbose struct {
	TxId          string           `json:"txid"`
	Hash          string           `json:"hash"`
	Version       int64            `json:"version"`
	Size          int64            `json:"size"`
	Vsize         int64            `json:"vsize"`
	Weight        int64            `json:"weight"`
	Locktime      int64            `json:"locktime"`
	Vin           []TxVerboseInput `json:"vin"`
	Vout          []TxVerboseOutput `json:"vout"`
	Fee           float64          `json:"fee"`
	Hex           string           `json:"hex"`
	BlockHash     string           `json:"blockhash"`
	Confirmations int64            `json:"confirmations"`
	Time          int64            `json:"time"`
	Blocktime     int64            `json:"blocktime"`
}

// TxVerboseInput verbose 交易输入
type TxVerboseInput struct {
	Coinbase    string              `json:"coinbase"`
	TxId        string              `json:"txid"`
	Vout        int64               `json:"vout"`
	ScriptSig   TxVerboseScriptSig  `json:"scriptSig"`
	TxInWitness []string            `json:"txinwitness"`
	PrevOut     TxVerbosePrevOut    `json:"prevout"`
	Sequence    int64               `json:"sequence"`
}

// TxVerboseOutput verbose 交易输出
type TxVerboseOutput struct {
	Value        float64              `json:"value"`
	N            int64                `json:"n"`
	ScriptPubKey TxVerboseScriptPubKey `json:"scriptPubKey"`
}

// TxVerboseScriptSig 输入脚本
type TxVerboseScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// TxVerbosePrevOut 前序输出
type TxVerbosePrevOut struct {
	Generated    bool                 `json:"generated"`
	Height       int64                `json:"height"`
	Value        float64              `json:"value"`
	ScriptPubKey TxVerboseScriptPubKey `json:"scriptPubKey"`
}

// TxVerboseScriptPubKey 输出脚本
type TxVerboseScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

// AddressBalanceSats 地址余额（聪）
type AddressBalanceSats struct {
	Confirmed   float64 `json:"confirmed"`
	Unconfirmed float64 `json:"unconfirmed"`
	Total       float64 `json:"total"`
	Address     string  `json:"address"`
}

// AddressUTXO 地址 UTXO（统一结构，兼容多数据源字段）
type AddressUTXO struct {
	Value    int64  `json:"value"`
	TxId     string `json:"txid"`
	Vout     uint32 `json:"vout"`
	TxHash   string `json:"tx_hash"`
	TxPos    uint32 `json:"tx_pos"`
	Height   int64  `json:"height"`
	Confirmed   bool  `json:"confirmed"`
	BlockHeight int64 `json:"block_height"`
}

// AddressHistoryTx 地址交易历史摘要
type AddressHistoryTx struct {
	TxId   string `json:"txid"`
	TxHash string `json:"txhash"`
	Height int64  `json:"height"`
}

// AddressUnconfirmedTx 地址内存池交易摘要（electrumX）
type AddressUnconfirmedTx struct {
	TxHash string  `json:"txhash"`
	Height int64   `json:"height"`
	Fee    float64 `json:"fee"`
}
