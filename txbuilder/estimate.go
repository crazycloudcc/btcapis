// vsize/fee 估算，脚本类型识别
package txbuilder

// import (
// 	"github.com/crazycloudcc/btcapis/script"
// )

// // ClassifyScript - 简包装（脚本类型识别）
// func ClassifyScript(spk []byte) string {
// 	typ, _ := script.Classify(spk)
// 	return typ
// }

// // EstimateInputVSize - 经验估算（足够做费控）
// func EstimateInputVSize(scriptType string) int64 {
// 	switch scriptType {
// 	case "p2wpkh", "v0_p2wpkh":
// 		return 68
// 	case "p2sh-p2wpkh":
// 		return 91
// 	case "p2tr", "v1_p2tr", "tr":
// 		return 57
// 	case "p2wsh", "v0_p2wsh":
// 		return 110
// 	case "p2pkh":
// 		return 148 // 非 segwit（仅估算用）
// 	case "p2sh":
// 		return 108
// 	default:
// 		return 110
// 	}
// }

// // EstimateOutputVSize - 经验估算
// func EstimateOutputVSize(scriptType string) int64 {
// 	switch scriptType {
// 	case "p2wpkh", "v0_p2wpkh":
// 		return 31
// 	case "p2tr", "v1_p2tr", "tr":
// 		return 43
// 	case "p2sh":
// 		return 32
// 	case "p2pkh":
// 		return 34
// 	default:
// 		return 34
// 	}
// }
