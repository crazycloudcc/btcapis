// Package types PSBT相关类型定义
package types

import (
	"time"
)

// PSBT 表示部分签名比特币交易
type PSBT struct {
	GlobalTx  *Transaction  `json:"global_tx"`
	Inputs    []PSBTInput   `json:"inputs"`
	Outputs   []PSBTOutput  `json:"outputs"`
	Version   uint32        `json:"version"`
	Unknown   []PSBTUnknown `json:"unknown"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// PSBTInput 表示PSBT输入
type PSBTInput struct {
	NonWitnessUtxo     *Transaction      `json:"non_witness_utxo"`
	WitnessUtxo        *TxOutput         `json:"witness_utxo"`
	PartialSigs        []PartialSig      `json:"partial_sigs"`
	SighashType        SigHashType       `json:"sighash_type"`
	RedeemScript       []byte            `json:"redeem_script"`
	WitnessScript      []byte            `json:"witness_script"`
	Bip32Derivation    []BIP32Derivation `json:"bip32_derivation"`
	FinalScriptSig     []byte            `json:"final_script_sig"`
	FinalScriptWitness [][]byte          `json:"final_script_witness"`
	Ripemd160Hashes    []byte            `json:"ripemd160_hashes"`
	Sha256Hashes       []byte            `json:"sha256_hashes"`
	Hash160Hashes      []byte            `json:"hash160_hashes"`
	Hash256Hashes      []byte            `json:"hash256_hashes"`
	Unknown            []PSBTUnknown     `json:"unknown"`
}

// PSBTOutput 表示PSBT输出
type PSBTOutput struct {
	RedeemScript    []byte            `json:"redeem_script"`
	WitnessScript   []byte            `json:"witness_script"`
	Bip32Derivation []BIP32Derivation `json:"bip32_derivation"`
	Unknown         []PSBTUnknown     `json:"unknown"`
}

// PartialSig 表示部分签名
type PartialSig struct {
	PubKey    []byte `json:"pub_key"`
	Signature []byte `json:"signature"`
}

// BIP32Derivation 表示BIP32派生信息
type BIP32Derivation struct {
	PubKey            []byte   `json:"pub_key"`
	MasterFingerprint []byte   `json:"master_fingerprint"`
	Path              []uint32 `json:"path"`
}

// PSBTUnknown 表示未知的PSBT字段
type PSBTUnknown struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// PSBTStatus 表示PSBT状态
type PSBTStatus string

const (
	PSBTStatusCreated   PSBTStatus = "created"
	PSBTStatusSigning   PSBTStatus = "signing"
	PSBTStatusSigned    PSBTStatus = "signed"
	PSBTStatusFinalized PSBTStatus = "finalized"
	PSBTStatusBroadcast PSBTStatus = "broadcast"
)
