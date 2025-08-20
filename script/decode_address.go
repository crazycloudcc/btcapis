// 通过btcd库, 解析钱包地址, 获取钱包地址对应的类型, 锁定脚本, 脚本哈希等信息.
package script

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/crazycloudcc/btcapis/types"
)

func DecodeAddress(addr string) (*types.AddressScriptInfo, error) {
	decodeAddr, err := btcutil.DecodeAddress(addr, types.CurrentNetworkParams)
	if err != nil {
		return nil, fmt.Errorf("decode address: %w", err)
	}
	if !decodeAddr.IsForNet(types.CurrentNetworkParams) {
		return nil, fmt.Errorf("address not for network %s", types.CurrentNetwork)
	}

	// 1) 生成 scriptPubKey
	pkScript, err := txscript.PayToAddrScript(decodeAddr)
	if err != nil {
		return nil, fmt.Errorf("make scriptPubKey: %w", err)
	}

	// 2) 分类（模板识别）
	class, _, _, _ := txscript.ExtractPkScriptAddrs(pkScript, types.CurrentNetworkParams)
	stype := class.String()

	// 3) 反汇编
	asm, _ := txscript.DisasmString(pkScript)

	info := &types.AddressScriptInfo{
		Address:         addr,
		Typ:             types.AddressType(stype),
		ScriptPubKeyHex: pkScript,
		ScriptAsm:       asm,
	}

	// 4) 通用：从地址提取底层 “脚本参数”（hash160 / program）
	scriptParam := decodeAddr.ScriptAddress() // P2PKH/P2SH: 20B 哈希；SegWit：witness program

	switch class {
	case txscript.PubKeyHashTy:
		info.PubKeyHashHex = scriptParam // 20B
	case txscript.ScriptHashTy:
		info.RedeemScriptHashHex = scriptParam // 20B
	}

	// 5) SegWit / Taproot 解析 witness version & program
	// 规范形式：
	//   v0:   [OP_0  <20|32>    program]
	//   v1+:  [OP_(1..16) <2..40> program] （Taproot: v=1, 32B）
	if ver, prog, ok := extractWitness(pkScript); ok {
		info.IsWitness = true
		info.WitnessVersion = ver
		info.WitnessProgramHex = prog
		info.WitnessProgramLen = len(prog)

		enc := "bech32"
		if ver >= 1 {
			enc = "bech32m" // BIP-350
		}
		info.BechEncoding = enc

		// Taproot（v=1, 32B）→ 输出 key（x-only）
		if ver == 1 && len(prog) == 32 {
			info.Typ = types.AddrP2TR
			info.TaprootOutputKeyHex = prog
		}
	}

	printDecodeAddress(addr, info)
	return info, nil
}

// 提取 witness version & program；只做规范脚本的快速解析。
func extractWitness(pk []byte) (version int, program []byte, ok bool) {
	if len(pk) < 4 {
		return
	}
	// v0: OP_0, v1..16: OP_1..OP_16
	switch pk[0] {
	case txscript.OP_0:
		// 期望：OP_0 <pushlen> <program>
		if len(pk) < 2 {
			return
		}
		push := int(pk[1])
		if push < 2 || push > 40 || 2+push != len(pk) {
			return
		}
		return 0, pk[2:], true
	case txscript.OP_1, txscript.OP_2, txscript.OP_3, txscript.OP_4, txscript.OP_5,
		txscript.OP_6, txscript.OP_7, txscript.OP_8, txscript.OP_9, txscript.OP_10,
		txscript.OP_11, txscript.OP_12, txscript.OP_13, txscript.OP_14, txscript.OP_15, txscript.OP_16:
		if len(pk) < 2 {
			return
		}
		push := int(pk[1])
		if push < 2 || push > 40 || 2+push != len(pk) {
			return
		}
		return int(pk[0]-txscript.OP_1) + 1, pk[2:], true
	default:
		return
	}
}

func printDecodeAddress(addr string, info *types.AddressScriptInfo) {

	// txscript.PayToAddrScript 根据地址类型生成相应的锁定脚本
	decodeAddre, _ := btcutil.DecodeAddress(addr, types.CurrentNetworkParams)
	pkScript, _ := txscript.PayToAddrScript(decodeAddre)

	ops, asm, err := DisasmScript(pkScript)
	if err != nil {
		fmt.Printf("[DisasmScript] %s\n", err)
	}

	// 打印详细的解析结果，便于调试和验证
	fmt.Printf("Addr2ScriptHash ===================================\n")
	fmt.Printf("[pkScript %d] %x\n", len(pkScript), pkScript)
	fmt.Printf("[DisasmScriptOps] %v\n", ops)
	fmt.Printf("[DisasmScript] %s\n", asm)

	fmt.Printf("[Network] %s\n", types.CurrentNetwork)
	fmt.Printf("[Address] %s\n", addr)
	fmt.Printf("[AddressType] %s\n", info.Typ)

	fmt.Printf("[PubKeyHash %d] %x\n", len(info.PubKeyHashHex), info.PubKeyHashHex)
	fmt.Printf("[RedeemScript %d] %x\n", len(info.RedeemScriptHashHex), info.RedeemScriptHashHex)

	fmt.Printf("[ScriptPubKeyHex(pkScript) %d] %x\n", len(info.ScriptPubKeyHex), info.ScriptPubKeyHex)
	fmt.Printf("[ScriptAsm] %s\n", info.ScriptAsm)

	fmt.Printf("[IsWitness] %t\n", info.IsWitness)
	fmt.Printf("[WitnessVersion] %d\n", info.WitnessVersion)
	fmt.Printf("[WitnessScript %d] %x\n", len(info.WitnessProgramHex), info.WitnessProgramHex)
	fmt.Printf("[WitnessProgramLen] %d\n", info.WitnessProgramLen)
	fmt.Printf("[BechEncoding] %s\n", info.BechEncoding)

	fmt.Printf("[TaprootKey %d] %x\n", len(info.TaprootOutputKeyHex), info.TaprootOutputKeyHex)
	fmt.Printf("Addr2ScriptHash ===================================\n")
}
