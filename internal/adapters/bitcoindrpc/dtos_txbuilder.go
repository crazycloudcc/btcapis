// 构建tx结构 涉及到动态key值, 需要进行转换处理.
package bitcoindrpc

import (
	"encoding/json"
	"fmt"
)

type TxOutputScriptPubKeyCreateRawDTO struct {
	// 二选一：常用为 Hex。也允许携带 address（Core 会验证/转换）。
	Address string `json:"address,omitempty"`
	Hex     string `json:"hex,omitempty"`
}

// 构建交易输入
type TxInputCreateRawDTO struct {
	TxID     string `json:"txid"`     // 交易ID
	Vout     uint32 `json:"vout"`     // 输出索引
	Sequence uint32 `json:"sequence"` // 可选, 默认 0xffffffff
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
	Locktime    *int64                 `json:"locktime"`    // 可选, 默认 0; 非0值时, 交易锁定时间生效
	Replaceable *bool                  `json:"replaceable"` // 可选, 默认 true; 是否可替换
}

// MarshalJSON 按 Core 期望的“位置参数数组”编码： [vin, vout, locktime, replaceable]
func (p TxCreateRawDTO) MarshalJSON() ([]byte, error) {
	var arr []any
	arr = append(arr, p.Inputs)
	arr = append(arr, p.Outputs)
	// 仅当你需要第三/第四参数时才追加（避免与旧版/默认冲突）
	if p.Locktime != nil || p.Replaceable != nil {
		if p.Locktime == nil {
			var zero int64
			arr = append(arr, zero)
		} else {
			arr = append(arr, *p.Locktime)
		}
	}
	if p.Replaceable != nil {
		arr = append(arr, *p.Replaceable)
	}
	return json.Marshal(arr)
}

// 便捷构造器
func NewTxInput(txid string, vout uint32) TxInputCreateRawDTO {
	return TxInputCreateRawDTO{TxID: txid, Vout: vout}
}

func NewPayToAddress(addr string, amountBtc float64) TxOutputCreateRawDTO {
	return TxOutputCreateRawDTO{Address: addr, Amount: amountBtc}
}

func NewPayToScriptHex(scriptHex string, amountBtc float64) TxOutputCreateRawDTO {
	return TxOutputCreateRawDTO{
		Script: &TxOutputScriptPubKeyCreateRawDTO{Hex: scriptHex},
		Amount: amountBtc,
	}
}

func NewOpReturn(hexData string) TxOutputCreateRawDTO {
	return TxOutputCreateRawDTO{DataHex: hexData}
}

// 填充交易费用
type TxFundOptionsDTO struct {
	AddInputs      bool   `json:"add_inputs"`     // 可选, 默认 true; 是否自动添加更多输入
	IncludeUnsafe  bool   `json:"include_unsafe"` // 可选, 默认 false; 是否包含未确认的交易
	MinConf        int    `json:"minconf"`        // 可选, 默认 0; 如果 add_inputs 为 true, 则需要输入至少这么多确认
	MaxConf        int    `json:"maxconf"`        // 可选, 默认 0; 如果 add_inputs 为 true, 则需要输入最多这么多确认
	ChangeAddress  string `json:"changeAddress"`  // 可选, 默认 自动; 接收找零的地址
	ChangePosition int    `json:"changePosition"` // 可选, 默认 随机; 找零输出的索引
	// ChangeType             string     `json:"change_type"`            // 可选, 默认 由 -changetype 设置; 找零输出的类型
	IncludeWatching bool    `json:"includeWatching"` // 可选, 默认 true; 是否包含 watch-only 的输入
	LockUnspents    bool    `json:"lockUnspents"`    // 可选, 默认 false; 锁定选定的未花费输出
	FeeRateSats     float64 `json:"fee_rate"`        // 可选, 默认 0; 指定费用率, 单位 sat/vB
	// FeeRateBtc             float64    `json:"feeRate"`                // 可选, 默认 0; 指定费用率, 单位 BTC/kvB
	// SubtractFeeFromOutputs []int      `json:"subtractFeeFromOutputs"` // 可选, 默认 []; 从指定输出的金额中扣除费用
	// InputWeights           []struct { // 可选, 默认 []; 输入和对应的权重
	// 	TxID   string `json:"txid"`   // 交易ID
	// 	Vout   int    `json:"vout"`   // 输出索引
	// 	Weight int    `json:"weight"` // 权重
	// } `json:"input_weights"`
	// MaxTxWeight int `json:"max_tx_weight"` // 可选, 默认 400000; 最大可接受的交易权重
	// ConfTarget   int      `json:"conf_target"`   // 可选, 默认 由 -txconfirmtarget 设置; 确认目标区块数
	// EstimateMode string   `json:"estimate_mode"` // 可选, 默认 "unset"; 费用估计模式, 可选值: "unset", "economical", "conservative"
	Replaceable bool `json:"replaceable"` // 可选, 默认 false; 是否可替换
	// SolvingData struct { // 可选, 默认 {}; 需要签名数据
	// 	Pubkeys     []string `json:"pubkeys"`     // 公钥
	// 	Scripts     []string `json:"scripts"`     // 脚本
	// 	Descriptors []string `json:"descriptors"` // 描述符
	// } `json:"solving_data"`
}

// 填充交易费用结果
type TxFundRawResultDTO struct {
	Hex       string  `json:"hex"`       // 交易数据
	Fee       float64 `json:"fee"`       // 交易费用
	ChangePos int     `json:"changepos"` // 找零输出索引
	Change    string  `json:"change"`    // 找零地址
}
