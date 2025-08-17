package script

import (
	"encoding/hex"
	"fmt"

	"github.com/crazycloudcc/btcapis/types"
)

// IsTapScriptPathWitness: 简单判断 witness 最后一项是否是 control block
func IsTapScriptPathWitness(w [][]byte) bool {
	if len(w) < 2 {
		return false
	}
	cb := w[len(w)-1]
	if len(cb) < 33 || (len(cb)-33)%32 != 0 {
		return false
	}
	// header: arbitrary; leaf version must be even; we don't strictly validate here
	return true
}

func ParseControlBlock(cb []byte) (types.TapControlBlock, error) {
	if len(cb) < 33 || (len(cb)-33)%32 != 0 {
		return types.TapControlBlock{}, fmt.Errorf("invalid control block length=%d", len(cb))
	}
	header := cb[0]
	leafVer := header & 0xfe
	parity := int(header>>7) & 1
	intKey := hex.EncodeToString(cb[1:33])
	var branches []string
	for i := 33; i < len(cb); i += 32 {
		branches = append(branches, hex.EncodeToString(cb[i:i+32]))
	}
	return types.TapControlBlock{
		Header: header, LeafVersion: leafVer, Parity: parity,
		InternalKey: intKey, MerkleHashes: branches,
	}, nil
}

// ExtractTapScriptPath: 返回(栈元素, 脚本, 控制块)
func ExtractTapScriptPath(w [][]byte) (stack [][]byte, script []byte, control []byte, ok bool) {
	if !IsTapScriptPathWitness(w) {
		return nil, nil, nil, false
	}
	cb := w[len(w)-1]
	script = w[len(w)-2]
	stack = w[:len(w)-2]
	return stack, script, cb, true
}
