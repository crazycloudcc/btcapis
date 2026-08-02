# BTCAPIs

[English](README.md) | **简体中文**

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go 语言比特币链上操作库，面向 [chainbox](https://github.com/chainboxapp) 等上层服务。提供地址查询、交易构建/广播、PSBT、链 RPC、mempool 查询、费率估算与脚本解析，支持 Bitcoin Core RPC、mempool.space、ElectrumX 多数据源降级。

当前版本：**v0.5.0**（详见 [CHANGELOG.md](CHANGELOG.md)）

## 特性

- **地址**：余额/UTXO 查询、地址校验、钱包生成；支持 ElectrumX 批量查询与自定义 Provider 注入
- **交易**：原始交易查询、verbose 交易、mempool 预检、PSBT v0 创建/校验/广播
- **链 RPC**：网络信息、链状态、区块头/区块、按高度查 hash、verbosity 2 紧凑有序交易
- **费率**：统一 `types.ChainFees`（sat/vB）；默认降级链 + apis 扩展链（unisat/okx）
- **Mempool**：统计、交易状态、推荐费率、内存池 txid/entry
- **脚本/解码**：地址↔脚本互转、ASM 解析、原始交易解码
- **网络**：mainnet / testnet / signet

## 安装

```bash
go get github.com/crazycloudcc/btcapis@v0.5.0
```

## 快速开始

### 创建客户端

```go
package main

import (
    "context"
    "log"

    "github.com/crazycloudcc/btcapis"
)

func main() {
    client := btcapis.New(&btcapis.Config{
        Network:         "testnet",
        Timeout:         30,
        RPCUrl:          "http://127.0.0.1:18332",
        RPCUser:         "rpcuser",      // 从环境变量或配置文件注入，禁止硬编码
        RPCPass:         "rpcpassword",
        MempoolSpaceUrl: "https://mempool.space/testnet",
        ElectrumXUrl:    "", // 可选，如 https://blockstream.info/testnet/api
    })
    if client == nil {
        log.Fatal("client init failed")
    }

    ctx := context.Background()
    info, err := client.GetBlockChainInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("chain=%s blocks=%d", info.Chain, info.Blocks)
}
```

快捷构造（自动按网络填充 mempool.space URL）：

```go
client := btcapis.NewWithElectrumX(
    "testnet",
    "http://127.0.0.1:18332", "rpcuser", "rpcpassword",
    "wss://electrum.example.com:50002", // ElectrumX，可传空
    30,
)
```

> **安全**：RPC 密码、Unisat/OKX API Key 等凭证必须通过环境变量或本地配置文件注入，**禁止提交到 git**。见 `.gitignore`。

### 地址查询

```go
// 默认：mempool.space → bitcoin core
confirmed, mempool, err := client.GetAddressBalance(ctx, "tb1q...")
utxos, err := client.GetAddressUTXOs(ctx, "tb1q...")

// 带 ElectrumX 自定义 Provider（如 TCP 连接，由调用方实现 types.ElectrumXAddressProvider）
balance, err := client.GetAddressBalanceSatsWithOptions(ctx, "tb1q...", types.AddressProviderOptions{
    ElectrumX: myElectrumProvider,
})

addrType, err := client.DecodeAddressToType("tb1q...")
```

### 交易与 PSBT

```go
import "github.com/crazycloudcc/btcapis/types"

tx, err := client.GetTx(ctx, "txid...")
verbose, err := client.GetTxVerbose(ctx, "txid...", 1)

psbtBase64, err := client.CreatePSBT(ctx, &types.TxInputParams{
    FromAddress:   []string{"tb1p..."},
    ToAddress:     []string{"tb1q..."},
    AmountBTC:     []float64{0.001},
    FeeRate:       1.0,
    ChangeAddress: "tb1p...",
    Replaceable:   true,
})

txid, err := client.FinalizePSBTAndBroadcast(ctx, signedPSBT)
```

### 费率估算

```go
// 默认降级：mempool.space → electrumX → bitcoin core
fees, err := client.EstimateChainFees(ctx, 6)
// fees.High / Medium / Low / FeeRate（sat/vB，保留 2 位小数）

// apis 完整降级链（需自行创建扩展客户端并传入 API Key）
import (
    "github.com/crazycloudcc/btcapis/extensions/okx"
    "github.com/crazycloudcc/btcapis/extensions/unisat"
)

unisatClient, _ := unisat.New(unisat.Config{APIKey: os.Getenv("UNISAT_API_KEY")})
okxClient, _ := okx.New(okx.Config{
    APIKey:     os.Getenv("OKX_API_KEY"),
    SecretKey:  os.Getenv("OKX_SECRET_KEY"),
    Passphrase: os.Getenv("OKX_PASSPHRASE"),
})
fees, err = client.EstimateChainFeesForAPI(ctx, 6, btcapis.APIFeesProviderOptions{
    Unisat: unisatClient,
    OKX:    okxClient,
})

// 兼容旧接口
high, low, err := client.EstimateFeeRate(ctx, 6)
```

## 公开 API 一览

### 地址 `btcapis_address.go` / `btcapis_mempool.go`

| 方法 | 说明 |
|------|------|
| `CreateNewWallet` | 生成新钱包（助记词/地址/私钥） |
| `GetAddressBalance` | 确认/未确认余额（BTC） |
| `GetAddressUTXOs` | 地址 UTXO 列表 |
| `GetAddressBalanceSats` | 余额（聪，经默认 Provider） |
| `GetAddressUTXOsForAddress` | UTXO（`types.AddressUTXO` 格式） |
| `GetAddressBalanceSatsWithOptions` | 余额，electrumX → mempool.space |
| `GetAddressUTXOsWithOptions` | UTXO，electrumX → mempool.space |
| `GetAddressHistoryTxs` | 交易历史（需 `ElectrumXAddressProvider`） |
| `GetAddressUnconfirmedTxs` | 未确认交易（需 ElectrumX） |
| `BatchGetAddressBalanceSats` | 批量余额（需 ElectrumX） |
| `GetAddressBalanceWithElectrumX` | ElectrumX 余额 |
| `GetAddressBalanceWithElectrumXByXPRV` | 通过 xprv 派生地址查余额 |
| `GetAddressBalanceWithElectrumXByPrivateKey` | 通过 WIF 查余额 |
| `FilterAddressesWithBalanceWithElectrumX` | 批量过滤有余额地址 |
| `BatchGetBalancesWithElectrumX` | ElectrumX 批量余额 |
| `ValidateAddress` | Bitcoin Core `validateaddress` |

### 交易 `btcapis_tx.go`

| 方法 | 说明 |
|------|------|
| `GetTx` / `GetTxRaw` | 查询并解码交易 |
| `GetTxVerbose` | verbose 交易（core → mempool 降级） |
| `CreatePSBT` | 创建 PSBT v0 |
| `FinalizePSBTAndBroadcast` | 完成 PSBT 并广播 |
| `BroadcastRawTx` | 广播原始交易 |
| `ValidateUnsignedPsbtBase64` | 校验未签名 PSBT |
| `ValidateSignedPsbtBase64` | 校验已签名 PSBT |
| `TransferAllToNewAddress` | 全额转出（WIF 签名） |
| `TestMempoolAccept` | 内存池接受预检 |
| `GetMempoolTxIds` / `GetMempoolTxEntry` | 内存池交易查询 |

### 链 `btcapis_chain.go`

| 方法 | 说明 |
|------|------|
| `EstimateChainFees` | 默认费率降级链 |
| `EstimateChainFeesWithOptions` | 支持自定义 ElectrumX 费率 Provider |
| `EstimateChainFeesWithProviders` | 完全自定义 Provider 顺序 |
| `EstimateChainFeesForAPI` | apis 完整降级链 |
| `EstimateFeeRate` | 兼容旧接口（high/low sat/vB） |
| `GetNetworkInfo` / `GetBlockChainInfo` | 节点/链状态 |
| `GetBlockStats` / `GetChainTips` | 区块统计/分叉信息 |
| `GetBlockHeader` / `GetBlock` | 区块头/区块（verbosity=1） |
| `GetBlockTransactions` | 紧凑有序交易（verbosity=2，精确 fee sats） |
| `GetBlockHashByHeight` | 按高度查区块 hash |

### Mempool `btcapis_mempool.go`

| 方法 | 说明 |
|------|------|
| `GetMempoolStats` | 内存池统计 |
| `GetMempoolTxStatus` | 交易确认状态 |
| `GetMempoolFeesRecommend` | mempool.space 推荐费率 |

### 脚本/解码 `btcapis_scripts.go`

| 方法 | 说明 |
|------|------|
| `DecodeAddressToScriptInfo` / `DecodeAddressToPkScript` / `DecodeAddressToType` | 地址解析 |
| `DecodePkScriptToAddressInfo` / `DecodePKScriptToType` | 脚本→地址 |
| `DecodePkScriptToAsmString` | 脚本→操作码/ASM |
| `DecodeRawTx` / `DecodeRawTxString` | 原始交易解码 |

### 格式校验 `btcapis_checkers.go`

`CheckFormatAddress` / `CheckFormatTxid` / `CheckFormatHex` / `CheckFormatBase64`

### 扩展包

| 包 | 用途 |
|----|------|
| `extensions/unisat` | Unisat Open API 费率 Provider |
| `extensions/okx` | OKX Web3 API 费率 Provider |

### 自定义 Provider 接口 `types/`

```go
// 地址数据源（调用方实现，如 TCP ElectrumX）
type ElectrumXAddressProvider interface {
    Name() string
    GetBalanceSats(ctx context.Context, addr string) (*AddressBalanceSats, error)
    GetUTXOs(ctx context.Context, addr string) ([]AddressUTXO, error)
    GetHistoryTxs(ctx context.Context, addr string) ([]AddressHistoryTx, error)
    GetUnconfirmedTxs(ctx context.Context, addr string) ([]AddressUnconfirmedTx, error)
    BatchGetBalanceSats(ctx context.Context, addrs []string) ([]AddressBalanceSats, error)
}

// 费率数据源
type ChainFeesEstimator interface { ... }
```

## 数据源配置

`Config` 字段按需填写，未配置的适配器不会初始化：

| 字段 | 说明 |
|------|------|
| `Network` | `mainnet` / `testnet` / `signet` |
| `RPCUrl` / `RPCUser` / `RPCPass` | Bitcoin Core JSON-RPC |
| `MempoolSpaceUrl` | mempool.space API 根地址 |
| `ElectrumXUrl` | ElectrumX HTTP 端点 |
| `Timeout` | HTTP/RPC 超时（秒） |

mempool.space 默认地址（`NewWithElectrumX` 自动选择）：

- mainnet → `https://mempool.space`
- testnet → `https://mempool.space/testnet`
- signet → `https://mempool.space/signet`

## 地址类型

| 类型 | 前缀 | 支持 |
|------|------|------|
| P2PKH | `1...` / `m...` / `n...` | ✅ |
| P2SH | `3...` / `2...` | ✅ |
| P2WPKH / P2WSH | `bc1q...` / `tb1q...` | ✅ |
| P2TR (Taproot) | `bc1p...` / `tb1p...` | ✅ |

## 项目结构

```
btcapis.go / btcapis_*.go     # 公开门面 API
types/                        # 对外数据类型与 Provider 接口
extensions/unisat|okx/        # 可选费率扩展
internal/
  adapters/
    bitcoindrpc/              # Bitcoin Core RPC
    mempoolapis/              # mempool.space REST
    electrumx/                # ElectrumX HTTP
  address/                    # 地址查询与 orchestration
  chain/                      # 链 RPC、费率、mempool
  tx/                         # 交易构建、PSBT、广播
  decoders/                   # 地址/脚本/交易解码
  ordinals/ / runes/          # 内部实现，暂未导出公开 API
pkg/logger/                   # 日志
docs/                         # ElectrumX 文档与 RPC 接口清单
```

```
┌─────────────────────────────────────────┐
│  门面层  btcapis.go, btcapis_*.go        │
├─────────────────────────────────────────┤
│  领域层  internal/address|chain|tx       │
├─────────────────────────────────────────┤
│  适配器  internal/adapters/*             │
│  扩展    extensions/unisat|okx           │
└─────────────────────────────────────────┘
```

## 开发

```bash
git clone https://github.com/crazycloudcc/btcapis.git
cd btcapis
export GOCACHE=$PWD/.gocache
go mod tidy
go test ./...
golangci-lint run   # 需本地安装
```

集成测试可使用 `TestClient`（`btcapis_tests.go`），直接访问底层 bitcoindrpc/mempool/electrumx 客户端，需配置真实节点。

## 文档

| 文件 | 内容 |
|------|------|
| [CHANGELOG.md](CHANGELOG.md) | 版本变更 |
| [docs/electrumx_quickstart.md](docs/electrumx_quickstart.md) | ElectrumX 快速入门 |
| [docs/electrumx_implementation.md](docs/electrumx_implementation.md) | ElectrumX 实现说明 |
| [docs/2.bitcoin-core-rpc接口清单.txt](docs/2.bitcoin-core-rpc接口清单.txt) | Bitcoin Core RPC 清单 |
| [docs/3.mempool-space接口清单.txt](docs/3.mempool-space接口清单.txt) | mempool.space API 清单 |

## 依赖

- [btcsuite/btcd](https://github.com/btcsuite/btcd) — 交易/脚本/地址底层
- [shopspring/decimal](https://github.com/shopspring/decimal) — 精度计算
- [sirupsen/logrus](https://github.com/sirupsen/logrus) — 日志

公共 API 使用 `types/` 自定义类型，与 btcd 类型解耦。

## 开源协议

本项目采用 [MIT License](LICENSE) 开源。
