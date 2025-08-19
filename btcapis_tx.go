package btcapis

import (
	"github.com/crazycloudcc/btcapis/tx"
	"github.com/crazycloudcc/btcapis/types"
)

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
