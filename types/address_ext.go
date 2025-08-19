// Package btcscripts: "地址 → 脚本模板与花费步骤" 常量库
// 适配 btcsuite 生态（btcd/txscript、btcutil、chaincfg）。
//
// 功能概览：
// 1) 按地址生成锁定脚本 scriptPubKey（含 P2PKH/P2SH/P2WPKH/P2WSH/P2TR）。
// 2) 各地址类型的 OP 列表模板（锁定/解锁）与要点说明（TemplateMap）。
// 3) 常用脚本构造：P2WSH scriptPubKey、P2TR scriptPubKey、P2SH 的 redeemScript（P2WPKH/P2WSH 包装）。
// 4) 示例解锁数据生成（scriptSig / witness），便于构造与调试。
// 5) 复杂脚本示例：2-of-3 多签、CLTV/CSV、HTLC tapscript 生成函数。
//
// 注意：
// - 仅由 P2SH 地址本身无法区分是否 P2SH-P2WPKH/P2WSH；需要提供 redeemScript。
// - Taproot Script-path 需要调用方提供 tapscript 与 control block；本库只负责组织 witness。
// - 本库不实现 Schnorr/ECDSA 签名，仅接受外部生成的签名原始字节。
package types

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// Template 描述锁定/解锁模板（OP 列表与说明）。
type Template struct {
	Name         string
	LockingOPs   []string // scriptPubKey 的 OP 序列（描述）
	UnlockingOPs []string // scriptSig/witness 的形态描述
	Notes        string   // 重要说明
}

// TemplateMap: 各地址类型的模板速查。
var TemplateMap = map[AddressType]Template{
	AddrP2PK: {
		Name:       "P2PK",
		LockingOPs: []string{"<pubkey>", "OP_CHECKSIG"},
		UnlockingOPs: []string{
			"scriptSig: <signature>",
		},
		Notes: "早期形态；现极少用。",
	},
	AddrP2PKH: {
		Name:       "P2PKH (1…)",
		LockingOPs: []string{"OP_DUP", "OP_HASH160", "<20B pubKeyHash>", "OP_EQUALVERIFY", "OP_CHECKSIG"},
		UnlockingOPs: []string{
			"scriptSig: <signature> <pubkey>",
		},
		Notes: "传统最常见；可塑性高；Base58Check。",
	},
	AddrP2SH: {
		Name:       "P2SH (3…)",
		LockingOPs: []string{"OP_HASH160", "<20B HASH160(redeemScript)>", "OP_EQUAL"},
		UnlockingOPs: []string{
			"scriptSig: <arg1> .. <argk> <redeemScript>",
		},
		Notes: "脚本原文锁时隐藏，花费披露；仅地址无法区分是否 P2SH-P2WPKH/P2WSH。",
	},
	AddrP2WPKH: {
		Name:       "P2WPKH (bc1q…, v0/20)",
		LockingOPs: []string{"OP_0", "<20B HASH160(pubkey)>"},
		UnlockingOPs: []string{
			"scriptSig: 空",
			"witness: [<signature>, <pubkey>]",
		},
		Notes: "SegWit v0；BIP-143 摘要；低费、抗可塑性。",
	},
	AddrP2WSH: {
		Name:       "P2WSH (bc1q…, v0/32)",
		LockingOPs: []string{"OP_0", "<32B SHA256(witnessScript)>"},
		UnlockingOPs: []string{
			"scriptSig: 空",
			"witness: [<args…>, <witnessScript>] (最后一项是脚本原文)",
		},
		Notes: "复杂脚本优选；BIP-143 摘要。",
	},
	AddrP2TR: {
		Name:       "P2TR (bc1p…, v1/32)",
		LockingOPs: []string{"OP_1", "<32B x-only output key Q>"},
		UnlockingOPs: []string{
			"Key-path witness: [<schnorr sig>(+opt sighash), annex?]",
			"Script-path witness: [<args…>, annex?, <tapscript>, <control_block>]",
		},
		Notes: "Taproot；BIP-341/342；Schnorr；Bech32m。",
	},
}

// DetectAddressType 解码地址并返回地址类型与 btcutil.Address。
func DetectAddressType(addrStr string, params *chaincfg.Params) (AddressType, btcutil.Address, error) {
	addr, err := btcutil.DecodeAddress(addrStr, params)
	if err != nil {
		return "", nil, err
	}
	switch a := addr.(type) {
	case *btcutil.AddressPubKeyHash:
		return AddrP2PKH, a, nil
	case *btcutil.AddressScriptHash:
		return AddrP2SH, a, nil
	case *btcutil.AddressWitnessPubKeyHash:
		return AddrP2WPKH, a, nil
	case *btcutil.AddressWitnessScriptHash:
		return AddrP2WSH, a, nil
	case *btcutil.AddressTaproot:
		return AddrP2TR, a, nil
	default:
		return "", a, fmt.Errorf("unknown address type: %T", addr)
	}
}

