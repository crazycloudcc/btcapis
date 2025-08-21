// 提供一些通用的解析函数
package decoders

import (
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/txscript"
)

// DER ECDSA 签名（可能带末尾 1B sighash）解析
func parseDERSignatureWithSigHash(b []byte) (rHex, sHex string, sighash byte, ok bool) {
	if len(b) < 8 || b[0] != 0x30 {
		return "", "", 0, false
	}
	// 优先当作 [DER || sighash] 处理
	var sig struct{ R, S *big.Int }
	der := b[:len(b)-1]
	if _, err := asn1.Unmarshal(der, &sig); err == nil && sig.R != nil && sig.S != nil {
		return sig.R.Text(16), sig.S.Text(16), b[len(b)-1], true
	}
	// 兜底：没有 sighash 的情况（极少见）
	if _, err := asn1.Unmarshal(b, &sig); err == nil && sig.R != nil && sig.S != nil {
		return sig.R.Text(16), sig.S.Text(16), 0xff, true // 0xff 表示未知/无
	}
	return "", "", 0, false
}

// Taproot Schnorr 签名解析（64B [r||s] 或 65B [r||s||sighash]）
func parseSchnorrSignature(b []byte) (rHex, sHex string, sighash byte, ok bool) {
	if len(b) != 64 && len(b) != 65 {
		return "", "", 0, false
	}
	r := b[:32]
	s := b[32:64]
	var sh byte = 0x00 // DEFAULT
	if len(b) == 65 {
		sh = b[64]
	}
	return hex.EncodeToString(r), hex.EncodeToString(s), sh, true
}

// 解析公钥并打印坐标（压缩/非压缩）
func parsePubKeyCoords(b []byte) (compressed bool, xHex, yHex string, ok bool) {
	pk, err := btcec.ParsePubKey(b)
	if err != nil {
		return false, "", "", false
	}
	compressed = (len(b) == 33)
	xHex = pk.X().Text(16)
	yHex = pk.Y().Text(16)
	return compressed, xHex, yHex, true
}

// SIGHASH 名称（Taproot 支持 DEFAULT=0x00）
func parseSigHash(sh byte, taproot bool) string {
	var base string
	t := sh & 0x03
	if taproot && sh == 0x00 {
		base = "DEFAULT"
	} else {
		switch t {
		case 0x01:
			base = "ALL"
		case 0x02:
			base = "NONE"
		case 0x03:
			base = "SINGLE"
		default:
			if taproot {
				base = "RESERVED/UNKNOWN"
			} else {
				base = "UNKNOWN"
			}
		}
	}
	if (sh & 0x80) != 0 {
		base += "|ANYONECANPAY"
	}
	return base
}

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

// 判定前序锁定脚本是否为 Taproot P2TR
func isTaprootPkScript(pkScript []byte) bool {
	ver, prog, err := txscript.ExtractWitnessProgramInfo(pkScript)
	if err != nil {
		return false
	}
	return ver == 1 && len(prog) == 32
}
