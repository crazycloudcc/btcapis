// 简单 UTXO 选择 + 收敛
package txbuilder

// import "github.com/crazycloudcc/btcapis/types"

// // SelectInputsGreedy - 最小骨架：直接使用请求里给定的 Inputs；返回所选及总额
// func SelectInputsGreedy(req *types.BuildTxRequest) (selected []types.UTXO, sum int64) {
// 	selected = make([]types.UTXO, len(req.Inputs))
// 	copy(selected, req.Inputs)
// 	for _, in := range selected {
// 		sum += in.Value // 若为 0，后续可由上层补齐
// 	}
// 	return
// }
