# BTC APIs - Go 比特币 API 合集

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个专为 Go 项目设计的比特币 API 合集，提供完整的比特币区块链功能。**一次导入，所有功能立即可用！**

## 🚀 核心特色

- **🎯 一次导入**: `import "github.com/yourusername/btcapis"` 即可使用所有功能
- **专注 BTC**: 专门为比特币区块链设计，功能完整
- **模块化设计**: 清晰的包结构，易于集成和维护
- **类型安全**: 完整的 Go 类型定义，编译时错误检查
- **详细注释**: 每个函数和类型都有详细的中文注释
- **测试覆盖**: 包含完整的测试用例和性能基准测试
- **配置灵活**: 支持配置文件和环境变量配置
- **多地址类型**: 支持 Legacy、P2SH、Bech32 等地址类型

## 📁 项目结构

```
btcapis/
├── btcapis.go           # 🎯 统一入口文件（一次导入所有功能）
├── pkg/api/             # 核心API包
│   ├── btc/            # 比特币相关API
│   ├── common/         # 通用类型和结构
│   ├── utils/          # 工具函数
│   └── config/         # 配置管理
├── examples/            # 使用示例
├── tests/              # 测试用例
├── config/             # 配置文件
├── go.mod              # Go模块文件
├── Makefile            # 开发工具命令
├── .gitignore          # Git忽略文件
└── README.md           # 项目说明
```

## 🛠️ 安装

### 前置要求

- Go 1.21 或更高版本
- Git

### 安装步骤

1. 克隆项目

```bash
git clone https://github.com/yourusername/btcapis.git
cd btcapis
```

2. 初始化 Go 模块

```bash
go mod tidy
```

3. 运行测试

```bash
go test ./...
```

## 📖 使用方法

### 🎯 超简单导入

```go
import "github.com/yourusername/btcapis"

// 现在可以直接使用所有功能！
```

### 完整使用示例

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/btcapis"
)

func main() {
    // 生成BTC地址
    address, err := btcapis.GenerateAddress()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("地址: %s\n", address.Address)
    fmt.Printf("私钥: %s\n", address.PrivateKey)

    // 验证地址
    if btcapis.ValidateAddress(address.Address) {
        fmt.Println("✅ 地址验证通过")
    }

    // 计算交易费用
    fee := btcapis.CalculateTransactionFee(2, 2, 10.0)
    fmt.Printf("交易费用: %.8f BTC\n", fee)

    // 生成随机数
    randomHex, _ := btcapis.GenerateRandomHex(32)
    fmt.Printf("随机数: %s\n", randomHex)

    // 计算哈希
    hash, _ := btcapis.CalculateHash([]byte("Hello Bitcoin"), btcapis.SHA256)
    fmt.Printf("哈希值: %s\n", hash)

    // 创建响应
    resp := btcapis.NewSuccessResponse(address)
    fmt.Printf("响应: %+v\n", resp)
}
```

### 🎯 所有可用功能

#### BTC 核心功能

- `btcapis.GenerateAddress()` - 生成比特币地址
- `btcapis.GenerateAddressWithType(type)` - 生成指定类型地址
- `btcapis.ValidateAddress(address)` - 验证地址格式
- `btcapis.GetAddressType(address)` - 获取地址类型
- `btcapis.CalculateTransactionFee(inputs, outputs, rate)` - 计算交易费用
- `btcapis.ValidatePrivateKey(key)` - 验证私钥

#### 工具函数

- `btcapis.GenerateRandomHex(length)` - 生成随机十六进制
- `btcapis.CalculateHash(data, type)` - 计算哈希值
- `btcapis.ValidateHexString(str)` - 验证十六进制字符串
- `btcapis.HexToBytes(hex)` - 十六进制转字节
- `btcapis.BytesToHex(bytes)` - 字节转十六进制

#### 通用类型和常量

- `btcapis.BTCAddress` - 比特币地址结构
- `btcapis.Transaction` - 交易结构
- `btcapis.Block` - 区块结构
- `btcapis.Mainnet/Testnet/Regtest` - 网络类型常量
- `btcapis.Pending/Confirmed/Failed` - 交易状态常量

#### 配置管理

- `btcapis.LoadDefaultConfig()` - 加载默认配置
- `btcapis.LoadConfig(path)` - 从文件加载配置

## ⚙️ 配置

### 环境变量

| 变量名                 | 描述                    | 默认值                            |
| ---------------------- | ----------------------- | --------------------------------- |
| `BTC_MAINNET_RPC`      | BTC 主网 RPC 地址       | `https://btc.getblock.io/mainnet` |
| `BTC_TESTNET_RPC`      | BTC 测试网 RPC 地址     | `https://btc.getblock.io/testnet` |
| `BTC_REGTEST_RPC`      | BTC 回归测试网 RPC 地址 | `http://localhost:18443`          |
| `BTC_NETWORK_TYPE`     | BTC 网络类型            | `mainnet`                         |
| `BTC_DEFAULT_FEE_RATE` | 默认费率(sat/byte)      | `10.0`                            |
| `BTC_CONFIRMATIONS`    | 默认确认数              | `6`                               |
| `BTC_ADDRESS_TYPE`     | 默认地址类型            | `legacy`                          |
| `API_PORT`             | API 服务端口            | `8080`                            |
| `API_HOST`             | API 服务主机            | `0.0.0.0`                         |
| `LOG_LEVEL`            | 日志级别                | `info`                            |

