# BTC APIs

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

BTC APIs 是一个功能完整的比特币区块链 API 统一接口库，采用端口/适配器/门面架构，支持多种后端服务，提供统一的比特币区块链操作接口。

## ✨ 主要特性

- 🔄 **多后端聚合**：支持 Bitcoin Core RPC、mempool.space、Electrum 等
- 🚀 **智能路由**：自动选择最佳后端，支持故障转移和负载均衡
- 📊 **完整交易处理**：交易解析、验证、广播、费率估算
- 🏠 **地址支持**：Base58、Bech32、Taproot 等多种地址类型
- 📜 **脚本处理**：脚本分类、构建、反汇编、签名哈希计算
- 🔐 **PSBT 支持**：完整的 PSBT 构建、签名、最终化流程
- 🏗️ **模块化架构**：高度可扩展，易于集成新后端和功能
- 🧪 **全面测试**：包含单元测试和集成测试

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        门面层 (Facade)                      │
│                    btcapis.go, *_facade.go                  │
├─────────────────────────────────────────────────────────────┤
│                        端口层 (Ports)                       │
│                    chain/backend.go                         │
├─────────────────────────────────────────────────────────────┤
│                        适配器层 (Adapters)                  │
│              providers/bitcoindrpc, mempoolspace           │
└─────────────────────────────────────────────────────────────┘
```

### 核心模块

- **chain**: 路由器和后端接口定义
- **providers**: 具体后端实现（Bitcoin Core、mempool.space 等）
- **tx**: 交易解析和处理
- **address**: 地址编码和验证
- **script**: 脚本分析和构建
- **psbt**: PSBT 处理
- **types**: 核心数据类型定义

## 🚀 快速开始

### 安装

```bash
go get github.com/crazycloudcc/btcapis
```

### 基础使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/crazycloudcc/btcapis"
)

func main() {
    // 创建客户端，配置多个后端
    c := btcapis.New(
        btcapis.WithBitcoindRPC(
            os.Getenv("BITCOIND_URL"),
            os.Getenv("BITCOIND_USER"),
            os.Getenv("BITCOIND_PASS"),
        ),
        btcapis.WithMempoolSpace("https://mempool.space"),
    )

    // 获取交易信息
    tx, err := c.GetTransaction(context.Background(), "your-txid-here")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("交易ID: %s\n", tx.TxID)
    fmt.Printf("输入数量: %d\n", len(tx.Vin))
    fmt.Printf("输出数量: %d\n", len(tx.Vout))
    fmt.Printf("虚拟大小: %d\n", tx.Vsize)
}
```

### 环境配置

创建 `.env` 文件：

```bash
# Bitcoin Core RPC
BITCOIND_URL=http://localhost:8332
BITCOIND_USER=your_username
BITCOIND_PASS=your_password

# mempool.space (可选)
MEMPOOLSPACE_URL=https://mempool.space
```

## 📚 API 参考

### 交易操作

```go
// 获取交易
tx, err := client.GetTransaction(ctx, txid)

// 获取原始交易数据
rawTx, err := client.GetRawTransaction(ctx, txid)

// 广播交易
txid, err := client.Broadcast(ctx, rawTx)

// 估算费率
feeRate, err := client.EstimateFeeRate(ctx, targetBlocks)
```

### 地址操作

```go
// 地址分类
addrType := address.Classify(scriptPubKey)

// Base58 编码/解码
encoded := address.Base58Encode(data)
decoded := address.Base58Decode(encoded)

// Bech32 编码/解码
encoded := address.Bech32Encode(hrp, data)
decoded := address.Bech32Decode(encoded)
```

### 脚本操作

```go
// 脚本分类
scriptType, addresses := script.Classify(pkScript)

// 脚本反汇编
asm := script.Disasm(pkScript)

// 构建脚本
pkScript := script.Builder{}.
    AddOp(script.OP_DUP).
    AddOp(script.OP_HASH160).
    AddData(hash160).
    AddOp(script.OP_EQUALVERIFY).
    AddOp(script.OP_CHECKSIG).
    Build()
```

### PSBT 操作

```go
// 创建 PSBT
psbt := psbt.New()

// 添加输入
psbt.AddInput(prevTx, vout, scriptPubKey, amount)

// 添加输出
psbt.AddOutput(scriptPubKey, amount)

// 签名
psbt.SignInput(inputIndex, privateKey, sighashType)

// 最终化
finalTx := psbt.Finalize()
```

## 🔧 配置选项

### Bitcoin Core RPC 选项

```go
client := btcapis.New(
    btcapis.WithBitcoindRPC(
        "http://localhost:8332",
        "username",
        "password",
        bitcoindrpc.WithHTTPClient(customHTTPClient),
        bitcoindrpc.WithTimeout(10*time.Second),
    ),
)
```

### mempool.space 选项

```go
client := btcapis.New(
    btcapis.WithMempoolSpace(
        "https://mempool.space",
        mempoolspace.WithHTTPClient(customHTTPClient),
        mempoolspace.WithTimeout(8*time.Second),
    ),
)
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./chain/...
go test ./providers/...

# 运行集成测试
go test ./test/...
```

### 测试覆盖率

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📁 项目结构

```
btcapis/
├── chain/           # 路由器和后端接口
├── providers/       # 后端实现
│   ├── bitcoindrpc/ # Bitcoin Core RPC
│   ├── mempoolspace/ # mempool.space API
│   └── electrum/    # Electrum 协议
├── tx/              # 交易处理
├── address/         # 地址处理
├── script/          # 脚本处理
├── psbt/            # PSBT 处理
├── types/           # 类型定义
├── internal/        # 内部工具
├── examples/        # 使用示例
└── docs/            # 文档
```

## 🔌 扩展后端

实现新的后端服务：

```go
type CustomBackend struct {
    // 实现 chain.Backend 接口
}

func (b *CustomBackend) GetTransaction(ctx context.Context, txid string) (*types.Tx, error) {
    // 实现具体逻辑
}

func (b *CustomBackend) Capabilities(ctx context.Context) (chain.Capabilities, error) {
    return chain.Capabilities{
        HasMempool:     true,
        HasFeeEstimate: false,
        Network:        types.Mainnet,
    }, nil
}

// 添加到客户端
client := btcapis.New(
    btcapis.WithCustomBackend(&CustomBackend{}),
)
```

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [btcd](https://github.com/btcsuite/btcd) - 比特币协议实现
- [btcutil](https://github.com/btcsuite/btcutil) - 比特币工具库
- [mempool.space](https://mempool.space) - 内存池数据服务

## 📞 支持

如有问题或建议，请：

- 提交 [Issue](https://github.com/crazycloudcc/btcapis/issues)
- 查看 [文档](docs/)
- 参考 [示例](examples/)

---

**BTC APIs** - 让比特币区块链开发更简单 🚀
