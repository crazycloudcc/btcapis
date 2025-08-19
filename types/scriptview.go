package types

// ScriptOp 表示一条脚本指令；若为数据推送，DataHex/DataLen 会被填充。
type ScriptOp struct {
	Op      string `json:"op"`                 // 操作码
	DataHex string `json:"data_hex,omitempty"` // 数据十六进制
	DataLen int    `json:"data_len,omitempty"` // 数据长度
}
