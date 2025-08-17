// 提供一些通用的解析函数
package btcdapis

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"
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

// 将 opcode 映射为可读名称
func parseOpCode(op byte) string {
	switch {
	case op >= 0x01 && op <= 0x4b:
		// 1..75 字节的直接推送
		return fmt.Sprintf("OP_DATA_%d", int(op))
	}
	switch op {
	case txscript.OP_0:
		return "OP_0" // 将空字节数组推入栈中
	case txscript.OP_PUSHDATA1:
		return "OP_PUSHDATA1" // 下一个字节包含要推送的数据长度（1字节）
	case txscript.OP_PUSHDATA2:
		return "OP_PUSHDATA2" // 下两个字节包含要推送的数据长度（2字节，小端序）
	case txscript.OP_PUSHDATA4:
		return "OP_PUSHDATA4" // 下四个字节包含要推送的数据长度（4字节，小端序）
	case txscript.OP_1:
		return "OP_1" // 将数字1推入栈中
	case txscript.OP_2:
		return "OP_2" // 将数字2推入栈中
	case txscript.OP_3:
		return "OP_3" // 将数字3推入栈中
	case txscript.OP_4:
		return "OP_4" // 将数字4推入栈中
	case txscript.OP_5:
		return "OP_5" // 将数字5推入栈中
	case txscript.OP_6:
		return "OP_6" // 将数字6推入栈中
	case txscript.OP_7:
		return "OP_7" // 将数字7推入栈中
	case txscript.OP_8:
		return "OP_8" // 将数字8推入栈中
	case txscript.OP_9:
		return "OP_9" // 将数字9推入栈中
	case txscript.OP_10:
		return "OP_10" // 将数字10推入栈中
	case txscript.OP_11:
		return "OP_11" // 将数字11推入栈中
	case txscript.OP_12:
		return "OP_12" // 将数字12推入栈中
	case txscript.OP_13:
		return "OP_13" // 将数字13推入栈中
	case txscript.OP_14:
		return "OP_14" // 将数字14推入栈中
	case txscript.OP_15:
		return "OP_15" // 将数字15推入栈中
	case txscript.OP_16:
		return "OP_16" // 将数字16推入栈中

	case txscript.OP_DUP:
		return "OP_DUP" // 复制栈顶元素
	case txscript.OP_HASH160:
		return "OP_HASH160" // 对栈顶元素进行RIPEMD160(SHA256())哈希运算
	case txscript.OP_EQUAL:
		return "OP_EQUAL" // 比较栈顶两个元素是否相等，结果推入栈中
	case txscript.OP_EQUALVERIFY:
		return "OP_EQUALVERIFY" // 比较栈顶两个元素是否相等，不相等则脚本失败
	case txscript.OP_CHECKSIG:
		return "OP_CHECKSIG" // 验证签名，使用栈顶的公钥和签名验证消息
	case txscript.OP_CHECKMULTISIG:
		return "OP_CHECKMULTISIG" // 验证多重签名，支持M-of-N签名验证
	case txscript.OP_CHECKLOCKTIMEVERIFY:
		return "OP_CHECKLOCKTIMEVERIFY" // 检查锁定时间，确保交易在指定时间后才能被确认
	case txscript.OP_CHECKSEQUENCEVERIFY:
		return "OP_CHECKSEQUENCEVERIFY" // 检查序列号，用于相对时间锁定
	case txscript.OP_RETURN:
		return "OP_RETURN" // 标记交易输出为不可花费，用于存储数据
	case txscript.OP_DROP:
		return "OP_DROP" // 移除栈顶元素
	case txscript.OP_SWAP:
		return "OP_SWAP" // 交换栈顶两个元素的位置
	case txscript.OP_IF:
		return "OP_IF" // 条件执行开始，如果栈顶元素为真则执行后续代码
	case txscript.OP_NOTIF:
		return "OP_NOTIF" // 条件执行开始，如果栈顶元素为假则执行后续代码
	case txscript.OP_ELSE:
		return "OP_ELSE" // 条件执行的分支，当IF条件为假时执行
	case txscript.OP_ENDIF:
		return "OP_ENDIF" // 结束条件执行块
	case txscript.OP_VERIFY:
		return "OP_VERIFY" // 验证栈顶元素，如果为假则脚本失败

	// Taproot/Tapscript 常用
	case txscript.OP_CHECKSIGADD:
		return "OP_CHECKSIGADD" // Tapscript操作码，验证签名并累加计数器，用于Schnorr多重签名
	}

	// 兜底
	return fmt.Sprintf("OP_%d", int(op))
}
