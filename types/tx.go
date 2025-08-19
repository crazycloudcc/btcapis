// // Package types 交易相关类型定义
package types

// // Tx 交易结构体
// type Tx struct {
// 	TxID     string // 交易哈希
// 	Version  int32  // 交易版本
// 	LockTime uint32 // 交易锁定时间
// 	Weight   int64  // 交易权重
// 	Vsize    int64  // 交易大小

// 	Inputs  []TxIn  // 交易输入
// 	Outputs []TxOut // 交易输出
// }

// // TxIn 交易输入结构体
// type TxIn struct {
// 	TxID      string   // 交易哈希
// 	Vout      uint32   // 交易输出索引
// 	Sequence  uint32   // 交易序列号
// 	ScriptSig []byte   // 交易输入脚本
// 	Witness   [][]byte // 交易见证

// 	// 可选：解析出的 scriptSig/witness 语义字段
// }

// // TxOut 交易输出结构体
// type TxOut struct {
// 	TxID         string   // 交易哈希
// 	Vout         uint32   // 交易输出索引
// 	Value        int64    // sats
// 	ScriptPubKey []byte   // 交易输出脚本
// 	Type         string   // p2pkh/p2sh/p2wpkh/p2wsh/p2tr/...
// 	Addresses    []string // 交易输出地址
// }

// // UTXO 统一的未花费输出表示，供各 provider 返回给上层使用。
// // 这里只放通用字段；如需更多信息（确认数、是否 coinbase、区块高度等），
// // 可按需追加可选字段，保持向后兼容。
// type UTXO struct {
// 	Output       TxOut  // 交易哈希 + vout 索引
// 	Value        int64  // 金额（sats）
// 	ScriptPubKey []byte // 输出脚本（spk）

// 	// 可选元数据（provider 有就填；没有就留默认值）
// 	Height        int64 // 包含此 UTXO 的区块高度；未确认可为 0 或 -1
// 	Confirmations int64 // 确认数；未确认为 0
// 	Coinbase      bool  // 是否来自 coinbase 交易

// 	// 用户创建交易时使用
// 	TxID       string // 交易哈希
// 	Vout       uint32 // 交易输出索引
// 	ScriptType string // p2pkh/p2sh/p2wpkh/p2wsh/p2tr/...
// }

// // Output 交易输出结构体
// type Output struct {
// 	Address string // 或者你也可以扩展支持 ScriptPubKey
// 	Value   int64  // sats
// }

// // BuildTxRequest 构建交易请求结构体
// type BuildTxRequest struct {
// 	Network         Network  // 网络
// 	Inputs          []UTXO   // 输入
// 	Outputs         []Output // 输出
// 	ChangeAddress   string   // 找零地址
// 	FeeRateSatPerVb float64  // <=0 时使用后端估算/兜底
// 	EnableRBF       bool     // 默认 true
// 	LockTime        uint32   // 默认为0
// 	MinChange       int64    // <该值就并入手续费（避免尘埃找零），默认 0
// 	DustLimit       int64    // 默认 546（可按脚本类型细化）
// 	SighashType     uint32   // 缺省 SIGHASH_ALL
// }

// // BuildTxResult 构建交易结果结构体
// type BuildTxResult struct {
// 	UnsignedTxHex string // 未签名交易
// 	PSBTBase64    string // 部分签名交易
// 	SelectedUTXO  []UTXO // 选中的输入
// 	VSizeEstimate int64  // 交易大小估计
// 	FeePaid       int64  // 手续费
// 	ChangeValue   int64  // 找零金额
// }
