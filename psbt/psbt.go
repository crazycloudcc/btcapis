// Package psbt PSBT处理
package psbt

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// Version 常量：支持 v0(BIP174) / v2(BIP370)
const (
	VersionV0 = 0
	VersionV2 = 2
)

// Packet 表示一个（简化版）PSBT 数据结构，覆盖 v0/v2 关键字段以满足本库的角色操作。
// 注：此结构关注“功能”与“互操作”，序列化/反序列化可后续扩展。
type Packet struct {
	// 版本：0 或 2
	Version int

	// v0：未签名交易的完整结构
	UnsignedTx *wire.MsgTx

	// v2：交易元数据（不直接包含完整 MsgTx）
	TxVersion int32
	LockTime  uint32

	Inputs  []*Input
	Outputs []*Output
}

// Input 表示单个输入的 PSBT 附加数据
type Input struct {
	// 引用的前序输出（v0 可从 UnsignedTx 读取；v2 需要显式设置）
	PrevTxID  chainhash.Hash
	PrevIndex uint32
	Sequence  uint32

	// UTXO 信息（二选一或同时提供，遵循 BIP174）
	NonWitnessUtxo *wire.MsgTx
	WitnessUtxo    *wire.TxOut

	// 脚本与派生信息
	RedeemScript  []byte
	WitnessScript []byte
	BIP32         []BIP32Derivation

	// 签名相关
	SighashType uint32            // txscript.SigHashType
	PartialSigs map[string][]byte // pubkey(hex) -> sig (DER+hashtype 或 Schnorr+hashtype)

	// 最终化产物
	FinalScriptSig     []byte
	FinalScriptWitness wire.TxWitness

	// Taproot 脚本路径（tapscript）相关
	// 若使用脚本路径花费，应提供：TapLeafScript、TapControlBlock；
	// 可选提供：TapAnnex（首字节应为 0x50）、TapScriptStack（额外入栈元素，不含脚本与控制块）。
	TapLeafScript   []byte
	TapControlBlock []byte
	TapAnnex        []byte
	TapScriptStack  [][]byte
}

// Output 表示单个输出的 PSBT 附加数据
type Output struct {
	Value        int64
	ScriptPubKey []byte

	RedeemScript  []byte
	WitnessScript []byte
	BIP32         []BIP32Derivation
}

// BIP32Derivation 记录单钥的派生路径
type BIP32Derivation struct {
	PubKey      []byte   // 压缩公钥(33)
	Fingerprint [4]byte  // 父键指纹
	Path        []uint32 // 绝对或相对路径（m/.. 样式外部表示由调用方维护）
}

// NewV0FromUnsignedTx 创建 v0 PSBT：填入未签名交易；Inputs/Outputs 根据 UnsignedTx 初始化空 maps。
func NewV0FromUnsignedTx(unsigned *wire.MsgTx) *Packet {
	if unsigned == nil {
		return &Packet{Version: VersionV0}
	}
	p := &Packet{Version: VersionV0, UnsignedTx: unsigned}
	p.Inputs = make([]*Input, len(unsigned.TxIn))
	for i := range p.Inputs {
		in := unsigned.TxIn[i]
		p.Inputs[i] = &Input{
			PrevTxID:    in.PreviousOutPoint.Hash,
			PrevIndex:   in.PreviousOutPoint.Index,
			Sequence:    in.Sequence,
			PartialSigs: make(map[string][]byte),
		}
	}
	p.Outputs = make([]*Output, len(unsigned.TxOut))
	for i := range p.Outputs {
		out := unsigned.TxOut[i]
		p.Outputs[i] = &Output{Value: out.Value, ScriptPubKey: append([]byte(nil), out.PkScript...)}
	}
	return p
}

// NewV2 创建 v2 PSBT：设置元数据与 I/O 计数，并初始化空 maps。
func NewV2(txVersion int32, lockTime uint32, inCount, outCount int) *Packet {
	p := &Packet{Version: VersionV2, TxVersion: txVersion, LockTime: lockTime}
	p.Inputs = make([]*Input, inCount)
	for i := 0; i < inCount; i++ {
		p.Inputs[i] = &Input{PartialSigs: make(map[string][]byte)}
	}
	p.Outputs = make([]*Output, outCount)
	for i := 0; i < outCount; i++ {
		p.Outputs[i] = &Output{}
	}
	return p
}

// MustInput 返回第 i 个输入，越界报错。
func (p *Packet) MustInput(i int) *Input {
	if i < 0 || i >= len(p.Inputs) {
		panic(fmt.Errorf("psbt: input index out of range: %d", i))
	}
	return p.Inputs[i]
}

// MustOutput 返回第 i 个输出，越界报错。
func (p *Packet) MustOutput(i int) *Output {
	if i < 0 || i >= len(p.Outputs) {
		panic(fmt.Errorf("psbt: output index out of range: %d", i))
	}
	return p.Outputs[i]
}

// IsV0 返回是否 v0
func (p *Packet) IsV0() bool { return p != nil && p.Version == VersionV0 }

// IsV2 返回是否 v2
func (p *Packet) IsV2() bool { return p != nil && p.Version == VersionV2 }
