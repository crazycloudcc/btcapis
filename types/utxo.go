package types

// UTXO 统一的未花费输出表示，供各 provider 返回给上层使用。
// 这里只放通用字段；如需更多信息（确认数、是否 coinbase、区块高度等），
// 可按需追加可选字段，保持向后兼容。
type UTXO struct {
	OutPoint     OutPoint // 交易哈希 + vout 索引
	Value        int64    // 金额（sats）
	ScriptPubKey []byte   // 输出脚本（spk）

	// 可选元数据（provider 有就填；没有就留默认值）
	Height        int64 // 包含此 UTXO 的区块高度；未确认可为 0 或 -1
	Confirmations int64 // 确认数；未确认为 0
	Coinbase      bool  // 是否来自 coinbase 交易
}
