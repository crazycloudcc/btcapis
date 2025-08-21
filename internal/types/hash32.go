package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Hash32 是 32 字节哈希（双 SHA256 结果）。
// 约定：内部按“比特币内部顺序”存（通常等同于小端存放），
// JSON 显示时按“大端十六进制字符串”（常见 txid 书写）输出。
type Hash32 [32]byte

// TxIDBE 返回人类可读的大端十六进制字符串（txid/块哈希常见写法）
func (h Hash32) TxIDBE() string {
	var be [32]byte
	for i := 0; i < 32; i++ {
		be[i] = h[31-i]
	}
	return hex.EncodeToString(be[:])
}

// FromTxIDBE 把大端十六进制字符串写回内部字节序（反转）
func (h *Hash32) FromTxIDBE(s string) error {
	s = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(s)), "0x")
	if len(s) != 64 {
		return fmt.Errorf("invalid hash length: %d", len(s))
	}
	be, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	for i := 0; i < 32; i++ {
		h[i] = be[31-i]
	}
	return nil
}

// 自定义 JSON：输出为十六进制字符串（大端）
func (h Hash32) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.TxIDBE())
}

// 允许用十六进制字符串（大端）反序列化
func (h *Hash32) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return h.FromTxIDBE(s)
}
