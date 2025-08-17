// Package types 交易相关类型定义
package types

// types/tx.go
type Tx struct {
	TxID     string
	Version  int32
	LockTime uint32
	Weight   int64
	Vsize    int64

	Vin  []TxIn
	Vout []TxOut
}

type TxIn struct {
	TxID      string
	Vout      uint32
	Sequence  uint32
	ScriptSig []byte
	Witness   [][]byte
	// 可选：解析出的 scriptSig/witness 语义字段
}

type TxOut struct {
	Value        int64 // sats
	ScriptPubKey []byte
	Type         string // p2pkh/p2sh/p2wpkh/p2wsh/p2tr/...
	Addresses    []string
}

type OutPoint struct {
	Hash string
	N    uint32
}
