// 提供一些通用的解析函数
package script

import (
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

// 判定前序锁定脚本是否为 Taproot P2TR
func isTaprootPkScript(pkScript []byte) bool {
	ver, prog, err := txscript.ExtractWitnessProgramInfo(pkScript)
	if err != nil {
		return false
	}
	return ver == 1 && len(prog) == 32
}
