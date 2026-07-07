# BTCAPIs

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-detailed-green.svg)](docs/)

一个功能全面的 Go 语言比特币 API 库，提供地址处理、交易操作、PSBT 管理、脚本解析等核心功能。支持多种数据源，包括 Bitcoin Core RPC 和 Mempool.space API。

## 🚀 特性

### 核心功能

- **地址处理** - 支持所有比特币地址类型（P2PKH、P2SH、P2WPKH、P2WSH、P2TR）
- **交易操作** - 创建、签名、广播比特币交易
- **PSBT 支持** - 完整的 PSBT（部分签名比特币交易）工厂实现
- **脚本解析** - 比特币脚本编码/解码与分析
- **UTXO 管理** - 查询和管理未花费交易输出
- **费率估算** - 智能费率计算和优化

### 数据源支持

- **Bitcoin Core RPC** - 直接连接比特币核心节点
- **Mempool.space API** - 支持主网、测试网、Signet
- **多后端架构** - 可扩展的提供商系统

### 网络支持

- ✅ **主网 (Mainnet)**
- ✅ **测试网 (Testnet)**
- ✅ **Signet**

## 📦 安装

```bash
go get github.com/crazycloudcc/btcapis
```

## 🔧 快速开始

### 基础初始化

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/crazycloudcc/btcapis"
)

func main() {
    // 创建客户端连接
    client := btcapis.New(
        "testnet",                    // 网络: mainnet, testnet, signet
        "http://localhost:18332",     // Bitcoin Core RPC URL
        "rpcuser",                    // RPC 用户名
        "rpcpassword",               // RPC 密码
        30,                          // 超时时间(秒)
    )

    ctx := context.Background()

    // 检查连接
    blockCount, err := client.GetBlockCount(ctx)
    if err != nil {
        log.Fatal("连接失败:", err)
    }

    log.Printf("当前区块高度: %d", blockCount)
}
```

### 地址操作

```go
// 查询地址余额
confirmed, mempool, err := client.GetAddressBalance(ctx, "tb1q...")
if err != nil {
    log.Fatal(err)
}
log.Printf("确认余额: %d satoshi, 未确认: %d satoshi", confirmed, mempool)

// 获取地址 UTXO
utxos, err := client.GetAddressUTXOs(ctx, "tb1q...")
if err != nil {
    log.Fatal(err)
}
log.Printf("UTXO 数量: %d", len(utxos))

// 地址类型解析
addrType, err := client.DecodeAddressToType("tb1q...")
if err != nil {
    log.Fatal(err)
}
log.Printf("地址类型: %s", addrType)
```

### 交易操作

```go
// 查询交易信息
tx, err := client.GetTx(ctx, "交易ID")
if err != nil {
    log.Fatal(err)
}
log.Printf("交易版本: %d, 输入数: %d, 输出数: %d",
    tx.Version, len(tx.TxIn), len(tx.TxOut))

// 获取原始交易数据
rawTx, err := client.GetTxRaw(ctx, "交易ID")
if err != nil {
    log.Fatal(err)
}
log.Printf("原始交易大小: %d bytes", len(rawTx))
```

### PSBT 交易创建

```go
import "github.com/crazycloudcc/btcapis/types"

// 构建交易参数
txParams := &types.TxInputParams{
    FromAddress: []string{"tb1p..."},          // 发送地址
    ToAddress:   []string{"tb1q..."},          // 接收地址
    AmountBTC:   []float64{0.001},             // 金额 (BTC)
    FeeRate:     1.0,                          // 费率 (sat/vB)
    Locktime:    0,                            // 锁定时间
    Replaceable: true,                         // 支持 RBF
    Data:        "Hello Bitcoin",              // 可选数据 (OP_RETURN)
    PublicKey:   "公钥十六进制",                // 公钥
    ChangeAddress: "tb1p...",                  // 找零地址
}

// 创建 PSBT
psbtBase64, err := client.CreatePSBT(ctx, txParams)
if err != nil {
    log.Fatal(err)
}
log.Printf("PSBT: %s", psbtBase64)

// 完成签名并广播 (需要外部签名)
signedPSBT := "..." // 签名后的 PSBT
txid, err := client.FinalizePSBTAndBroadcast(ctx, signedPSBT)
if err != nil {
    log.Fatal(err)
}
log.Printf("交易已广播: %s", txid)
```

### 链信息查询

```go
// 统一费率结构（sat/vB，保留 2 位小数）
fees, err := client.EstimateChainFees(ctx, 6)
if err != nil {
    log.Fatal(err)
}
log.Printf("high=%.2f medium=%.2f low=%.2f feerate=%.2f blocks=%d",
    fees.High, fees.Medium, fees.Low, fees.FeeRate, fees.Blocks)

