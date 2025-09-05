package tests

import "github.com/crazycloudcc/btcapis/types"

func getTxid() string {
	txid := ""
	switch types.CurrentNetwork {
	case types.Mainnet:
		txid = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b" // btc链创世交易
	// txid := "8e2955284d9f66f56df9c89fedd50103fe68f84ec9f138e4e20c67db15de68ee" // main 839639
	// txid := "2aafdfd00e5a27ef211d7e6ed76a14dd7edd6b88a19df2e535f0963c05372a8d" // main 910115
	case types.Signet:
		txid = "12ad394188be183dccc53357a4374a7eab067810c3535063012986dc437e8a3c"
		// txid := "0d31e59675c85f17d942f4510bb4760d9ed4b661df22af3b7cd5ef3c2116626b" // 测试交易
	case types.Testnet:
		txid = "" // 空交易
	case types.Regtest:
		txid = ""
	}
	return txid
}

func getAddr() string {
	addr := "" // 测试地址
	switch types.CurrentNetwork {
	case types.Mainnet:
		addr = "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq"
		// addr := "1xxxxxxxxxxxxxxxxxxxxxxxxxy1kmdGr" // P2PKH 示例地址
		// addr := "3DHgSaYsxCj62UKU2yFG3vKbwjua3ViHUS" // P2SH 示例地址
		// addr := "bc1qgnmdx4pyaxrkhtgeqgh0g93cvar7achq8kjtnm" // P2WPKH 示例地址 - 0个UTXO
		// addr := "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq" // P2TR 示例地址
	case types.Signet:
		addr = "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw"
	case types.Testnet:
		addr = "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw"
	case types.Regtest:
		addr = ""
	}
	return addr
}
