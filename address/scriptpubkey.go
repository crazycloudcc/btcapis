// Package address 脚本公钥处理
package address

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/crazycloudcc/btcapis/types"
)

// ScriptPubKeyFromAddress 透传到 types 包中的实现，便于上层按地址生成锁定脚本。
func ScriptPubKeyFromAddress(addr string, params *chaincfg.Params) ([]byte, types.AddressType, error) {
	return types.ScriptPubKeyFromAddress(addr, params)
}
