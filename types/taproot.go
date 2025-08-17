package types

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
