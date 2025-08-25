package tx

// 通用的转账交易输入参数
type TxInputParams struct {
	FromAddress []string  `json:"from_address"` // 来源地址数组-可以是多个, 但是目前版本只支持1个地址
	ToAddress   []string  `json:"to_address"`   // 目标地址数组-可以是多个, 但是要和Amount一一对应
	AmountBTC   []float64 `json:"amount"`       // 金额-单位BTC
	FeeRate     float64   `json:"fee_rate"`     // 费用率(sat/vB)
	Locktime    int64     `json:"locktime"`     // 锁定时间(秒)
	Replaceable bool      `json:"replaceable"`  // 是否可替换RBF
	Data        string    `json:"data"`         // 可选 交付附加数据
	PublicKey   string    `json:"public_key"`   // 公钥 => 从OKX获取, 后续要删除, 改用其他方式录入钱包
}
