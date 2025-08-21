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

// 内部：LE 存储
func (h Hash32) BytesLE() []byte {
	b := make([]byte, 32)
	copy(b, h[:])
	return b
}
func (h Hash32) BytesBE() []byte {
	b := make([]byte, 32)
	for i := 0; i < 32; i++ {
		b[i] = h[31-i]
	}
	return b
}

// 人类可读（与 mempool.space / bitcoin-cli 对齐）：BE hex
func (h Hash32) String() string { // 或 TxIDBE()
	return hex.EncodeToString(h.BytesLE()) // gpt-5推荐BE, 但是mempool.space和LE一致.
}

// 解析来自外部（BE 字符串）
func (h *Hash32) FromBEHex(s string) error {
	s = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(s), "0x"))
	if len(s) != 64 {
		return fmt.Errorf("invalid hash length: %d", len(s))
	}
	be, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	for i := 0; i < 32; i++ {
		h[i] = be[31-i]
	} // 反转写入 LE
	return nil
}

// JSON 始终输出 BE；输入接受 BE
func (h Hash32) MarshalJSON() ([]byte, error) { return json.Marshal(h.String()) }
func (h *Hash32) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return h.FromBEHex(s)
}
