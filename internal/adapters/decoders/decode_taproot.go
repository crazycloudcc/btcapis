package decoders

import (
	"encoding/hex"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/types"
)

// 严格：必须能成功解析 control block（长度规则 + 头字节 leaf version 合法）
func IsTapScriptPathWitness(w [][]byte) bool {
	if len(w) < 2 {
		return false
	}
	cb := w[len(w)-1]
	_, err := ParseControlBlock(cb)
	return err == nil
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

// TapLeafHash computes the tagged hash for a tapscript leaf per BIP-342.
// H_TapLeaf(leafVer || varint(len(script)) || script)
func TapLeafHash(leafVersion byte, script []byte) [32]byte {
	data := []byte{leafVersion}
	data = append(data, encodeVarInt(uint64(len(script)))...)
	data = append(data, script...)
	return taggedHash("TapLeaf", data)
}
