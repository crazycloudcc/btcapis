// 示例代码备用参考: 锁定脚本/解锁脚本/自定义逻辑脚等
package types

// import (
// 	"crypto/sha256"
// 	"errors"
// 	"fmt"

// 	"github.com/btcsuite/btcd/btcutil"
// 	"github.com/btcsuite/btcd/chaincfg"
// 	"github.com/btcsuite/btcd/txscript"
// 	"github.com/btcsuite/btcd/wire"
// )

// // Template 描述锁定/解锁模板（OP 列表与说明）。
// type Template struct {
// 	Name         string
// 	LockingOPs   []string // scriptPubKey 的 OP 序列（描述）
// 	UnlockingOPs []string // scriptSig/witness 的形态描述
// 	Notes        string   // 重要说明
// }

// // TemplateMap: 各地址类型的模板速查。
// var TemplateMap = map[AddressType]Template{
// 	AddrP2PK: {
// 		Name:       "P2PK",
// 		LockingOPs: []string{"<pubkey>", "OP_CHECKSIG"},
// 		UnlockingOPs: []string{
// 			"scriptSig: <signature>",
// 		},
// 		Notes: "早期形态；现极少用。",
// 	},
// 	AddrP2PKH: {
// 		Name:       "P2PKH (1…)",
// 		LockingOPs: []string{"OP_DUP", "OP_HASH160", "<20B pubKeyHash>", "OP_EQUALVERIFY", "OP_CHECKSIG"},
// 		UnlockingOPs: []string{
// 			"scriptSig: <signature> <pubkey>",
// 		},
// 		Notes: "传统最常见；可塑性高；Base58Check。",
// 	},
// 	AddrP2SH: {
// 		Name:       "P2SH (3…)",
// 		LockingOPs: []string{"OP_HASH160", "<20B HASH160(redeemScript)>", "OP_EQUAL"},
// 		UnlockingOPs: []string{
// 			"scriptSig: <arg1> .. <argk> <redeemScript>",
// 		},
// 		Notes: "脚本原文锁时隐藏，花费披露；仅地址无法区分是否 P2SH-P2WPKH/P2WSH。",
// 	},
// 	AddrP2WPKH: {
// 		Name:       "P2WPKH (bc1q…, v0/20)",
// 		LockingOPs: []string{"OP_0", "<20B HASH160(pubkey)>"},
// 		UnlockingOPs: []string{
// 			"scriptSig: 空",
// 			"witness: [<signature>, <pubkey>]",
// 		},
// 		Notes: "SegWit v0；BIP-143 摘要；低费、抗可塑性。",
// 	},
// 	AddrP2WSH: {
// 		Name:       "P2WSH (bc1q…, v0/32)",
// 		LockingOPs: []string{"OP_0", "<32B SHA256(witnessScript)>"},
// 		UnlockingOPs: []string{
// 			"scriptSig: 空",
// 			"witness: [<args…>, <witnessScript>] (最后一项是脚本原文)",
// 		},
// 		Notes: "复杂脚本优选；BIP-143 摘要。",
// 	},
// 	AddrP2TR: {
// 		Name:       "P2TR (bc1p…, v1/32)",
// 		LockingOPs: []string{"OP_1", "<32B x-only output key Q>"},
// 		UnlockingOPs: []string{
// 			"Key-path witness: [<schnorr sig>(+opt sighash), annex?]",
// 			"Script-path witness: [<args…>, annex?, <tapscript>, <control_block>]",
// 		},
// 		Notes: "Taproot；BIP-341/342；Schnorr；Bech32m。",
// 	},
// }

