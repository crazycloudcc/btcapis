package tx

import (
	"fmt"

	"github.com/crazycloudcc/btcapis/src/decoders"
	"github.com/crazycloudcc/btcapis/src/types"
)

// DisasmScriptPubKey 反汇编交易输出脚本
func DisasmScriptPubKey(t *types.Tx, vout int) ([]types.ScriptOp, string, error) {
	if vout < 0 || vout >= len(t.TxOut) {
		return nil, "", fmt.Errorf("vout index out of range")
	}
	return decoders.DisasmScript(t.TxOut[vout].PkScript)
}
