package types

// ScriptOp 表示一条脚本指令；若为数据推送，DataHex/DataLen 会被填充。
type ScriptOp struct {
	Op      string `json:"op"`
	DataHex string `json:"data_hex,omitempty"`
	DataLen int    `json:"data_len,omitempty"`
}

// OrdinalsRecord 表示一条 TLV 记录（key/value 都用 hex 表示；方便外部再做自定义解析）
type OrdinalsRecord struct {
	KeyHex   string `json:"key_hex"`
	ValueHex string `json:"value_hex"`
}

// OrdinalsEnvelope 为还原后的 envelope 数据
type OrdinalsEnvelope struct {
	ContentType string           `json:"content_type,omitempty"` // 若 key=0x01 存在，ASCII 解码
	BodyHex     string           `json:"body_hex,omitempty"`     // key=0x00 的所有分片合并后的十六进制
	Records     []OrdinalsRecord `json:"records"`                // 完整 TLV（有序）
}

// TapControlBlock 解析后的控制块信息（P2TR 脚本路径花费使用）
type TapControlBlock struct {
	Header       byte     `json:"header"`        // 原始头字节
	LeafVersion  byte     `json:"leaf_version"`  // header & 0xfe
	Parity       int      `json:"parity"`        // (header >> 7) & 1
	InternalKey  string   `json:"internal_key"`  // 32B x-only pubkey (hex)
	MerkleHashes []string `json:"merkle_hashes"` // 0或多段32B（hex）
}

// TapscriptInfo 表示 witness 中的脚本路径花费解析结果
type TapscriptInfo struct {
	ScriptHex string            `json:"script_hex"`
	ASM       string            `json:"asm"`
	Ops       []ScriptOp        `json:"ops"`
	Control   TapControlBlock   `json:"control_block"`
	StackHex  []string          `json:"stack_hex"` // 脚本执行前的栈元素（不含 script 与 control block）
	Path      string            `json:"path"`      // "p2tr-script", "p2tr-key", "p2wpkh", ...
	Ord       *OrdinalsEnvelope `json:"ord,omitempty"`
}
