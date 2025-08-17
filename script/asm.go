// script/asm.go
package script

import (
	"github.com/btcsuite/btcd/txscript"
)

// Asm 返回类似 "OP_0 <20-bytes> ..." 的汇编字符串
/* 外部使用举例:
asm, _ := script.Asm(tx.Vout[0].ScriptPubKey)
fmt.Println(asm)
*/
func Asm(pkScript []byte) (string, error) {
	return txscript.DisasmString(pkScript)
}
