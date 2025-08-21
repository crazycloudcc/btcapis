package types

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