// ScriptPubKeyFromAddress 通过地址生成锁定脚本（scriptPubKey）与类型。
func ScriptPubKeyFromAddress(addrStr string, params *chaincfg.Params) ([]byte, AddressType, error) {
	kind, addr, err := DetectAddressType(addrStr, params)
	if err != nil {
		return nil, "", err
	}
	spk, err := txscript.PayToAddrScript(addr)
	return spk, kind, err
}

// BuildRedeemScriptP2WPKH: redeemScript = OP_0 <20B pkh>
func BuildRedeemScriptP2WPKH(pubKeyHash20 []byte) ([]byte, error) {
	if len(pubKeyHash20) != 20 {
		return nil, fmt.Errorf("pkh length != 20")
	}
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash20).Script()
}

// BuildRedeemScriptP2WSH: redeemScript = OP_0 <32B sha256(witnessScript)>
func BuildRedeemScriptP2WSH(witnessScript []byte) ([]byte, error) {
	h := sha256.Sum256(witnessScript)
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(h[:]).Script()
}

// BuildP2WSHScriptPubKey: OP_0 <32B sha256(witnessScript)>
func BuildP2WSHScriptPubKey(witnessScript []byte) ([]byte, error) {
	h := sha256.Sum256(witnessScript)
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(h[:]).Script()
}

// BuildP2TRScriptPubKey: OP_1 <32B x-only Q>
func BuildP2TRScriptPubKey(xOnlyQ32 []byte) ([]byte, error) {
	if len(xOnlyQ32) != 32 {
		return nil, fmt.Errorf("x-only pubkey must be 32 bytes")
	}
	return txscript.NewScriptBuilder().AddOp(txscript.OP_1).AddData(xOnlyQ32).Script()
}

// UnlockExample 描述一次花费所需的 scriptSig 与 witness（其一或两者）。
type UnlockExample struct {
	ScriptSig []byte   // 传统脚本，SegWit/Taproot 原生为空
	Witness   [][]byte // 见证栈，自底到顶顺序（wire.TxWitness 同序）
}

// --- 示例解锁构造器 -------------------------------------------------------

// P2PKH: scriptSig = <sig> <pubkey>
func ExampleUnlockP2PKH(sigDERPlusHashType, pubKey []byte) (UnlockExample, error) {
	ss, err := txscript.NewScriptBuilder().AddData(sigDERPlusHashType).AddData(pubKey).Script()
	return UnlockExample{ScriptSig: ss}, err
}

// P2SH: scriptSig = <args...> <redeemScript>
func ExampleUnlockP2SH(args [][]byte, redeemScript []byte) (UnlockExample, error) {
	b := txscript.NewScriptBuilder()
	for _, a := range args {
		b.AddData(a)
	}
	b.AddData(redeemScript)
	ss, err := b.Script()
	return UnlockExample{ScriptSig: ss}, err
}

// P2WPKH: witness = [<sig>, <pubkey>]
func ExampleUnlockP2WPKH(sigDERPlusHashType, pubKey []byte) UnlockExample {
	return UnlockExample{Witness: wire.TxWitness{sigDERPlusHashType, pubKey}}
}

// P2WSH: witness = [<args...>, <witnessScript>] （最后一项是脚本原文）
func ExampleUnlockP2WSH(args [][]byte, witnessScript []byte) UnlockExample {
	w := make(wire.TxWitness, 0, len(args)+1)
	for _, a := range args {
		w = append(w, a)
	}
	w = append(w, witnessScript)
	return UnlockExample{Witness: w}
}

// P2SH-P2WPKH: scriptSig=<redeemScript>; witness 与 P2WPKH 相同
func ExampleUnlockP2SHP2WPKH(redeemScript, sigDERPlusHashType, pubKey []byte) (UnlockExample, error) {
	ss, err := txscript.NewScriptBuilder().AddData(redeemScript).Script()
	return UnlockExample{ScriptSig: ss, Witness: wire.TxWitness{sigDERPlusHashType, pubKey}}, err
}

// P2SH-P2WSH: scriptSig=<redeemScript>; witness 与 P2WSH 相同
func ExampleUnlockP2SHP2WSH(redeemScript []byte, args [][]byte, witnessScript []byte) (UnlockExample, error) {
	ss, err := txscript.NewScriptBuilder().AddData(redeemScript).Script()
	w := make(wire.TxWitness, 0, len(args)+1)
	for _, a := range args {
		w = append(w, a)
	}
	w = append(w, witnessScript)
	return UnlockExample{ScriptSig: ss, Witness: w}, err
}

// P2TR-KeyPath: witness = [<schnorrSig>(+opt sighash), annex?]
// 注意：sig 可为 64B（SIGHASH_DEFAULT）或 65B（附 1 字节类型）。
func ExampleUnlockP2TRKeyPath(sig []byte, annexOptional []byte) UnlockExample {
	w := wire.TxWitness{sig}
	if len(annexOptional) > 0 {
		w = append(w, annexOptional) // 必须以 0x50 开头（调用方保证）
	}
	return UnlockExample{Witness: w}
}

