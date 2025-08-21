package script

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
