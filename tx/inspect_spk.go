package tx

import (
	"fmt"

	"github.com/crazycloudcc/btcapis/script"
	"github.com/crazycloudcc/btcapis/types"
)

func DisasmScriptPubKey(t *types.Tx, vout int) ([]types.ScriptOp, string, error) {
	if vout < 0 || vout >= len(t.Vout) {
		return nil, "", fmt.Errorf("vout index out of range")
	}
	return script.DisasmScript(t.Vout[vout].ScriptPubKey)
}
