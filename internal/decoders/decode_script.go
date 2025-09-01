package decoders

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/crazycloudcc/btcapis/types"
)

func PKScriptToType(pkScript []byte) types.AddressType {
	n := len(pkScript)
	if n == 0 {
		return types.AddrUnknown
	}

	// P2PKH: OP_DUP OP_HASH160 PUSH20 <20> OP_EQUALVERIFY OP_CHECKSIG
	if n == 25 && pkScript[0] == 0x76 && pkScript[1] == 0xa9 && pkScript[2] == 0x14 && pkScript[23] == 0x88 && pkScript[24] == 0xac {
		return types.AddrP2PKH
	}
	// P2SH: OP_HASH160 PUSH20 <20> OP_EQUAL
	if n == 23 && pkScript[0] == 0xa9 && pkScript[1] == 0x14 && pkScript[22] == 0x87 {
		return types.AddrP2SH
	}
	// P2WPKH v0: OP_0 PUSH20 <20>
	if n == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14 {
		return types.AddrP2WPKH
	}
	// P2WSH v0: OP_0 PUSH32 <32>
	if n == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20 {
		return types.AddrP2WSH
	}
	// P2TR v1: OP_1 PUSH32 <32>
	if n == 34 && pkScript[0] == 0x51 && pkScript[1] == 0x20 {
		return types.AddrP2TR
	}
	return types.AddrUnknown
}

func DecodePkScript(pkScript []byte) (*types.AddressInfo, error) {
	cls, addrs, reqSigs, err := txscript.ExtractPkScriptAddrs(pkScript, types.CurrentNetworkParams)
	if err != nil {
		return nil, err
	}
	fmt.Printf("DecodePkScript cls: %s\n", cls.String())
	out := &types.AddressInfo{PKScript: pkScript, Typ: PKScriptToType(pkScript), Cls: cls, ReqSigs: reqSigs}
	out.Addresses = make([]string, len(addrs))
	for i, a := range addrs {
		out.Addresses[i] = a.EncodeAddress()
	}

	printDecodePkScript(out)
	return out, nil
}

func printDecodePkScript(info *types.AddressInfo) {

	ops, asm, err := DecodeAsmScript(info.PKScript)
	if err != nil {
		fmt.Printf("[DecodeAsmScript] %s\n", err)
	}

	// 打印详细的解析结果，便于调试和验证
	fmt.Printf("PKScript2Address ===================================\n")
	fmt.Printf("[Network] %s\n", types.CurrentNetwork)
	fmt.Printf("[PKScript] %x\n", info.PKScript)
	fmt.Printf("[AsmScriptOps] %v\n", ops)
	fmt.Printf("[AsmScript] %s\n", asm)
	fmt.Printf("[AddressType] %s\n", info.Typ)
	fmt.Printf("[ReqSigs] %d\n", info.ReqSigs)
	fmt.Printf("[Addresses] %v\n", info.Addresses)
	fmt.Printf("PKScript2Address ===================================\n")
}

// // Script2Addr 仅做最小化识别：p2pkh/p2sh/p2wpkh/p2wsh/p2tr，其余返回 "unknown"。
// // 需要生成地址时，再扩展为传入 network，并用 bech32/base58 生成。
// func Script2Addr(pkScript []byte, network types.Network) (typ types.AddressType, addrs []string) {
// 	b := pkScript
// 	n := len(b)
// 	params := network.ToParams()

// 	// P2PKH: OP_DUP OP_HASH160 PUSH20 <20> OP_EQUALVERIFY OP_CHECKSIG
// 	if n == 25 && b[0] == 0x76 && b[1] == 0xa9 && b[2] == 0x14 && b[23] == 0x88 && b[24] == 0xac {
// 		if params == nil {
// 			return types.AddrP2PKH, nil
// 		}
// 		pkh20 := b[3:23]
// 		if a, err := btcutil.NewAddressPubKeyHash(pkh20, params); err == nil {
// 			return types.AddrP2PKH, []string{a.EncodeAddress()}
// 		}
// 		return types.AddrP2PKH, nil
// 	}
// 	// P2SH: OP_HASH160 PUSH20 <20> OP_EQUAL
// 	if n == 23 && b[0] == 0xa9 && b[1] == 0x14 && b[22] == 0x87 {
// 		if params == nil {
// 			return types.AddrP2SH, nil
// 		}
// 		sh20 := b[2:22]
// 		if a, err := btcutil.NewAddressScriptHashFromHash(sh20, params); err == nil {
// 			return types.AddrP2SH, []string{a.EncodeAddress()}
// 		}
// 		return types.AddrP2SH, nil
// 	}
// 	// P2WPKH v0: OP_0 PUSH20 <20>
// 	if n == 22 && b[0] == 0x00 && b[1] == 0x14 {
// 		if params == nil {
// 			return types.AddrP2WPKH, nil
// 		}
// 		wp20 := b[2:22]
// 		if a, err := btcutil.NewAddressWitnessPubKeyHash(wp20, params); err == nil {
// 			return types.AddrP2WPKH, []string{a.EncodeAddress()}
// 		}
// 		return types.AddrP2WPKH, nil
// 	}
// 	// P2WSH v0: OP_0 PUSH32 <32>
// 	if n == 34 && b[0] == 0x00 && b[1] == 0x20 {
// 		if params == nil {
// 			return types.AddrP2WSH, nil
// 		}
// 		wp32 := b[2:34]
// 		if a, err := btcutil.NewAddressWitnessScriptHash(wp32, params); err == nil {
// 			return types.AddrP2WSH, []string{a.EncodeAddress()}
// 		}
// 		return types.AddrP2WSH, nil
// 	}
// 	// P2TR v1: OP_1 PUSH32 <32>
// 	if n == 34 && b[0] == 0x51 && b[1] == 0x20 {
// 		if params == nil {
// 			return types.AddrP2TR, nil
// 		}
// 		xonly32 := b[2:34]
// 		if a, err := btcutil.NewAddressTaproot(xonly32, params); err == nil {
// 			return types.AddrP2TR, []string{a.EncodeAddress()}
// 		}
// 		return types.AddrP2TR, nil
// 	}
// 	return types.AddrUnknown, nil
// }
