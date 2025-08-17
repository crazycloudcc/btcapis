package types

// ScriptOp 表示一条脚本指令；若为数据推送，DataHex/DataLen 会被填充。
type ScriptOp struct {
	Op      string `json:"op"`
	DataHex string `json:"data_hex,omitempty"`
	DataLen int    `json:"data_len,omitempty"`
}
