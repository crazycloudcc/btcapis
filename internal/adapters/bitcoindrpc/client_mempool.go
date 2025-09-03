// 这是bitcoin core节点的内存池, 要和mempool.space区分开!
package bitcoindrpc

import (
	"context"
	"fmt"
)

type MempoolInfoDTO struct {
	Loaded              bool    `json:"loaded"`              // 是否加载
	Size                int     `json:"size"`                // 大小
	Bytes               int     `json:"bytes"`               // 字节
	Usage               int     `json:"usage"`               // 使用
	TotalFee            float64 `json:"total_fee"`           // 总费用
	MaxMempool          int     `json:"maxmempool"`          // 最大内存池
	MempoolMinFee       float64 `json:"mempoolminfee"`       // 最小内存池费用
	MinRelayTxFee       float64 `json:"minrelaytxfee"`       // 最小中继交易费用
	IncrementalRelayFee float64 `json:"incrementalrelayfee"` // 增量中继交易费用
	UnbroadcastCount    int     `json:"unbroadcastcount"`    // 未广播数量
	FullRBF             bool    `json:"fullrbf"`             // 全RBF
}

// 内存池交易信息数据结构
type MempoolTxDTO struct {
	Vsize           int    `json:"vsize"`           // 剥离大小
	Weight          int    `json:"weight"`          // 权重
	Time            int    `json:"time"`            // 时间
	Height          int    `json:"height"`          // 高度
	DescendantCount int    `json:"descendantcount"` // 后代数量
	DescendantSize  int    `json:"descendantsize"`  // 后代大小
	AncestorCount   int    `json:"ancestorcount"`   // 祖先数量
	AncestorSize    int    `json:"ancestorsize"`    // 祖先大小
	Wtxid           string `json:"wtxid"`           // 交易ID
	Fees            struct {
		Base       float64 `json:"base"`       // 基础费用
		Modified   float64 `json:"modified"`   // 修改后的费用
		Ancestor   float64 `json:"ancestor"`   // 祖先费用
		Descendant float64 `json:"descendant"` // 后代费用
	} `json:"fees"` // 费用
	Depends           []string `json:"depends"`            // 依赖
	Spentby           []string `json:"spentby"`            // 花费
	Bip125Replaceable bool     `json:"bip125-replaceable"` // 可替换
	Unbroadcast       bool     `json:"unbroadcast"`        // 未广播
}

// 获取内存池信息
func (c *Client) MempoolGetInfo(ctx context.Context) (*MempoolInfoDTO, error) {
	var res *MempoolInfoDTO
	if err := c.rpcCall(ctx, "getmempoolinfo", []any{}, &res); err != nil {
		return nil, fmt.Errorf("getmempoolinfo: %w", err)
	}
	return res, nil
}

// 获取内存池交易信息
func (c *Client) MempoolGetTxs(ctx context.Context) ([]string, error) {
	var res []string
	var flag bool = false // false-返回txid数组; true-返回tx的详细数据json;
	if err := c.rpcCall(ctx, "getrawmempool", []any{flag}, &res); err != nil {
		return nil, fmt.Errorf("getrawmempool: %w", err)
	}
	return res, nil
}

// 获取内存池交易信息
func (c *Client) MempoolGetTx(ctx context.Context, txid string) (*MempoolTxDTO, error) {
	var res *MempoolTxDTO
	if err := c.rpcCall(ctx, "getmempoolentry", []any{txid}, &res); err != nil {
		return nil, fmt.Errorf("getmempoolentry: %w", err)
	}
	return res, nil
}
