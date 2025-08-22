package tx

// 交易输入参数
type TxInputParams struct {
	FromAddress string  `json:"from_address"` // 来源地址
	ToAddress   string  `json:"to_address"`   // 目标地址
	Amount      float64 `json:"amount"`       // 金额
	FeeRate     float64 `json:"fee_rate"`     // 费用率(sat/vB)
	Locktime    int64   `json:"locktime"`     // 锁定时间(秒)
	Replaceable bool    `json:"replaceable"`  // 是否可替换
}
