// Package types 交易相关类型定义
package types

import (
	"time"
)

// Transaction 表示比特币交易
type Transaction struct {
	TxID        string     `json:"txid"`
	Version     int32      `json:"version"`
	LockTime    uint32     `json:"lock_time"`
	Inputs      []TxInput  `json:"inputs"`
	Outputs     []TxOutput `json:"outputs"`
	Size        int        `json:"size"`
	Weight      int        `json:"weight"`
	Fee         int64      `json:"fee"`
	BlockHeight int64      `json:"block_height"`
	BlockHash   string     `json:"block_hash"`
	BlockTime   time.Time  `json:"block_time"`
	Confirmed   bool       `json:"confirmed"`
}

// TxInput 表示交易输入
type TxInput struct {
	OutPoint     OutPoint `json:"outpoint"`
	ScriptSig    []byte   `json:"script_sig"`
	Witness      [][]byte `json:"witness"`
	Sequence     uint32   `json:"sequence"`
	Value        int64    `json:"value"`
	ScriptPubKey []byte   `json:"script_pub_key"`
}

// TxOutput 表示交易输出
type TxOutput struct {
	Value        int64  `json:"value"`
	ScriptPubKey []byte `json:"script_pub_key"`
	Address      string `json:"address"`
	Spent        bool   `json:"spent"`
}

// SigHashType 表示签名哈希类型
type SigHashType uint32

const (
	SigHashAll          SigHashType = 0x01
	SigHashNone         SigHashType = 0x02
	SigHashSingle       SigHashType = 0x03
	SigHashAnyoneCanPay SigHashType = 0x80
)

// BlockHeader 表示区块头
type BlockHeader struct {
	Hash             string    `json:"hash"`
	Version          int32     `json:"version"`
	PreviousHash     string    `json:"previous_hash"`
	MerkleRoot       string    `json:"merkle_root"`
	Timestamp        time.Time `json:"timestamp"`
	Bits             uint32    `json:"bits"`
	Nonce            uint32    `json:"nonce"`
	Height           int64     `json:"height"`
	Size             int       `json:"size"`
	Weight           int       `json:"weight"`
	TransactionCount int       `json:"transaction_count"`
}

// Block 表示完整区块
type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"transactions"`
}
