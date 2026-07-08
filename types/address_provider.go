package types

import "context"

// ElectrumXAddressProvider 自定义 ElectrumX 地址数据源（如 TCP 连接）
type ElectrumXAddressProvider interface {
	Name() string
	GetBalanceSats(ctx context.Context, addr string) (*AddressBalanceSats, error)
	GetUTXOs(ctx context.Context, addr string) ([]AddressUTXO, error)
	GetHistoryTxs(ctx context.Context, addr string) ([]AddressHistoryTx, error)
	GetUnconfirmedTxs(ctx context.Context, addr string) ([]AddressUnconfirmedTx, error)
	BatchGetBalanceSats(ctx context.Context, addrs []string) ([]AddressBalanceSats, error)
}

// AddressProviderOptions 地址查询 Provider 配置
type AddressProviderOptions struct {
	ElectrumX ElectrumXAddressProvider
}
