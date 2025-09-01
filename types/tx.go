package types

// VarInt 序列化时采用比特币可变长度整型编码；此处仅作标注，实际可直接用 uint64 并在编解码层处理。
type VarInt = uint64

// OutPoint 唯一标识一个已存在的交易输出：<txid, vout>
type TxOutPoint struct {
	Hash  Hash32 // 上一笔交易的 txid（逻辑上大端显示；写入时小端序）
	Index uint32 // vout 索引，从 0 开始
}

// TxWitness 是 SegWit 的 per-input 见证栈；外层是项数量，内层每项为原始字节。
type TxWitness [][]byte

// Tx：完整的交易对象（不显式包含 marker/flag 字段；是否为 segwit 由 TxIn[].Witness 是否存在决定）
// 序列化顺序（无见证）：Version | vinCount | vin[...] | voutCount | vout[...] | LockTime
// 序列化顺序（有见证）：Version | 0x00 | 0x01 | vinCount | vin[...] | voutCount | vout[...] | witnesses(for each vin) | LockTime
type Tx struct {
	Version  int32   // 交易版本
	LockTime uint32  // 交易锁定时间
	TxIn     []TxIn  // 交易输入
	TxOut    []TxOut // 交易输出
	// 非序列化辅助字段（可选）：
	// CachedTxID  Hash32
	// CachedWtxID Hash32
}

// TxIn：交易输入
// - PreviousOutPoint：被花费的 UTXO 引用
// - ScriptSig：非隔离见证路径下的解锁脚本（如 P2PKH 的 <sig><pubkey> 等）
// - Sequence：nSequence；影响 RBF（<0xffffffff-1）与 CSV；默认 0xffffffff
// - Witness：隔离见证路径下的见证栈（P2WPKH/P2WSH/P2TR 等）
type TxIn struct {
	PreviousOutPoint TxOutPoint // 上一笔交易的输出点
	Sequence         uint32     // 交易序列号
	ScriptSig        []byte     // scriptSig
	Witness          TxWitness  // 若任一输入 Witness 非空，序列化需写入 marker/flag，并在所有 TxOut 之后写入全部 Witness
}

// TxOut：交易输出
// - Value：satoshi 数（int64，允许负数编码但务必在业务层校验非负）
// - PkScript：锁定脚本（scriptPubKey），如 P2PKH 的 OP_DUP OP_HASH160 <20b> OP_EQUALVERIFY OP_CHECKSIG
type TxOut struct {
	Value      int64  // satoshi
	PkScript   []byte // scriptPubKey
	ScriptType string // 可选, 用于显示 解析出的脚本类别
	Address    string // 可选, 用于显示 解析得到的人类可读地址
}

// UTXO：钱包/索引层常用的未花费输出结构（**链上共识并不定义该结构**，这是应用层抽象）
type TxUTXO struct {
	OutPoint TxOutPoint // 定位该 UTXO
	Value    int64      // satoshi
	PkScript []byte     // 原始 scriptPubKey
	Height   uint32     // 产出该 UTXO 的区块高度；mempool 可置 0 或特约定值
	Coinbase bool       // 该 UTXO 是否来自 coinbase 交易
	Address  string     // 解析得到的人类可读地址（可选）
	// Class    ScriptClass // 解析出的脚本类别（可选）
	// 额外可选元数据（根据业务需要扩展）：
	// Confirmations uint32
	// ScriptVersion uint16   // 保留位；目前主网 script 版本固定
	// Timestamp     int64    // 区块时间或首次见到时间（秒）
	// Spendable     bool     // 钱包层权限或策略控制
}

// 通用的转账交易输入参数
type TxInputParams struct {
	FromAddress   []string  `json:"from_address"`   // 来源地址数组-可以是多个, 但是目前版本只支持1个地址
	ToAddress     []string  `json:"to_address"`     // 目标地址数组-可以是多个, 但是要和Amount一一对应
	AmountBTC     []float64 `json:"amount"`         // 金额-单位BTC
	FeeRate       float64   `json:"fee_rate"`       // 费用率(sat/vB)
	Locktime      int64     `json:"locktime"`       // 锁定时间(秒)
	Replaceable   bool      `json:"replaceable"`    // 是否可替换RBF
	Data          string    `json:"data"`           // 可选 交付附加数据
	PublicKey     string    `json:"public_key"`     // 公钥 => 从OKX获取, 后续要删除, 改用其他方式录入钱包
	ChangeAddress string    `json:"change_address"` // 找零地址
}
