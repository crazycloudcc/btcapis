// 构建tx结构 涉及到动态key值, 需要进行转换处理.
package bitcoindrpc

import (
	"encoding/json"
	"fmt"
)

// 构建交易输入
type TxInputCreateRawDTO struct {
	TxID     string `json:"txid"`     // 交易ID
	Vout     uint32 `json:"vout"`     // 输出索引
	Sequence uint32 `json:"sequence"` // 可选, 默认 0xffffffff
}

type TxOutputScriptPubKeyCreateRawDTO struct {
	// 二选一：常用为 Hex。也允许携带 address（Core 会验证/转换）。
	Address string `json:"address,omitempty"`
	Hex     string `json:"hex,omitempty"`
}

// 构建交易输出
type TxOutputCreateRawDTO struct {
	Address string                            `json:"address,omitempty"` // 地址
	Script  *TxOutputScriptPubKeyCreateRawDTO `json:"script,omitempty"`  // 脚本公钥
	Amount  float64                           `json:"amount,omitempty"`  // 金额
	DataHex string                            `json:"datahex,omitempty"` // OP_RETURN 不带金额
}

// 单个输出对象的 JSON 编码（与 Core 接口格式对齐）
func (o TxOutputCreateRawDTO) MarshalJSON() ([]byte, error) {
	switch {
	// OP_RETURN：{"data":"<hex>"}
	case o.DataHex != "":
		return json.Marshal(map[string]any{
			"data": o.DataHex,
		})
	// 直付脚本：{"scriptPubKey":{"hex":"..."}, "amount": <btc number>}
	case o.Script != nil:
		return json.Marshal(map[string]any{
			"scriptPubKey": o.Script,
			"amount":       o.Amount,
		})
	// 直付地址：{"<address>": <btc number>}
	case o.Address != "":
		// 这里必须用动态 key
		m := map[string]any{o.Address: o.Amount}
		return json.Marshal(m)
	default:
		return nil, fmt.Errorf("empty CreateRawTxOutput")
	}
}

// 构建交易
type TxCreateRawDTO struct {
	Inputs      []TxInputCreateRawDTO  `json:"inputs"`      // 交易输入
	Outputs     []TxOutputCreateRawDTO `json:"outputs"`     // 交易输出
	Locktime    int64                  `json:"locktime"`    // 可选, 默认 0; 非0值时, 交易锁定时间生效
	Replaceable bool                   `json:"replaceable"` // 可选, 默认 true; 是否可替换
}
