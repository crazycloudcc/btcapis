# Changelog

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