// P2TR-ScriptPath: witness = [<args...>, annex?, <tapscript>, <control_block>]
func ExampleUnlockP2TRScriptPath(args [][]byte, tapscript []byte, controlBlock []byte, annexOptional []byte) UnlockExample {
	w := make(wire.TxWitness, 0, len(args)+2)
	for _, a := range args {
		w = append(w, a)
	}
	if len(annexOptional) > 0 {
		w = append(w, annexOptional) // 必须 0x50 开头
	}
	w = append(w, tapscript)
	w = append(w, controlBlock)
	return UnlockExample{Witness: w}
}

// --- 常见 witnessScript 生成（P2WSH/Tapscript） ---------------------------

// WitnessScript_2of3: 2-of-3 多签（SegWit v0 用于 P2WSH）。公钥建议按 BIP-67 排序（非共识）。
func WitnessScript_2of3(pub1, pub2, pub3 []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_2).
		AddData(pub1).AddData(pub2).AddData(pub3).
		AddOp(txscript.OP_3).
		AddOp(txscript.OP_CHECKMULTISIG).
		Script()
}

// WitnessScript_CLTV_Singlesig: 绝对时间锁 + 单签
func WitnessScript_CLTV_Singlesig(locktime int64, pub []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddInt64(locktime).AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).AddOp(txscript.OP_DROP).
		AddData(pub).AddOp(txscript.OP_CHECKSIG).
		Script()
}

// WitnessScript_CSV_Singlesig: 相对时间锁 + 单签（交易 version>=2，nSequence 满足 csvDelay）
func WitnessScript_CSV_Singlesig(csvDelay int64, pub []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddInt64(csvDelay).AddOp(txscript.OP_CHECKSEQUENCEVERIFY).AddOp(txscript.OP_DROP).
		AddData(pub).AddOp(txscript.OP_CHECKSIG).
		Script()
}

// WitnessScript_Hashlock_Singlesig (SHA256 版)：给出 R + 签名
func WitnessScript_Hashlock_Singlesig(H32 []byte, pub []byte) ([]byte, error) {
	if len(H32) != 32 {
		return nil, errors.New("H must be 32 bytes (SHA256)")
	}
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_SHA256).AddData(H32).AddOp(txscript.OP_EQUALVERIFY).
		AddData(pub).AddOp(txscript.OP_CHECKSIG).
		Script()
}

// Tapscript_HTLC 简化版（Script-path）：成功(收款人+R) / 退款(发送人+CLTV)
// 注意：这是 Tapscript 文本；用于 P2TR Script-path。
func Tapscript_HTLC(H32 []byte, recvPubXOnly []byte, locktime int64, sendPubXOnly []byte) ([]byte, error) {
	if len(H32) != 32 {
		return nil, errors.New("H must be 32 bytes (SHA256)")
	}
	b := txscript.NewScriptBuilder()
	b.AddOp(txscript.OP_IF)
	// 成功路径：收款人签名 + R 校验
	b.AddData(recvPubXOnly).AddOp(txscript.OP_CHECKSIGVERIFY)
	b.AddOp(txscript.OP_SIZE).AddInt64(32).AddOp(txscript.OP_EQUALVERIFY)
	b.AddOp(txscript.OP_SHA256).AddData(H32).AddOp(txscript.OP_EQUAL)
	b.AddOp(txscript.OP_ELSE)
	// 超时退款
	b.AddInt64(locktime).AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).AddOp(txscript.OP_DROP)
	b.AddData(sendPubXOnly).AddOp(txscript.OP_CHECKSIG)
	b.AddOp(txscript.OP_ENDIF)
	return b.Script()
}

// --- PSBT 相关小助手（签名侧需要的前序输出） ----------------------------

// PSBTInputForWitnessUTXO 生成 PSBT 里建议放入的 witness_utxo（供签名摘要用）。
func PSBTInputForWitnessUTXO(valueSats int64, scriptPubKey []byte) *wire.TxOut {
	return &wire.TxOut{Value: valueSats, PkScript: scriptPubKey}
}

// --- 使用示例 --------------------------------------------------------------
//
// 1) 从地址生成锁定脚本：
//    spk, kind, err := ScriptPubKeyFromAddress("bc1q...", &chaincfg.MainNetParams)
//
// 2) P2WPKH 花费 witness：
//    ex := ExampleUnlockP2WPKH(sigDER1, pub1)
//
// 3) P2WSH 多签：
//    wscript, _ := WitnessScript_2of3(pub1, pub2, pub3)
//    spk, _ := BuildP2WSHScriptPubKey(wscript)
//    ex := ExampleUnlockP2WSH([][]byte{[]byte(""), sig1, sig3}, wscript) // 注意 dummy ""
//
// 4) P2SH 包装 P2WPKH：
//    redeem, _ := BuildRedeemScriptP2WPKH(pkh20)
//    ss_wit, _ := ExampleUnlockP2SHP2WPKH(redeem, sigDER1, pub1)
//
// 5) Taproot：
//    // Key-path：ex := ExampleUnlockP2TRKeyPath(sig64or65, nil)
//    // Script-path：ex := ExampleUnlockP2TRScriptPath(args, tapscript, controlBlock, nil)
