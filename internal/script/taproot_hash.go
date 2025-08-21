package script

import (
	"crypto/sha256"
)

// taggedHash implements BIP-340 tagged hashing.
func taggedHash(tag string, data ...[]byte) [32]byte {
	var out [32]byte
	th := sha256.Sum256([]byte(tag))
	buf := make([]byte, 0, 64)
	buf = append(buf, th[:]...)
	buf = append(buf, th[:]...)
	for _, d := range data {
		buf = append(buf, d...)
	}
	out = sha256.Sum256(buf)
	return out
}

// encodeVarInt encodes a Bitcoin varint.
func encodeVarInt(v uint64) []byte {
	if v < 0xfd {
		return []byte{byte(v)}
	}
	if v <= 0xffff {
		return []byte{0xfd, byte(v), byte(v >> 8)}
	}
	if v <= 0xffffffff {
		return []byte{0xfe, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}
	}
	return []byte{0xff, byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24), byte(v >> 32), byte(v >> 40), byte(v >> 48), byte(v >> 56)}
}

// TapLeafHash computes the tagged hash for a tapscript leaf per BIP-342.
// H_TapLeaf(leafVer || varint(len(script)) || script)
func TapLeafHash(leafVersion byte, script []byte) [32]byte {
	data := []byte{leafVersion}
	data = append(data, encodeVarInt(uint64(len(script)))...)
	data = append(data, script...)
	return taggedHash("TapLeaf", data)
}
