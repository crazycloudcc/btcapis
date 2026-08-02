# Changelog

## v0.5.0 (2026-08-03)

### Added
- `Client.GetBlockTransactions`：固定调用 Bitcoin Core `getblock(hash, 2)`，返回紧凑有序交易与精确 `fee_sats`
- `types.BlockTransactions` / `types.BlockTransaction`
- `ErrBlockNotFound` / `ErrInvalidBlockData` 类型化错误

### Fixed
- Bitcoin Core JSON-RPC 请求改用固定 ID，移除并发调用的数据竞争
- 非 coinbase 缺失手续费、fractional sat、负数与 int64 溢出不再进入成功结果

## v0.4.0 (2026-07-08)

### Added
- 链 RPC 公开 API：`GetNetworkInfo` / `GetBlockStats` / `GetChainTips` / `GetBlockChainInfo` / `GetBlockHeader` / `GetBlock` / `GetBlockHashByHeight`
- 地址 orchestration：`GetAddressBalanceSatsWithOptions` / `GetAddressUTXOsWithOptions` / `BatchGetAddressBalanceSats` / `GetAddressHistoryTxs` / `GetAddressUnconfirmedTxs`
- mempool.space 公开 API：`GetMempoolStats` / `GetMempoolTxStatus` / `GetMempoolFeesRecommend`
- 交易：`GetTxVerbose` / `TestMempoolAccept` / `GetMempoolTxIds` / `GetMempoolTxEntry` / `ValidateAddress`
- `types.ElectrumXAddressProvider`：支持 apis TCP ElectrumX 注入

### Changed
- apis / schedulers 链通用能力委托至本库；本地重复 RPC 实现已清理

## v0.3.0 (2026-07-07)

### Added
- `types.ChainFees`：统一链上手续费结构（sat/vB，2 位小数）
- `Client.EstimateChainFees`：默认降级链 mempool.space → electrumX → bitcoin core
- `Client.EstimateChainFeesForAPI`：apis 完整降级链（含 unisat / okx 扩展）
- `extensions/unisat`、`extensions/okx`：推荐费率 Provider
- `internal/chain` 费率规范化与 Provider 接口

### Fixed
- Bitcoin Core / ElectrumX 费率单位：BTC/kB 正确转换为 sat/vB（×1e5）
- 旧版 `EstimateFeeRate` 误用 `BTCToSats` 处理 feerate 的问题
- `pkg/logger` 在未初始化时不 panic
- mempool 推荐费率 DTO 补齐 `economyFee` / `minimumFee`

### Changed
- `chain.Client` 现包含 electrumX 客户端，用于费率回退
- `internal/tx` 内部费率估算统一走 `chain.EstimateChainFeesDefault`
