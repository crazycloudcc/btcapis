package bitcoindrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// RPCWarnings 兼容 Bitcoin Core 新旧 RPC：
// - <28：warnings 为 string
// - >=28：warnings 为 []string
type RPCWarnings []string

func (w *RPCWarnings) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*w = nil
		return nil
	}
	switch trimmed[0] {
	case '"':
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return fmt.Errorf("warnings string: %w", err)
		}
		if s == "" {
			*w = nil
			return nil
		}
		*w = RPCWarnings{s}
		return nil
	case '[':
		var list []string
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return fmt.Errorf("warnings array: %w", err)
		}
		*w = RPCWarnings(list)
		return nil
	default:
		return fmt.Errorf("warnings: 期望 string 或 []string，实际: %s", string(trimmed))
	}
}

// AsAny 供对外 types.Warnings(any) 使用；空则 nil。
func (w RPCWarnings) AsAny() any {
	if len(w) == 0 {
		return nil
	}
	if len(w) == 1 {
		return w[0]
	}
	return []string(w)
}