// apis 完整降级链（mempool → unisat → okx → electrumX → bitcoin core）
fees, err = client.EstimateChainFeesForAPI(ctx, 6, btcapis.APIFeesProviderOptions{
    Unisat: unisatClient,
    OKX:    okxClient,
})
```

```go
// 兼容旧接口
fastRate, economyRate, err := client.EstimateFeeRate(ctx, 6)
if err != nil {
    log.Fatal(err)
}
log.Printf("快速费率: %.2f sat/vB, 经济费率: %.2f sat/vB", fastRate, economyRate)

// 获取最新区块哈希
bestHash, err := client.GetBestBlockHash(ctx)
if err != nil {
    log.Fatal(err)
}
log.Printf("最新区块: %s", bestHash)

// 查询区块信息
block, err := client.GetBlock(ctx, bestHash)
if err != nil {
    log.Fatal(err)
}
log.Printf("区块高度: %d, 交易数: %d", block.Height, len(block.Tx))
```

## 📋 完整功能列表

### 🏠 地址模块 (Address)

| 功能      | 方法                            | 描述                      |
| --------- | ------------------------------- | ------------------------- |
| 余额查询  | `GetAddressBalance()`           | 查询地址的确认/未确认余额 |
| UTXO 查询 | `GetAddressUTXOs()`             | 获取地址的未花费输出      |
| 地址解析  | `DecodeAddressToScriptInfo()`   | 解析地址的详细脚本信息    |
| 脚本转换  | `DecodeAddressToPkScript()`     | 地址转锁定脚本            |
| 类型识别  | `DecodeAddressToType()`         | 识别地址类型              |
| 脚本解析  | `DecodePkScriptToAddressInfo()` | 脚本转地址信息            |

### 💸 交易模块 (Transaction)

| 功能      | 方法                          | 描述                   |
| --------- | ----------------------------- | ---------------------- |
| 交易查询  | `GetTx()`                     | 获取交易详细信息       |
| 原始数据  | `GetTxRaw()`                  | 获取交易原始字节数据   |
| PSBT 创建 | `CreatePSBT()`                | 创建部分签名比特币交易 |
| PSBT 完成 | `FinalizePSBTAndBroadcast()`  | 完成签名并广播         |
| 交易广播  | `BroadcastRawTx()`            | 广播原始交易           |
| 地址导入  | `ImportAddressAndPublickey()` | 导入地址和公钥         |

### ⛓️ 区块链模块 (Chain)

| 功能      | 方法                 | 描述                 |
| --------- | -------------------- | -------------------- |
| 费率估算  | `EstimateFeeRate()`  | 估算交易费率         |
| UTXO 查询 | `GetUTXO()`          | 查询特定 UTXO 状态   |
| 区块统计  | `GetBlockCount()`    | 获取区块链高度       |
| 最新区块  | `GetBestBlockHash()` | 获取最新区块哈希     |
| 区块哈希  | `GetBlockHash()`     | 根据高度获取区块哈希 |
| 区块头    | `GetBlockHeader()`   | 获取区块头信息       |
| 区块数据  | `GetBlock()`         | 获取完整区块信息     |

### 🔧 脚本模块 (Script)

| 功能     | 方法                      | 描述             |
| -------- | ------------------------- | ---------------- |
| 脚本解析 | `DecodeScriptToOpcodes()` | 解析脚本为操作码 |
| ASM 转换 | `DecodeScriptToASM()`     | 脚本转汇编格式   |
| 类型检测 | `DecodePKScriptToType()`  | 检测脚本类型     |

## 🏗️ 架构设计

BTCAPIs 采用三层架构模式，提供清晰的关注点分离：

```
┌─────────────────────────────────────────────────────────────┐
│                    门面层 (Facade)                          │
│              btcapis.go, *_facade.go                        │
├─────────────────────────────────────────────────────────────┤
│                    端口层 (Ports)                           │
│                chain/backend.go                             │
├─────────────────────────────────────────────────────────────┤
│                   适配器层 (Adapters)                       │
│         providers/bitcoindrpc, mempoolspace                │
└─────────────────────────────────────────────────────────────┘
```

### 核心模块

- **`types/`** - 核心数据类型定义 (地址、交易、UTXO 等)
- **`internal/adapters/`** - 数据源适配器 (Bitcoin Core, Mempool.space)
- **`internal/address/`** - 地址处理逻辑
- **`internal/tx/`** - 交易构建和管理
- **`internal/chain/`** - 区块链交互
- **`internal/decoders/`** - 编码解码工具

## 🌐 数据源配置

### Bitcoin Core RPC

```go
client := btcapis.New(
    "mainnet",
    "http://localhost:8332",  // RPC 地址
    "rpcuser",                // 用户名
    "rpcpassword",           // 密码
    30,                      // 超时秒数
)
```

### Mempool.space API

自动根据网络配置：

- **主网**: `https://mempool.space`
- **测试网**: `https://mempool.space/testnet`
- **Signet**: `https://mempool.space/signet`

