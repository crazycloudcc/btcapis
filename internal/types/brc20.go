package types

// BRC20Action 结构体：存储 BRC-20 动作
type BRC20Action struct {
	Op   string `json:"op"`            // 操作类型
	Tick string `json:"tick"`          // 代币符号
	Amt  string `json:"amt,omitempty"` // 数量
	Max  string `json:"max,omitempty"` // 最大数量
	Lim  string `json:"lim,omitempty"` // 限制
}