### 配置文件

项目支持 JSON 格式的配置文件，详见 `config/config.json` 示例。

## 🧪 测试

### 运行所有测试

```bash
go test ./...
```

### 运行特定包的测试

```bash
go test ./pkg/api/btc
go test ./pkg/api/utils
```

### 运行性能基准测试

```bash
go test -bench=. ./...
```

### 生成测试覆盖率报告

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 开发指南

### 添加新的 BTC 功能

1. 在 `pkg/api/btc/` 下添加新的功能函数
2. 在 `btcapis.go` 中重新导出新功能
3. 实现核心功能逻辑
4. 添加详细的注释
5. 编写测试用例
6. 更新 README 文档

### 代码规范

- 使用中文注释说明功能
- 遵循 Go 官方代码规范
- 所有导出的函数和类型都要有注释
- 错误处理要详细和友好
- 专注于 BTC 相关功能
- 保持统一入口的简洁性

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 贡献步骤

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🆘 支持

如果您在使用过程中遇到问题，请：

1. 查看 [Issues](https://github.com/yourusername/btcapis/issues) 页面
2. 创建新的 Issue 描述问题
3. 联系维护者

## 🔮 路线图

- [ ] 添加更多 BTC 地址类型支持
- [ ] 实现交易签名和广播功能
- [ ] 添加 UTXO 管理功能
- [ ] 实现多重签名支持
- [ ] 添加闪电网络集成
- [ ] 提供 RESTful API 接口
- [ ] 添加监控和日志功能
- [ ] 支持更多 BTC 网络（如 Signet）

## 📊 项目状态

- **开发状态**: 活跃开发中
- **测试覆盖率**: >90%
- **Go 版本支持**: 1.21+
- **最后更新**: 2024 年 12 月
- **专注领域**: 比特币区块链
- **导入方式**: 一次导入，所有功能立即可用

## 🎯 核心功能

### 地址管理

- 生成多种类型的比特币地址
- 地址格式验证
- 地址类型检测
- 私钥验证

### 交易处理

- 交易费用计算
- 交易结构定义
- 交易状态管理

### 区块信息

- 区块结构定义
- 区块信息查询
- 交易确认管理

### 工具函数

- 随机数生成
- 哈希计算
- 十六进制处理
- 数据验证

## 🚀 快速开始

```go
package main

import (
    "fmt"
    "github.com/yourusername/btcapis"
)

func main() {
    // 一行导入，所有功能立即可用！
    address, _ := btcapis.GenerateAddress()
    fmt.Printf("BTC地址: %s\n", address.Address)
}
```

---

**🎉 特色**: 通过 `import "github.com/yourusername/btcapis"` 一次导入，即可使用所有比特币相关功能！

**注意**: 本项目仅用于学习和开发目的，在生产环境中使用前请充分测试。