## 📖 地址类型支持

| 地址类型           | 前缀      | 示例                                                             | 支持状态    |
| ------------------ | --------- | ---------------------------------------------------------------- | ----------- |
| P2PKH (Legacy)     | `1...`    | `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`                             | ✅ 完整支持 |
| P2SH (Script Hash) | `3...`    | `3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy`                             | ✅ 完整支持 |
| P2WPKH (SegWit v0) | `bc1q...` | `bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4`                     | ✅ 完整支持 |
| P2WSH (SegWit v0)  | `bc1q...` | `bc1qrp33g0q2c70qkn...`                                          | ✅ 完整支持 |
| P2TR (Taproot)     | `bc1p...` | `bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297` | ✅ 完整支持 |

## 🔐 PSBT 工作流程

### 1. 创建未签名交易

```go
psbtBase64, err := client.CreatePSBT(ctx, &types.TxInputParams{
    FromAddress:   []string{"发送地址"},
    ToAddress:     []string{"接收地址"},
    AmountBTC:     []float64{0.001},
    FeeRate:       1.0,
    ChangeAddress: "找零地址",
})
```

### 2. 外部签名 (如硬件钱包)

```go
// 使用外部钱包签名 PSBT
signedPSBT := signWithExternalWallet(psbtBase64)
```

### 3. 完成并广播

```go
txid, err := client.FinalizePSBTAndBroadcast(ctx, signedPSBT)
```

## ⚙️ 高级配置

### 自定义网络参数

```go
import "github.com/crazycloudcc/btcapis/types"

// 设置当前网络
types.SetCurrentNetwork("testnet")
```

### 费率策略

```go
// 获取推荐费率
fastRate, economyRate, err := client.EstimateFeeRate(ctx, 6)

// 使用自定义费率
txParams.FeeRate = 2.5  // sat/vB
```

## 🧪 测试示例

查看 `examples/` 目录获取完整示例：

```bash
cd examples
go run main.go
```

主要测试场景：

- **连接测试** - 验证 Bitcoin Core 和 Mempool.space 连接
- **地址操作** - 余额查询、UTXO 管理
- **交易创建** - PSBT 工作流程
- **脚本解析** - 地址和脚本转换

## 📚 文档

详细文档位于 `docs/` 目录：

- **[架构设计](docs/ARCHITECTURE.md)** - 系统架构和设计原则
- **[项目结构](docs/PROJECT_STRUCTURE.md)** - 目录组织和模块说明
- **[提供商指南](docs/PROVIDERS.md)** - 数据源接入指南
- **[btcd 兼容性](docs/COMPAT-BTCD.md)** - btcd 生态兼容策略

## 🤝 兼容性

### Go 版本支持

- **最低要求**: Go 1.23+
- **推荐版本**: Go 1.24+

### 依赖库

- **[btcsuite/btcd](https://github.com/btcsuite/btcd)** - 比特币协议实现
- **标准库** - 无其他外部依赖

### btcd 生态集成

内部使用 btcd 库进行：

- 交易编码/解码 (`wire`)
- 脚本处理 (`txscript`)
- 哈希计算 (`chainhash`)
- 地址工具 (`btcutil`)

公共 API 使用自定义类型，确保向前兼容。

## 🛠️ 开发指南

### 本地开发

```bash
# 克隆项目
git clone https://github.com/crazycloudcc/btcapis.git
cd btcapis

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 运行示例
cd examples && go run main.go
```

### 代码质量

```bash
# 代码检查
golangci-lint run

# 格式化代码
gofmt -w .

# 模块整理
go mod tidy
```

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🌟 贡献

欢迎贡献代码！请参考：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📞 支持

- **Issue 追踪**: [GitHub Issues](https://github.com/crazycloudcc/btcapis/issues)
- **讨论区**: [GitHub Discussions](https://github.com/crazycloudcc/btcapis/discussions)

## 🏷️ 版本历史

当前版本基于 Go modules `go 1.23.0`，支持：

- ✅ 完整的地址类型支持 (P2PKH, P2SH, P2WPKH, P2WSH, P2TR)
- ✅ PSBT v0 工作流程
- ✅ 多数据源架构 (Bitcoin Core + Mempool.space)
- ✅ RBF (Replace-By-Fee) 支持
- ✅ OP_RETURN 数据嵌入
- ✅ 智能费率估算

---

**BTCAPIs** - 构建现代比特币应用的可靠基础 🚀
