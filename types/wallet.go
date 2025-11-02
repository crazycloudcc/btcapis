package types

// WalletInfo 钱包信息结构
type WalletInfo struct {
	P2PKH      string `json:"p2pkh"`
	P2PSH      string `json:"p2psh"`
	P2WPKH     string `json:"p2wpkh"`
	P2TR       string `json:"p2tr"`
	Mnemonic   string `json:"mnemonic"`
	XPRV       string `json:"xprv"`
	BTCBalance int64  `json:"btc_balance,omitempty"`
}
