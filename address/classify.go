// Package address 地址分类
package address

import (
	"github.com/crazycloudcc/btcapis/script"
)

// Classify 对给定的锁定脚本进行最小化分类，返回脚本类型字符串。
// 该函数复用 `script.Classify` 的判定逻辑，仅返回类型字符串，方便 address 包对外暴露简单 API。
func Classify(pkScript []byte) string {
	typ, _ := script.Classify(pkScript)
	return typ
}
