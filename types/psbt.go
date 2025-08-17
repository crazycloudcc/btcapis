// Package types PSBT相关类型定义
// PSBT (Partially Signed Bitcoin Transaction) 是比特币交易的标准格式，
// 允许交易在多个参与者之间传递，逐步完成签名过程
package types

import (
	"time"
)

// PSBT 表示部分签名比特币交易
// 这是PSBT的核心结构，包含了交易的全局信息、输入输出详情以及元数据
type PSBT struct {
	// GlobalTx 全局交易信息，包含交易的版本、锁定时间等基本信息
	GlobalTx *Transaction `json:"global_tx"`
	// Inputs PSBT输入数组，每个输入包含UTXO信息和签名数据
	Inputs []PSBTInput `json:"inputs"`
	// Outputs PSBT输出数组，每个输出包含地址和金额信息
	Outputs []PSBTOutput `json:"outputs"`
	// Version PSBT版本号，当前标准版本为0
	Version uint32 `json:"version"`
	// Unknown 未知的全局字段，用于扩展性
	Unknown []PSBTUnknown `json:"unknown"`
	// CreatedAt PSBT创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt PSBT最后更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// PSBTInput 表示PSBT输入
// 包含UTXO信息、签名数据、脚本等输入相关的所有信息
type PSBTInput struct {
	// NonWitnessUtxo 非见证UTXO的完整交易信息
	// 用于非隔离见证输入，包含完整的交易数据
	NonWitnessUtxo *Transaction `json:"non_witness_utxo"`
	// WitnessUtxo 见证UTXO信息
	// 用于隔离见证输入，只包含输出信息
	WitnessUtxo *TxOutput `json:"witness_utxo"`
	// PartialSigs 部分签名数组
	// 存储已收集的部分签名，每个签名对应一个公钥
	PartialSigs []PartialSig `json:"partial_sigs"`
	// SighashType 签名哈希类型
	// 定义签名使用的哈希算法类型（如SIGHASH_ALL等）
	SighashType SigHashType `json:"sighash_type"`
	// RedeemScript 赎回脚本
	// 用于P2SH输入的赎回脚本
	RedeemScript []byte `json:"redeem_script"`
	// WitnessScript 见证脚本
	// 用于P2WSH输入的见证脚本
	WitnessScript []byte `json:"witness_script"`
	// Bip32Derivation BIP32派生路径信息
	// 用于确定性钱包的密钥派生路径
	Bip32Derivation []BIP32Derivation `json:"bip32_derivation"`
	// FinalScriptSig 最终的脚本签名
	// 完成签名后的最终脚本签名
	FinalScriptSig []byte `json:"final_script_sig"`
	// FinalScriptWitness 最终的见证脚本
	// 完成签名后的最终见证脚本
	FinalScriptWitness [][]byte `json:"final_script_witness"`
	// Ripemd160Hashes RIPEMD160哈希值数组
	// 用于脚本哈希验证
	Ripemd160Hashes []byte `json:"ripemd160_hashes"`
	// Sha256Hashes SHA256哈希值数组
	// 用于脚本哈希验证
	Sha256Hashes []byte `json:"sha256_hashes"`
	// Hash160Hashes Hash160哈希值数组（RIPEMD160(SHA256())）
	// 用于P2PKH地址验证
	Hash160Hashes []byte `json:"hash160_hashes"`
	// Hash256Hashes Hash256哈希值数组（SHA256(SHA256())）
	// 用于交易ID和区块头哈希
	Hash256Hashes []byte `json:"hash256_hashes"`
	// Unknown 未知的输入字段，用于扩展性
	Unknown []PSBTUnknown `json:"unknown"`
}

// PSBTOutput 表示PSBT输出
// 包含输出相关的脚本和派生路径信息
type PSBTOutput struct {
	// RedeemScript 赎回脚本
	// 用于P2SH输出的赎回脚本
	RedeemScript []byte `json:"redeem_script"`
	// WitnessScript 见证脚本
	// 用于P2WSH输出的见证脚本
	WitnessScript []byte `json:"witness_script"`
	// Bip32Derivation BIP32派生路径信息
	// 用于确定性钱包的密钥派生路径
	Bip32Derivation []BIP32Derivation `json:"bip32_derivation"`
	// Unknown 未知的输出字段，用于扩展性
	Unknown []PSBTUnknown `json:"unknown"`
}

// PartialSig 表示部分签名
// 存储公钥和对应的部分签名数据
type PartialSig struct {
	// PubKey 公钥字节数组
	// 用于验证签名的公钥
	PubKey []byte `json:"pub_key"`
	// Signature 签名字节数组
	// 对应公钥的签名数据
	Signature []byte `json:"signature"`
}

// BIP32Derivation 表示BIP32派生信息
// 用于确定性钱包的密钥派生路径管理
type BIP32Derivation struct {
	// PubKey 公钥字节数组
	// 派生出的公钥
	PubKey []byte `json:"pub_key"`
	// MasterFingerprint 主密钥指纹
	// 主密钥的指纹，用于标识密钥来源
	MasterFingerprint []byte `json:"master_fingerprint"`
	// Path 派生路径
	// BIP32派生路径，如 [44', 0', 0', 0, 0]
	Path []uint32 `json:"path"`
}

// PSBTUnknown 表示未知的PSBT字段
// 用于存储PSBT标准中未定义的扩展字段
type PSBTUnknown struct {
	// Key 字段键名
	// 未知字段的标识符
	Key []byte `json:"key"`
	// Value 字段值
	// 未知字段的数据内容
	Value []byte `json:"value"`
}

// PSBTStatus 表示PSBT状态
// 用于跟踪PSBT在签名流程中的当前状态
type PSBTStatus string

const (
	// PSBTStatusCreated PSBT已创建状态
	// 初始状态，PSBT刚被创建
	PSBTStatusCreated PSBTStatus = "created"
	// PSBTStatusSigning PSBT签名中状态
	// 正在进行签名收集过程
	PSBTStatusSigning PSBTStatus = "signing"
	// PSBTStatusSigned PSBT已签名状态
	// 所有必要的签名已完成
	PSBTStatusSigned PSBTStatus = "signed"
	// PSBTStatusFinalized PSBT已最终化状态
	// 所有脚本已最终化，可以广播
	PSBTStatusFinalized PSBTStatus = "finalized"
	// PSBTStatusBroadcast PSBT已广播状态
	// 交易已广播到网络
	PSBTStatusBroadcast PSBTStatus = "broadcast"
)
