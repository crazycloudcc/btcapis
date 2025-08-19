package tx

import (
	"fmt"

	"github.com/crazycloudcc/btcapis/script"
	"github.com/crazycloudcc/btcapis/types"
)

// DisasmScriptPubKey 反汇编交易输出脚本
func DisasmScriptPubKey(t *types.Tx, vout int) ([]types.ScriptOp, string, error) {
	if vout < 0 || vout >= len(t.TxOut) {
		return nil, "", fmt.Errorf("vout index out of range")
	}
	return script.DisasmScript(t.TxOut[vout].PkScript)
}
