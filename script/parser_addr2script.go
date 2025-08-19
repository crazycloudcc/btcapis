// 通过btcd库, 解析钱包地址, 获取钱包地址对应的类型, 锁定脚本, 脚本哈希等信息.
package script

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/crazycloudcc/btcapis/types"
)

// Addr2Script 函数：解析比特币地址并返回脚本信息
// 参数：
//   - addr: 要解析的比特币地址字符串
//   - params: 区块链网络参数（主网、测试网等）
//
// 返回值：
//   - *ScriptInfo: 解析后的脚本信息结构体
//   - error: 解析过程中的错误信息
func Addr2Script(addr string, params *chaincfg.Params) (*types.AddressScriptInfo, error) {

	// 使用 btcutil.DecodeAddress 自动识别地址格式
	// 支持：P2PKH(1...)、P2SH(3...)、P2WPKH(bc1q...)、P2WSH(bc1q...32B)、P2TR(bc1p...)
	a, err := btcutil.DecodeAddress(addr, params)
	if err != nil {
		return nil, err
	}

	// 创建返回结果结构体
	ret := &types.AddressScriptInfo{}

	// 使用类型断言判断具体的地址类型，并提取相应信息
	switch v := a.(type) {
	case *btcutil.AddressPubKeyHash: // 传统地址格式：1开头的地址
		ret.ScriptType = "P2PKH"           // 支付到公钥哈希
		ret.PubKeyHash = v.ScriptAddress() // 获取20字节的公钥哈希
		ret.WitnessVersion = -1            // -1表示非SegWit地址
	case *btcutil.AddressScriptHash: // 脚本哈希地址：3开头的地址（包含嵌套SegWit）
		ret.ScriptType = "P2SH"                  // 支付到脚本哈希
		ret.RedeemScriptHash = v.ScriptAddress() // 获取20字节的赎回脚本哈希
		ret.WitnessVersion = -1                  // -1表示非SegWit地址
	case *btcutil.AddressWitnessPubKeyHash: // 原生SegWit地址：bc1q开头的地址
		ret.ScriptType = "P2WPKH"                 // 支付到见证公钥哈希
		ret.WitnessScriptHash = v.ScriptAddress() // 获取20字节的见证脚本哈希
		ret.WitnessVersion = 0                    // 见证版本v0
		ret.WitnessLen = 20                       // 20字节脚本哈希
	case *btcutil.AddressWitnessScriptHash: // 原生SegWit脚本地址：bc1q开头的32字节地址
		ret.ScriptType = "P2WSH"                  // 支付到见证脚本哈希
		ret.WitnessScriptHash = v.ScriptAddress() // 获取32字节的见证脚本哈希
		ret.WitnessVersion = 0                    // 见证版本v0
		ret.WitnessLen = 32                       // 32字节脚本哈希
	case *btcutil.AddressTaproot: // Taproot地址：bc1p开头的地址
		ret.ScriptType = "P2TR"            // 支付到Taproot
		ret.WitnessVersion = 1             // 见证版本v1
		ret.WitnessLen = 32                // 32字节脚本哈希
		ret.TaprootKey = v.ScriptAddress() // 获取32字节的Taproot调整后公钥
	default:
		// 如果遇到未处理的地址类型，返回错误
		return nil, fmt.Errorf("unhandled address type %T", a)
	}

	printResult(addr, ret, params)
	return ret, nil
}

func printResult(addr string, info *types.AddressScriptInfo, network *chaincfg.Params) {

	// txscript.PayToAddrScript 根据地址类型生成相应的锁定脚本
	decodeAddre, _ := btcutil.DecodeAddress(addr, network)
	pkScript, _ := txscript.PayToAddrScript(decodeAddre)

	// 打印详细的解析结果，便于调试和验证
	fmt.Printf("Addr2ScriptHash ===================================\n")
	fmt.Printf("[Address] %s\n", addr)
	fmt.Printf("[AddressType] %s\n", info.ScriptType)
	fmt.Printf("[ScriptPubKey %d] %x\n", len(pkScript), pkScript)
	fmt.Printf("[PubKeyHash %d] %x\n", len(info.PubKeyHash), info.PubKeyHash)
	fmt.Printf("[RedeemScript %d] %x\n", len(info.RedeemScriptHash), info.RedeemScriptHash)
	fmt.Printf("[WitnessVersion] %x\n", info.WitnessVersion)
	fmt.Printf("[WitnessScript %d] %x\n", len(info.WitnessScriptHash), info.WitnessScriptHash)
	fmt.Printf("[TaprootKey %d] %x\n", len(info.TaprootKey), info.TaprootKey)
	fmt.Printf("Addr2ScriptHash ===================================\n")
}