// // DetectAddressType 解码地址并返回地址类型与 btcutil.Address。
// func DetectAddressType(addrStr string, params *chaincfg.Params) (AddressType, btcutil.Address, error) {
// 	addr, err := btcutil.DecodeAddress(addrStr, params)
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	switch a := addr.(type) {
// 	case *btcutil.AddressPubKeyHash:
// 		return AddrP2PKH, a, nil
// 	case *btcutil.AddressScriptHash:
// 		return AddrP2SH, a, nil
// 	case *btcutil.AddressWitnessPubKeyHash:
// 		return AddrP2WPKH, a, nil
// 	case *btcutil.AddressWitnessScriptHash:
// 		return AddrP2WSH, a, nil
// 	case *btcutil.AddressTaproot:
// 		return AddrP2TR, a, nil
// 	default:
// 		return "", a, fmt.Errorf("unknown address type: %T", addr)
// 	}
// }

// // ScriptPubKeyFromAddress 通过地址生成锁定脚本（scriptPubKey）与类型。
// func ScriptPubKeyFromAddress(addrStr string, params *chaincfg.Params) ([]byte, AddressType, error) {
// 	kind, addr, err := DetectAddressType(addrStr, params)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	spk, err := txscript.PayToAddrScript(addr)
// 	return spk, kind, err
// }

// // BuildRedeemScriptP2WPKH: redeemScript = OP_0 <20B pkh>
// func BuildRedeemScriptP2WPKH(pubKeyHash20 []byte) ([]byte, error) {
// 	if len(pubKeyHash20) != 20 {
// 		return nil, fmt.Errorf("pkh length != 20")
// 	}
// 	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash20).Script()
// }

// // BuildRedeemScriptP2WSH: redeemScript = OP_0 <32B sha256(witnessScript)>
// func BuildRedeemScriptP2WSH(witnessScript []byte) ([]byte, error) {
// 	h := sha256.Sum256(witnessScript)
// 	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(h[:]).Script()
// }

// // BuildP2WSHScriptPubKey: OP_0 <32B sha256(witnessScript)>
// func BuildP2WSHScriptPubKey(witnessScript []byte) ([]byte, error) {
// 	h := sha256.Sum256(witnessScript)
// 	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(h[:]).Script()
// }

// // BuildP2TRScriptPubKey: OP_1 <32B x-only Q>
// func BuildP2TRScriptPubKey(xOnlyQ32 []byte) ([]byte, error) {
// 	if len(xOnlyQ32) != 32 {
// 		return nil, fmt.Errorf("x-only pubkey must be 32 bytes")
// 	}
// 	return txscript.NewScriptBuilder().AddOp(txscript.OP_1).AddData(xOnlyQ32).Script()
// }

// // UnlockExample 描述一次花费所需的 scriptSig 与 witness（其一或两者）。
// type UnlockExample struct {
// 	ScriptSig []byte   // 传统脚本，SegWit/Taproot 原生为空
// 	Witness   [][]byte // 见证栈，自底到顶顺序（wire.TxWitness 同序）
// }

// // --- 示例解锁构造器 -------------------------------------------------------

// // P2PKH: scriptSig = <sig> <pubkey>
// func ExampleUnlockP2PKH(sigDERPlusHashType, pubKey []byte) (UnlockExample, error) {
// 	ss, err := txscript.NewScriptBuilder().AddData(sigDERPlusHashType).AddData(pubKey).Script()
// 	return UnlockExample{ScriptSig: ss}, err
// }

// // P2SH: scriptSig = <args...> <redeemScript>
// func ExampleUnlockP2SH(args [][]byte, redeemScript []byte) (UnlockExample, error) {
// 	b := txscript.NewScriptBuilder()
// 	for _, a := range args {
// 		b.AddData(a)
// 	}
// 	b.AddData(redeemScript)
// 	ss, err := b.Script()
// 	return UnlockExample{ScriptSig: ss}, err
// }

// // P2WPKH: witness = [<sig>, <pubkey>]
// func ExampleUnlockP2WPKH(sigDERPlusHashType, pubKey []byte) UnlockExample {
// 	return UnlockExample{Witness: wire.TxWitness{sigDERPlusHashType, pubKey}}
// }

// // P2WSH: witness = [<args...>, <witnessScript>] （最后一项是脚本原文）
// func ExampleUnlockP2WSH(args [][]byte, witnessScript []byte) UnlockExample {
// 	w := make(wire.TxWitness, 0, len(args)+1)
// 	for _, a := range args {
// 		w = append(w, a)
// 	}
// 	w = append(w, witnessScript)
// 	return UnlockExample{Witness: w}
// }

// // P2SH-P2WPKH: scriptSig=<redeemScript>; witness 与 P2WPKH 相同
// func ExampleUnlockP2SHP2WPKH(redeemScript, sigDERPlusHashType, pubKey []byte) (UnlockExample, error) {
// 	ss, err := txscript.NewScriptBuilder().AddData(redeemScript).Script()
// 	return UnlockExample{ScriptSig: ss, Witness: wire.TxWitness{sigDERPlusHashType, pubKey}}, err
// }

// // P2SH-P2WSH: scriptSig=<redeemScript>; witness 与 P2WSH 相同
// func ExampleUnlockP2SHP2WSH(redeemScript []byte, args [][]byte, witnessScript []byte) (UnlockExample, error) {
// 	ss, err := txscript.NewScriptBuilder().AddData(redeemScript).Script()
// 	w := make(wire.TxWitness, 0, len(args)+1)
// 	for _, a := range args {
// 		w = append(w, a)
// 	}
// 	w = append(w, witnessScript)
// 	return UnlockExample{ScriptSig: ss, Witness: w}, err
// }

// // P2TR-KeyPath: witness = [<schnorrSig>(+opt sighash), annex?]
// // 注意：sig 可为 64B（SIGHASH_DEFAULT）或 65B（附 1 字节类型）。
// func ExampleUnlockP2TRKeyPath(sig []byte, annexOptional []byte) UnlockExample {
// 	w := wire.TxWitness{sig}
// 	if len(annexOptional) > 0 {
// 		w = append(w, annexOptional) // 必须以 0x50 开头（调用方保证）
// 	}
// 	return UnlockExample{Witness: w}
// }

// // P2TR-ScriptPath: witness = [<args...>, annex?, <tapscript>, <control_block>]
// func ExampleUnlockP2TRScriptPath(args [][]byte, tapscript []byte, controlBlock []byte, annexOptional []byte) UnlockExample {
// 	w := make(wire.TxWitness, 0, len(args)+2)
// 	for _, a := range args {
// 		w = append(w, a)
// 	}
// 	if len(annexOptional) > 0 {
// 		w = append(w, annexOptional) // 必须 0x50 开头
// 	}
// 	w = append(w, tapscript)
// 	w = append(w, controlBlock)
// 	return UnlockExample{Witness: w}
// }

// // --- 常见 witnessScript 生成（P2WSH/Tapscript） ---------------------------

// // WitnessScript_2of3: 2-of-3 多签（SegWit v0 用于 P2WSH）。公钥建议按 BIP-67 排序（非共识）。
// func WitnessScript_2of3(pub1, pub2, pub3 []byte) ([]byte, error) {
// 	return txscript.NewScriptBuilder().
// 		AddOp(txscript.OP_2).
// 		AddData(pub1).AddData(pub2).AddData(pub3).
// 		AddOp(txscript.OP_3).
// 		AddOp(txscript.OP_CHECKMULTISIG).
// 		Script()
// }

// // WitnessScript_CLTV_Singlesig: 绝对时间锁 + 单签
// func WitnessScript_CLTV_Singlesig(locktime int64, pub []byte) ([]byte, error) {
// 	return txscript.NewScriptBuilder().
// 		AddInt64(locktime).AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).AddOp(txscript.OP_DROP).
// 		AddData(pub).AddOp(txscript.OP_CHECKSIG).
// 		Script()
// }

// // WitnessScript_CSV_Singlesig: 相对时间锁 + 单签（交易 version>=2，nSequence 满足 csvDelay）
// func WitnessScript_CSV_Singlesig(csvDelay int64, pub []byte) ([]byte, error) {
// 	return txscript.NewScriptBuilder().
// 		AddInt64(csvDelay).AddOp(txscript.OP_CHECKSEQUENCEVERIFY).AddOp(txscript.OP_DROP).
// 		AddData(pub).AddOp(txscript.OP_CHECKSIG).
// 		Script()
// }

// // WitnessScript_Hashlock_Singlesig (SHA256 版)：给出 R + 签名
// func WitnessScript_Hashlock_Singlesig(H32 []byte, pub []byte) ([]byte, error) {
// 	if len(H32) != 32 {
// 		return nil, errors.New("H must be 32 bytes (SHA256)")
// 	}
// 	return txscript.NewScriptBuilder().
// 		AddOp(txscript.OP_SHA256).AddData(H32).AddOp(txscript.OP_EQUALVERIFY).
// 		AddData(pub).AddOp(txscript.OP_CHECKSIG).
// 		Script()
// }

// // Tapscript_HTLC 简化版（Script-path）：成功(收款人+R) / 退款(发送人+CLTV)
// // 注意：这是 Tapscript 文本；用于 P2TR Script-path。
// func Tapscript_HTLC(H32 []byte, recvPubXOnly []byte, locktime int64, sendPubXOnly []byte) ([]byte, error) {
// 	if len(H32) != 32 {
// 		return nil, errors.New("H must be 32 bytes (SHA256)")
// 	}
// 	b := txscript.NewScriptBuilder()
// 	b.AddOp(txscript.OP_IF)
// 	// 成功路径：收款人签名 + R 校验
// 	b.AddData(recvPubXOnly).AddOp(txscript.OP_CHECKSIGVERIFY)
// 	b.AddOp(txscript.OP_SIZE).AddInt64(32).AddOp(txscript.OP_EQUALVERIFY)
// 	b.AddOp(txscript.OP_SHA256).AddData(H32).AddOp(txscript.OP_EQUAL)
// 	b.AddOp(txscript.OP_ELSE)
// 	// 超时退款
// 	b.AddInt64(locktime).AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).AddOp(txscript.OP_DROP)
// 	b.AddData(sendPubXOnly).AddOp(txscript.OP_CHECKSIG)
// 	b.AddOp(txscript.OP_ENDIF)
// 	return b.Script()
// }

// // --- PSBT 相关小助手（签名侧需要的前序输出） ----------------------------

// // PSBTInputForWitnessUTXO 生成 PSBT 里建议放入的 witness_utxo（供签名摘要用）。
// func PSBTInputForWitnessUTXO(valueSats int64, scriptPubKey []byte) *wire.TxOut {
// 	return &wire.TxOut{Value: valueSats, PkScript: scriptPubKey}
// }

// // --- 使用示例 --------------------------------------------------------------
// //
// // 1) 从地址生成锁定脚本：
// //    spk, kind, err := ScriptPubKeyFromAddress("bc1q...", &chaincfg.MainNetParams)
// //
// // 2) P2WPKH 花费 witness：
// //    ex := ExampleUnlockP2WPKH(sigDER1, pub1)
// //
// // 3) P2WSH 多签：
// //    wscript, _ := WitnessScript_2of3(pub1, pub2, pub3)
// //    spk, _ := BuildP2WSHScriptPubKey(wscript)
// //    ex := ExampleUnlockP2WSH([][]byte{[]byte(""), sig1, sig3}, wscript) // 注意 dummy ""
// //
// // 4) P2SH 包装 P2WPKH：
// //    redeem, _ := BuildRedeemScriptP2WPKH(pkh20)
// //    ss_wit, _ := ExampleUnlockP2SHP2WPKH(redeem, sigDER1, pub1)
// //
// // 5) Taproot：
// //    // Key-path：ex := ExampleUnlockP2TRKeyPath(sig64or65, nil)
// //    // Script-path：ex := ExampleUnlockP2TRScriptPath(args, tapscript, controlBlock, nil)

// // --- vsize/weight 估算与序列化尺寸工具 ------------------------------------

// // TxBaseSize 返回剥离 witness 后的序列化大小（B）。
// func TxBaseSize(tx *wire.MsgTx) int { return tx.SerializeSizeStripped() }

// // TxTotalSize 返回包含 witness 的完整序列化大小（B）。
// func TxTotalSize(tx *wire.MsgTx) int { return tx.SerializeSize() }

// // TxWeight 依据 BIP-141：weight = base*3 + total。
// func TxWeight(tx *wire.MsgTx) int {
// 	base := TxBaseSize(tx)
// 	total := TxTotalSize(tx)
// 	return base*3 + total
// }

// // TxVSize = ceil(weight/4)。
// func TxVSize(tx *wire.MsgTx) int {
// 	w := TxWeight(tx)
// 	return (w + 3) / 4
// }

// // WeightFromBaseAndWitness: 当你只知道 base_size 与 witness 字节总量时的估算。
// // 注意：包含 witness 的交易还需加上 2 字节的 marker+flag。
// func WeightFromBaseAndWitness(baseSize int, witnessBytes int, hasWitness bool) int {
// 	extra := witnessBytes
// 	if hasWitness {
// 		extra += 2
// 	}
// 	return baseSize*4 + extra
// }

// // VSizeFromWeight 按向上取整规则计算 vsize。
// func VSizeFromWeight(weight int) int { return (weight + 3) / 4 }

// // ---- 见证与脚本尺寸估算通用函数 -------------------------------------------

// // compactSizeLen 返回 Bitcoin CompactSize(varint) 的长度（B）。
// func compactSizeLen(n int) int {
// 	switch {
// 	case n < 0xfd:
// 		return 1
// 	case n <= 0xffff:
// 		return 3
// 	case n <= 0xffffffff:
// 		return 5
// 	default:
// 		return 9
// 	}
// }

// // WitnessPushesSize 计算 witness 栈按序列化后的总字节（含项个数 varint + 每项长度 varint）。
// func WitnessPushesSize(lengths ...int) int {
// 	size := compactSizeLen(len(lengths))
// 	for _, l := range lengths {
// 		size += compactSizeLen(l)
// 		size += l
// 	}
// 	return size
// }

// // RoughWitnessSizeP2WPKH 估算 P2WPKH witness 大小；sigLen 典型 71~73，pubkey=33。
// func RoughWitnessSizeP2WPKH(sigLen int) int {
// 	if sigLen == 0 {
// 		sigLen = 73
// 	}
// 	return WitnessPushesSize(sigLen, 33)
// }

// // RoughWitnessSizeP2WSH 估算 P2WSH witness 大小；argsLen 为按顺序的参数字节长度，wscriptLen 为脚本原文字节数。
// func RoughWitnessSizeP2WSH(argsLen []int, wscriptLen int) int {
// 	lens := append(append([]int{}, argsLen...), wscriptLen)
// 	return WitnessPushesSize(lens...)
// }

// // RoughWitnessSizeP2TRKeyPath 估算 P2TR Key-path witness 大小；sigLen=64|65；annexLen 可为 0。
// func RoughWitnessSizeP2TRKeyPath(sigLen, annexLen int) int {
// 	lens := []int{sigLen}
// 	if annexLen > 0 {
// 		lens = append(lens, annexLen)
// 	}
// 	return WitnessPushesSize(lens...)
// }

// // RoughWitnessSizeP2TRScriptPath 估算 P2TR Script-path witness 大小。
// func RoughWitnessSizeP2TRScriptPath(argLens []int, tapscriptLen, controlBlockLen, annexLen int) int {
// 	lens := append([]int{}, argLens...)
// 	if annexLen > 0 {
// 		lens = append(lens, annexLen)
// 	}
// 	lens = append(lens, tapscriptLen, controlBlockLen)
// 	return WitnessPushesSize(lens...)
// }

// // BaseInputSize 返回“无 witness 序列化”下单个输入基大小：outpoint(36) + scriptsig varint + scriptsig + sequence(4)。
// func BaseInputSize(scriptSigLen int) int { return 36 + compactSizeLen(scriptSigLen) + scriptSigLen + 4 }

// // TxOutSize: amount(8) + pkScript varint + pkScript。
// func TxOutSize(pkScriptLen int) int { return 8 + compactSizeLen(pkScriptLen) + pkScriptLen }

// // --- Tapscript 多签（CHECKSIGADD 阈值）构造器 ------------------------------

// // Tapscript_ThresholdCHECKSIGADD 生成 k-of-n 的阈值脚本：
// // OP_0; for each pub: <pub> OP_CHECKSIGADD; <k> OP_NUMEQUAL
// // 见证侧应为每个 pubkey 提供一个签名（或空字节表示 0），顺序与 pubkeys 对应。
// func Tapscript_ThresholdCHECKSIGADD(pubkeys [][]byte, k int64) ([]byte, error) {
// 	if len(pubkeys) == 0 {
// 		return nil, errors.New("no pubkeys")
// 	}
// 	if k <= 0 || int(k) > len(pubkeys) {
// 		return nil, fmt.Errorf("invalid k: %d", k)
// 	}
// 	b := txscript.NewScriptBuilder().AddOp(txscript.OP_0)
// 	for _, pk := range pubkeys {
// 		b.AddData(pk).AddOp(txscript.OP_CHECKSIGADD)
// 	}
// 	b.AddInt64(k).AddOp(txscript.OP_NUMEQUAL)
// 	return b.Script()
// }

// // ExampleUnlock_TapThresholdWitness 依据给定签名（或空字节）顺序，组织 Script-path 的 witness。
// // 注意：signatures 的长度必须与 pubkeys 数量一致；用空切片表示某个 pub 没有签名。
// func ExampleUnlock_TapThresholdWitness(signatures [][]byte, tapscript []byte, controlBlock []byte) UnlockExample {
// 	w := make(wire.TxWitness, 0, len(signatures)+2)
// 	for _, s := range signatures {
// 		w = append(w, s)
// 	}
// 	w = append(w, tapscript)
// 	w = append(w, controlBlock)
// 	return UnlockExample{Witness: w}
// }
