# btcapis

一个高性能的比特币 API 库，采用端口/适配器/门面架构，支持多种后端服务。

## 架构特点

- **三层架构**: 端口(ports) + 适配器(adapters) + 门面(facade)
- **多后端支持**: Bitcoin Core RPC、mempool.space、Electrum 等
- **智能路由**: 故障转移、负载均衡、并发查询
- **统一接口**: 一致的 API 设计，后端透明切换

## 快速开始

```go
import "github.com/yourusername/btcapis"

func main() {
    client := btcapis.New(
        btcapis.WithBitcoindRPC("http://127.0.0.1:8332", "user", "pass"),
        btcapis.WithMempoolSpace("https://mempool.space/api"),
    )

    // 地址解析（纯计算，无需后端）
    info, err := btcapis.Address.Parse("bc1q...", btcapis.Mainnet)

    // 费率估算（自动降级）
    fee, err := client.EstimateFeeRate(context.Background(), 6)
}
```

## 目录结构

```
btcapis/
├─ btcapis.go                 # 根包门面
├─ address_facade.go          # 地址模块门面
├─ script_facade.go           # 脚本模块门面
├─ tx_facade.go               # 交易模块门面
├─ psbt_facade.go             # PSBT模块门面
├─ chain_facade.go            # 链上数据门面
├─ errors.go                  # 错误定义
├─ types/                     # 共享类型定义
├─ address/                   # 地址编解码实现
├─ script/                    # 脚本处理实现
├─ tx/                        # 交易处理实现
├─ psbt/                      # PSBT工具实现
├─ chain/                     # 后端接口定义
├─ providers/                 # 后端适配器实现
├─ internal/                  # 内部基础设施
├─ examples/                  # 使用示例
├─ test/                      # 测试文件
└─ docs/                      # 文档
```

## 特性

- 🚀 **高性能**: 并发查询、智能缓存、连接池
- 🔄 **高可用**: 自动故障转移、降级策略
- 🛡️ **可靠性**: 重试机制、超时控制、错误处理
- 🔌 **可扩展**: 插件化后端支持
- 📚 **易使用**: 简洁的 API 设计

## 文档

- [架构设计](docs/ARCHITECTURE.md)
- [后端接入指南](docs/PROVIDERS.md)
- [btcd 兼容性](docs/COMPAT-BTCD.md)

## 许可证

MIT License
