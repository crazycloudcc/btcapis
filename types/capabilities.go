// Package types 后端能力定义
package types

// Capabilities 表示后端支持的能力
type Capabilities struct {
	// 基础能力
	HasChainReader  bool `json:"has_chain_reader"`
	HasBroadcaster  bool `json:"has_broadcaster"`
	HasFeeEstimator bool `json:"has_fee_estimator"`
	HasMempoolView  bool `json:"has_mempool_view"`

	// 高级能力
	HasBatchQueries bool `json:"has_batch_queries"`
	HasWebSocket    bool `json:"has_websocket"`
	HasCompression  bool `json:"has_compression"`
	HasCaching      bool `json:"has_caching"`

	// 网络支持
	Network         Network `json:"network"`
	SupportsSegWit  bool    `json:"supports_segwit"`
	SupportsTaproot bool    `json:"supports_taproot"`

	// 性能指标
	MaxConcurrentRequests int `json:"max_concurrent_requests"`
	RequestTimeout        int `json:"request_timeout"` // 秒
	RateLimit             int `json:"rate_limit"`      // 请求/秒

	// 数据一致性
	ProvidesConfirmedData bool `json:"provides_confirmed_data"`
	ProvidesMempoolData   bool `json:"provides_mempool_data"`
	DataFreshness         int  `json:"data_freshness"` // 秒
}

// IsReadOnly 检查是否为只读后端
func (c *Capabilities) IsReadOnly() bool {
	return c.HasChainReader && !c.HasBroadcaster
}

// IsFullNode 检查是否为全节点
func (c *Capabilities) IsFullNode() bool {
	return c.HasChainReader && c.HasBroadcaster && c.HasFeeEstimator && c.HasMempoolView
}

// SupportsNetwork 检查是否支持指定网络
func (c *Capabilities) SupportsNetwork(network Network) bool {
	return c.Network == network
}
