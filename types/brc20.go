package types

// BRC20Action represents a parsed BRC-20 inscription body.
type BRC20Action struct {
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt,omitempty"`
	Max  string `json:"max,omitempty"`
	Lim  string `json:"lim,omitempty"`
}
