package script

// Classify 仅做最小化识别：p2pkh/p2sh/p2wpkh/p2wsh/p2tr，其余返回 "unknown"。
// 需要生成地址时，再扩展为传入 network，并用 bech32/base58 生成。
func Classify(pkScript []byte) (typ string, addrs []string) {
	b := pkScript
	n := len(b)
	// P2PKH: OP_DUP OP_HASH160 PUSH20 <20> OP_EQUALVERIFY OP_CHECKSIG
	if n == 25 && b[0] == 0x76 && b[1] == 0xa9 && b[2] == 0x14 && b[23] == 0x88 && b[24] == 0xac {
		return "p2pkh", nil
	}
	// P2SH: OP_HASH160 PUSH20 <20> OP_EQUAL
	if n == 23 && b[0] == 0xa9 && b[1] == 0x14 && b[22] == 0x87 {
		return "p2sh", nil
	}
	// P2WPKH v0: OP_0 PUSH20 <20>
	if n == 22 && b[0] == 0x00 && b[1] == 0x14 {
		return "p2wpkh", nil
	}
	// P2WSH v0: OP_0 PUSH32 <32>
	if n == 34 && b[0] == 0x00 && b[1] == 0x20 {
		return "p2wsh", nil
	}
	// P2TR v1: OP_1 PUSH32 <32>
	if n == 34 && b[0] == 0x51 && b[1] == 0x20 {
		return "p2tr", nil
	}
	return "unknown", nil
}
