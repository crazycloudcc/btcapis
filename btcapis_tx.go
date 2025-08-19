package btcapis

import (
	"github.com/crazycloudcc/btcapis/tx"
	"github.com/crazycloudcc/btcapis/types"
)

// // 1) 创建未签名交易（返回：unsigned tx hex + PSBT base64）
// func (TxModule) CreateTransaction(ctx context.Context, c *Client, req *types.BuildTxRequest) (*types.BuildTxResult, error) {
// 	// 费率优先从 req.FeeRateSatPerVb；若<=0 则尝试 c.router.EstimateFeeRate(ctx, 3)，失败给个兜底比如 5.0
// 	return txbuilder.Build(ctx, c.router, req)
// }

// // 2) 广播已签名交易（参数是 signed tx hex）
// func (TxModule) BroadcastSigned(ctx context.Context, c *Client, signedHex string) (string, error) {
// 	raw, err := hex.DecodeString(strings.TrimSpace(signedHex))
// 	if err != nil {
// 		return "", err
// 	}
// 	// 使用 router.Broadcast —— 内部先尝试 bitcoind，失败再试 mempool.space
// 	return c.router.Broadcast(ctx, raw)
// }

// 对外门面：使用输入数据在数组中的索引 解析某个输入的脚本/控制块为 OP code 列表
func (TxModule) AnalyzeTxInWithIdx(t *types.Tx, idx int) (*types.TapscriptInfo, error) {
	return tx.AnalyzeTxInWithIdx(t, idx)
}

// 对外门面：使用输入数据 解析某个输入的脚本/控制块为 OP code 列表
func (TxModule) AnalyzeTxIn(t *types.Tx, in *types.TxIn) (*types.TapscriptInfo, error) {
	return tx.AnalyzeTxIn(t, in)
}

// 对外门面：使用交易输出索引 反汇编某个交易的特定输出脚本
func (TxModule) DisasmScriptPubKey(t *types.Tx, vout int) ([]types.ScriptOp, string, error) {
	return tx.DisasmScriptPubKey(t, vout)
}

// 对外门面：扫描交易中的所有输入，提取BRC-20动作
func (TxModule) ExtractBRC20(t *types.Tx) []types.BRC20Action {
	return tx.ExtractBRC20(t)
}

// 对外门面：扫描交易中的所有输入，提取Runestone动作
func (TxModule) ExtractRunes(t *types.Tx) []types.Runestone {
	return tx.ExtractRunes(t)
}
