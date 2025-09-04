package types

// 根据地址类型速查交易输入的虚拟字节大小（vsize）
// 返回值：估算值上限.
func GetInSize(addrType AddressType) int {
	switch addrType {
	case AddrP2PKH: // sig + pubkey | (sig 70–72 + pubkey 33；典型 ≈148。)
		return 148
	case AddrP2SH: // sig(s) + redeemScript | (依赖 m-of-n 可达几百字节, 不同脚本差别极大。2-of-3 ≈ 297。)
		return 297 // 仅估算 2-of-3 多签
	// case P2SH-P2WPKH: // redeemScript (22字节) + witness(sig+pub) | (base 部分固定：32+4+1+23+4=64；witness ≈107；总权重 ≈ 272 → vsize ≈ 68。)
	// 	return 依赖脚本, 依赖脚本
	case AddrP2WPKH: // native segwit witness(sig+pub) | (base: 41；witness: 107；总权重 ≈ 272 → vsize ≈ 68。)
		return 68
	case AddrP2WSH: // witness(sig(s)+script)
		return 109 // 仅估算 2-of-3 多签
	case AddrP2TR: // schnorr sig | (base 41；witness: 65；总权重 ≈ 208 → vsize ≈ 52。若 script path，则更大。)
		return 57
	default:
		return 148 // 默认按 P2PKH 估算
	}
}

// 根据地址类型速查交易输出的字节大小
// 返回值：1-scriptPubKey 长度, 2-输出大小vsize.
func GetOutSize(addrType AddressType) (int, int) {
	switch addrType {
	case AddrP2PKH: // 8 (金额) + 1 (len) + 25。
		return 25, 34
	case AddrP2SH: // 8 + 1 + 23。
		return 23, 32
	case AddrP2WPKH: // 8 + 1 + 22。
		return 22, 31
	case AddrP2WSH: // 8 + 1 + 34。
		return 34, 43
	case AddrP2TR: // 8 + 1 + 34。
		return 34, 43
	default:
		return 34, 43 // 默认按 AddrP2TR 估算
	}
}

// GetOpReturnSize 估算 OP_RETURN 输出大小
// 返回值：输出大小vsize.
func GetOpReturnSize(dataLen int) int {
	if dataLen <= 75 { // 1(OP_RETURN)+1(PUSHDATA)+dataLen
		return 11 + dataLen
	} else if dataLen <= 255 { // 1(OP_RETURN)+2(PUSHDATA1)+dataLen
		return 12 + dataLen
	} else if dataLen <= 65535 { // dataLen <= 520  // 1(OP_RETURN)+3(PUSHDATA2)+dataLen
		return 13 + dataLen
	} else { // 1(OP_RETURN)+4(PUSHDATA4)+dataLen
		return 14 + dataLen
	}
}
