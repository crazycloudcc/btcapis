# BTCAPIs 项目目录结构

## 项目概述

BTCAPIs 是一个比特币相关的 Go 语言项目，提供地址处理、交易操作、PSBT 管理、脚本处理等功能。

## 完整目录结构

```
btcapis/
├── .gitignore                    # Git 忽略文件配置
├── .golangci.yml                # Go 代码质量检查配置
├── LICENSE                       # 项目许可证
├── README.md                     # 项目说明文档
├── go.mod                        # Go 模块依赖文件
├── go.sum                        # Go 模块校验文件
│
├── btcapis.go                    # 主程序入口文件
├── btcapis_address.go            # 地址相关操作
├── btcapis_addressops.go         # 地址操作接口
├── btcapis_txops.go              # 交易操作接口
│
├── address/                       # 地址处理模块
│   ├── base58.go                 # Base58 编码/解码
│   ├── bech32.go                 # Bech32 编码/解码
│   ├── classify.go                # 地址分类
│   │── common.go                   # 通用功能
│   │── parse2script.go            # 地址解析
│   ├── scriptpubkey.go           # 脚本公钥处理
│   └── testdata/                 # 测试数据目录
│
├── chain/                         # 区块链相关模块
│   ├── address.go                # 链上地址处理
│   ├── backend.go                # 后端接口
│   ├── capabilities.go            # 功能能力定义
│   ├── errors.go                 # 错误处理
│   └── router.go                 # 路由处理
│
├── docs/                          # 项目文档目录
│   ├── 1.init需求.txt            # 初始化需求文档
│   ├── ARCHITECTURE.md           # 架构设计文档
│   ├── COMPAT-BTCD.md            # BTCD 兼容性文档
│   ├── PROVIDERS.md              # 提供商文档
│   ├── psbt优化1.txt             # PSBT 优化文档1
│   ├── psbt优化2.txt             # PSBT 优化文档2
│   ├── tx1.txt                   # 交易文档1
│   ├── tx2.txt                   # 交易文档2
│   ├── 测试交易数据结构-ord数据.txt  # 测试数据结构文档
│   └── PROJECT_STRUCTURE.md      # 项目结构文档（本文件）
│
├── examples/                      # 示例代码目录
│   ├── basic/                    # 基础示例
│   │   └── main.go
│   ├── bitcoind_fee/             # Bitcoin Core 费用示例
│   │   └── main.go
│   ├── mempool_mix/              # 内存池混合示例
│   │   └── main.go
│   └── txdump/                   # 交易转储示例
│       └── main.go
│
├── internal/                      # 内部工具模块
│   ├── assert/                   # 断言工具
│   │   └── assert.go
│   ├── bytespool/                # 字节池管理
│   │   └── pool.go
│   ├── cache/                    # 缓存管理
│   │   └── cache.go
│   ├── httpx/                    # HTTP 扩展工具
│   │   ├── client.go             # HTTP 客户端
│   │   ├── limiter.go            # 限流器
│   │   ├── middleware.go         # 中间件
│   │   └── retry.go              # 重试机制
│   ├── jsonrpc/                  # JSON-RPC 客户端
│   │   └── client.go
│   └── trace/                    # 追踪工具
│       └── trace.go
│
├── providers/                     # 数据提供商模块
│   ├── bitcoindrpc/              # Bitcoin Core RPC 提供商
│   │   ├── client.go             # RPC 客户端
│   │   ├── mapper.go             # 数据映射器
│   │   ├── methods.go            # RPC 方法
│   │   └── options.go            # 配置选项
│   ├── blockstream/              # Blockstream 提供商
│   │   └── client.go
│   ├── electrum/                 # Electrum 提供商
│   │   └── client.go
│   └── mempoolspace/             # Mempool.space 提供商
│       ├── address.go            # 地址相关 API
│       ├── client.go             # 客户端实现
│       ├── mapper.go             # 数据映射器
│       └── schema.go             # 数据模式定义
│
├── psbt/                          # PSBT (Partially Signed Bitcoin Transaction) 模块
│   ├── analyze.go                # PSBT 分析
│   ├── finalize.go               # PSBT 完成
│   ├── psbt.go                   # PSBT 核心功能
│   ├── psbt_test.go              # PSBT 测试
│   └── roles.go                  # PSBT 角色定义
│
├── script/                        # 脚本处理模块
│   ├── asm.go                    # 脚本汇编
│   ├── builder.go                # 脚本构建器
│   ├── classify.go                # 脚本分类
│   ├── decompile.go              # 脚本反编译
│   ├── disasm.go                 # 脚本反汇编
│   ├── opcodes.go                # 操作码定义
│   ├── ordinal.go                # Ordinal 相关脚本
│   ├── sighash.go                # 签名哈希
│   ├── taproot.go                # Taproot 脚本
│   └── taproot_hash.go           # Taproot 哈希
│
├── test/                          # 测试目录
│   └── integration_test.go       # 集成测试
│
├── tx/                            # 交易处理模块
│   ├── decoder.go                # 交易解码器
│   ├── inspect.go                # 交易检查器
│   ├── inspect_spk.go            # 脚本公钥检查器
│   └── protocol.go               # 交易协议
│
├── types/                         # 类型定义模块
│   ├── address.go                # 地址类型
│   ├── brc20.go                  # BRC-20 代币类型
│   ├── json_hex.go               # JSON 十六进制类型
│   ├── network.go                # 网络类型
│   ├── ordinal.go                # Ordinal 类型
│   ├── runes.go                  # Runes 类型
│   ├── scriptview.go             # 脚本视图类型
│   ├── taproot.go                # Taproot 类型
│   ├── tx.go                     # 交易类型
│   └── utxo.go                   # UTXO 类型
│
└── x/                             # 扩展模块
    └── tapleaf/                   # Tapleaf 扩展
```

## 模块功能说明

### 核心模块

- **address**: 比特币地址处理，支持 Base58 和 Bech32 编码
- **chain**: 区块链交互，提供后端接口和路由功能
- **psbt**: PSBT 交易处理，支持分析、完成和角色管理
- **script**: 比特币脚本处理，包括汇编、反汇编、分类等
- **tx**: 交易处理，提供解码、检查和协议支持
- **types**: 核心数据类型定义，涵盖地址、交易、UTXO 等

### 提供商模块

- **bitcoindrpc**: Bitcoin Core RPC 接口
- **blockstream**: Blockstream API 接口
- **btcdapis**: BTCD APIs 接口
- **electrum**: Electrum 协议接口
- **mempoolspace**: Mempool.space API 接口

### 工具模块

- **internal**: 内部工具，包括 HTTP 客户端、缓存、断言等
- **examples**: 使用示例，涵盖各种应用场景

## 项目特点

1. **模块化设计**: 清晰的模块分离，便于维护和扩展
2. **多提供商支持**: 支持多种比特币数据源
3. **完整功能覆盖**: 涵盖地址、交易、脚本、PSBT 等核心功能
4. **测试完备**: 包含单元测试和集成测试
5. **文档齐全**: 提供详细的架构和使用文档

## 使用建议

- 新用户可从 `examples/` 目录开始学习
- 开发时参考 `docs/` 目录中的架构文档
- 核心功能在 `types/` 和 `tx/` 模块中
- 扩展功能可通过 `providers/` 模块接入不同数据源
