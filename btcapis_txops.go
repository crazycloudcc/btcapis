package btcapis

import (
	"github.com/crazycloudcc/btcapis/tx"
	"github.com/crazycloudcc/btcapis/types"
)

// 对外门面：解析某个输入的脚本/控制块为 OP code 列表
func (TxModule) AnalyzeInput(t *types.Tx, idx int) (*types.TapscriptInfo, error) {
	return tx.AnalyzeInput(t, idx)
}

func (TxModule) DisasmScriptPubKey(t *types.Tx, vout int) ([]types.ScriptOp, string, error) {
	return tx.DisasmScriptPubKey(t, vout)
}
