// 查询 bitcoind core 节点状态
package bitcoindrpc

import (
	"context"
	"fmt"
)

// 网络信息数据结构
type NetworkInfoDTO struct {
	Version            int      `json:"version"`            // 版本
	Subversion         string   `json:"subversion"`         // 子版本
	ProtocolVersion    int      `json:"protocolversion"`    // 协议版本
	LocalServices      string   `json:"localservices"`      // 本地服务
	LocalServicesNames []string `json:"localservicesnames"` // 本地服务名称
	LocalRelay         bool     `json:"localrelay"`         // 本地中继
	TimeOffset         int      `json:"timeoffset"`         // 时间偏移
	NetworkActive      bool     `json:"networkactive"`      // 网络活跃
	Connections        int      `json:"connections"`        // 连接数
	Connectionsin      int      `json:"connections_in"`     // 入站连接数
	Connectionsout     int      `json:"connections_out"`    // 出站连接数
	Networks           []struct {
		Name                      string `json:"name"`                        // 网络名称
		Limited                   bool   `json:"limited"`                     // 是否有限
		Reachable                 bool   `json:"reachable"`                   // 是否可达
		Proxy                     string `json:"proxy"`                       // 代理
		ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"` // 代理随机化凭证
	} `json:"networks"` // 网络列表
	Relayfee       float64 `json:"relayfee"`       // 中继费
	Incrementalfee float64 `json:"incrementalfee"` // 增量费
	Localaddresses []struct {
		Address string `json:"address"` // 地址
		Port    int    `json:"port"`    // 端口
		Score   int    `json:"score"`   // 分数
	} `json:"localaddresses"` // 本地地址列表
	Warnings string `json:"warnings"` // 警告
}

// 链信息数据结构
type ChainInfoDTO struct {
	Chain                string   `json:"chain"`                // 链
	Blocks               int      `json:"blocks"`               // 区块数
	Headers              int      `json:"headers"`              // 头数
	Bestblockhash        string   `json:"bestblockhash"`        // 最佳区块哈希
	Difficulty           float64  `json:"difficulty"`           // 难度
	Time                 int      `json:"time"`                 // 时间
	MedianTime           int      `json:"mediantime"`           // 中位时间
	Verificationprogress float64  `json:"verificationprogress"` // 验证进度
	Initialblockdownload bool     `json:"initialblockdownload"` // 初始区块下载
	Chainwork            string   `json:"chainwork"`            // 链工作量
	Sizeondisk           int      `json:"size_on_disk"`         // 磁盘大小
	Pruned               bool     `json:"pruned"`               // 是否修剪
	Warnings             []string `json:"warnings"`             // 警告列表
}

// 区块统计信息数据结构
type BlockStatsDTO struct {
	Avgfee             int    `json:"avgfee"`               // 平均交易费
	Avgfeerate         int    `json:"avgfeerate"`           // 平均交易费率
	Avgtxsize          int    `json:"avgtxsize"`            // 平均交易大小
	Blockhash          string `json:"blockhash"`            // 区块哈希
	Feeratepercentiles []int  `json:"feerate_percentiles"`  // 交易费率百分比
	Height             int    `json:"height"`               // 区块高度
	Ins                int    `json:"ins"`                  // 输入数
	Maxfee             int    `json:"maxfee"`               // 最大交易费
	Maxfeerate         int    `json:"maxfeerate"`           // 最大交易费率
	Maxtxsize          int    `json:"maxtxsize"`            // 最大交易大小
	Medianfee          int    `json:"medianfee"`            // 中位交易费
	Mediantime         int    `json:"mediantime"`           // 中位时间
	Mediantxsize       int    `json:"mediantxsize"`         // 中位交易大小
	Minfee             int    `json:"minfee"`               // 最小交易费
	Minfeerate         int    `json:"minfeerate"`           // 最小交易费率
	Mintxsize          int    `json:"mintxsize"`            // 最小交易大小
	Outs               int    `json:"outs"`                 // 输出数
	Subsidy            int    `json:"subsidy"`              // 挖矿奖励
	SwtotalSize        int    `json:"swtotal_size"`         // 简化交易大小
	SwtotalWeight      int    `json:"swtotal_weight"`       // 简化交易权重
	Swtxs              int    `json:"swtxs"`                // 简化交易数
	Time               int    `json:"time"`                 // 时间
	TotalOut           int    `json:"total_out"`            // 总输出
	TotalSize          int    `json:"total_size"`           // 总大小
	TotalWeight        int    `json:"total_weight"`         // 总权重
	Totalfee           int    `json:"totalfee"`             // 总交易费
	Txs                int    `json:"txs"`                  // 交易数
	UtxoIncrease       int    `json:"utxo_increase"`        // 未花费输出增加
	UtxoSizeInc        int    `json:"utxo_size_inc"`        // 未花费输出大小增加
	UtxoIncreaseActual int    `json:"utxo_increase_actual"` // 实际未花费输出增加
	UtxoSizeIncActual  int    `json:"utxo_size_inc_actual"` // 实际未花费输出大小增加
}

// 链顶信息数据结构
type ChainTipDTO struct {
	Height int64  `json:"height"` // 高度
	Hash   string `json:"hash"`   // 哈希
	Branch string `json:"branch"` // 分支
	Status string `json:"status"` // 状态
}

// 获取节点网络信息
func (c *Client) GetNetworkInfo(ctx context.Context) (NetworkInfoDTO, error) {
	var res NetworkInfoDTO
	if err := c.rpcCall(ctx, "getnetworkinfo", []any{}, &res); err != nil {
		return NetworkInfoDTO{}, fmt.Errorf("getnetworkinfo: %w", err)
	}

	return res, nil
}

// 获取链信息
func (c *Client) GetChainInfo(ctx context.Context) (ChainInfoDTO, error) {
	var res ChainInfoDTO
	if err := c.rpcCall(ctx, "getchaininfo", []any{}, &res); err != nil {
		return ChainInfoDTO{}, fmt.Errorf("getchaininfo: %w", err)
	}

	return res, nil
}

// 获取区块统计信息
func (c *Client) GetBlockStats(ctx context.Context, height int64) (BlockStatsDTO, error) {
	var res BlockStatsDTO
	if err := c.rpcCall(ctx, "getblockstats", []any{height}, &res); err != nil {
		return BlockStatsDTO{}, fmt.Errorf("getblockstats: %w", err)
	}

	return res, nil
}

func (c *Client) GetChainTip(ctx context.Context) (ChainTipDTO, error) {
	var res ChainTipDTO
	if err := c.rpcCall(ctx, "getchaintip", []any{}, &res); err != nil {
		return ChainTipDTO{}, fmt.Errorf("getchaintip: %w", err)
	}

	return res, nil
}
