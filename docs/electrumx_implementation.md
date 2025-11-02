# ElectrumX 实现说明文档

## 概述

本文档详细说明了 ElectrumX 适配器的实现细节，包括设计理念、代码结构、与其他适配器的对比，以及使用指南。

## 设计理念

### 1. 遵循项目规范

ElectrumX 适配器的实现严格遵循 `bitcoindrpc` 和 `mempoolapis` 的代码规范和设计模式：

- **客户端结构**: 使用相同的客户端初始化模式
- **错误处理**: 统一的错误处理和日志记录机制
- **上下文管理**: 所有 API 都支持 context.Context 进行超时控制
- **DTO 设计**: 清晰的数据传输对象定义

### 2. 代码组织

```
electrumx/
├── client.go     # 客户端核心实现和 RPC 调用
├── apis.go       # API 接口实现
├── dtos.go       # 数据传输对象定义
└── README.md     # 使用文档
```

### 3. 模块化设计

- **client.go**: 提供基础的 JSON-RPC 调用能力
- **apis.go**: 实现具体的业务逻辑和 API 封装
- **dtos.go**: 定义所有的数据结构

## 核心实现

### 1. 客户端初始化

```go
type Client struct {
    url    string       // ElectrumX 服务器地址
    http   *http.Client // HTTP 客户端
    idSeed int          // RPC 请求 ID 种子
}

func New(baseURL string, timeout int) *Client {
    return &Client{
        url:  baseURL,
        http: &http.Client{Timeout: time.Duration(timeout) * time.Second},
    }
}
```

**设计要点:**

- 简洁的构造函数
- 可配置的超时时间
- 自增的请求 ID 管理

### 2. JSON-RPC 调用

```go
func (c *Client) rpcCall(ctx context.Context, method string, params []interface{}, out interface{}) error
```

**实现特点:**

- 完整的错误处理和日志记录
- 支持 context 取消和超时
- 类型安全的结果解析
- 与 Bitcoin Core RPC 类似的实现模式

### 3. 地址到脚本哈希转换

ElectrumX 使用脚本哈希而非地址进行索引，这是一个关键的实现细节：

```go
func addressToScriptHash(addr string) (string, error) {
    // 1. 使用 decoders 模块将地址转换为 scriptPubKey
    pkScript, err := decoders.AddressToPkScript(addr)
    if err != nil {
        return "", fmt.Errorf("address to pkscript: %w", err)
    }

    // 2. 计算 SHA256 哈希并反序
    scriptHash := computeScriptHash(pkScript)
    return scriptHash, nil
}

func computeScriptHash(script []byte) string {
    // 计算 SHA256 哈希
    hash := sha256.Sum256(script)

    // 反序字节序（ElectrumX 使用 little-endian）
    result := hash[:]
    for i := 0; i < len(result)/2; i++ {
        result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
    }

    return hex.EncodeToString(result)
}
```

**关键点:**

- 复用 `decoders` 模块，避免重复实现
- 正确处理字节序转换（little-endian）
- 支持所有主流地址类型（P2PKH, P2SH, P2WPKH, P2WSH, P2TR）

## API 实现详解

### 1. 地址相关 API

#### AddressGetBalance

```go
func (c *Client) AddressGetBalance(ctx context.Context, addr string) (int64, int64, error)
```

**实现流程:**

1. 将地址转换为脚本哈希
2. 调用 `blockchain.scripthash.get_balance`
3. 返回已确认和未确认余额

**特点:**

- 直接返回 satoshi 单位，避免精度问题
- 同时返回已确认和未确认余额

#### AddressGetHistory

```go
func (c *Client) AddressGetHistory(ctx context.Context, addr string) ([]HistoryDTO, error)
```

**返回数据:**

- 完整的交易历史
- 包含区块高度和手续费信息
- 按时间顺序排列

#### AddressGetUTXOs

```go
func (c *Client) AddressGetUTXOs(ctx context.Context, addr string) ([]UTXODTO, error)
```

**特点:**

- 返回所有未花费输出
- 包含确认状态
- 可直接用于构建交易

### 2. 交易相关 API

#### TransactionBroadcast

```go
func (c *Client) TransactionBroadcast(ctx context.Context, rawTxHex string) (string, error)
```

**用途:**

- 广播签名后的交易
- 返回交易 ID
- 支持所有交易类型（包括 SegWit 和 Taproot）

#### TransactionGet

```go
func (c *Client) TransactionGet(ctx context.Context, txid string, verbose bool) (interface{}, error)
```

**灵活性:**

- `verbose=false`: 返回原始十六进制
- `verbose=true`: 返回详细的交易信息
- 支持已确认和未确认的交易

### 3. 手续费估算 API

#### EstimateFee

```go
func (c *Client) EstimateFee(ctx context.Context, blocks int) (float64, error)
```

**返回值:**

- 单位: BTC/KB
- 可转换为 sat/vB: `satPerVB = feeRate * 100000`

**使用建议:**

```go
// 估算下一个区块的手续费
feeRate, err := client.EstimateFee(ctx, 1)
if err != nil {
    log.Fatal(err)
}

// 转换为 sat/vB
satPerVB := feeRate * 100000
fmt.Printf("建议手续费率: %.2f sat/vB\n", satPerVB)
```

### 4. 服务器管理 API

实现了完整的服务器管理接口：

- `ServerVersion`: 获取版本信息
- `ServerFeatures`: 获取服务器功能
- `ServerPing`: 心跳检测
- `ServerBanner`: 获取欢迎信息

## 与其他适配器对比

### 1. 与 Bitcoin Core RPC 对比

| 特性          | ElectrumX                 | Bitcoin Core RPC      |
| ------------- | ------------------------- | --------------------- |
| **部署**      | 需要独立 ElectrumX 服务器 | 需要完整 Bitcoin 节点 |
| **资源占用**  | 低（索引数据库）          | 高（完整区块链）      |
| **查询速度**  | 快（优化的索引）          | 慢（需要扫描）        |
| **地址查询**  | 原生支持                  | 需要 `-addressindex`  |
| **UTXO 查询** | 即时                      | 需要全链扫描或索引    |
| **交易历史**  | 完整支持                  | 有限支持              |
| **钱包功能**  | 无                        | 完整支持              |
| **交易构建**  | 需要客户端实现            | 内置支持              |

**使用场景:**

- **ElectrumX**: 轻量级钱包、查询服务、区块浏览器
- **Bitcoin Core**: 矿池、完整节点、需要钱包功能

### 2. 与 Mempool.space 对比

| 特性           | ElectrumX           | Mempool.space          |
| -------------- | ------------------- | ---------------------- |
| **协议**       | JSON-RPC            | REST API               |
| **部署**       | 自建服务器          | 使用公共 API           |
| **稳定性**     | 依赖自建服务器      | 依赖第三方服务         |
| **功能丰富度** | 标准 ElectrumX 协议 | 扩展功能（图表、统计） |
| **隐私性**     | 高（自建）          | 低（第三方）           |
| **实时性**     | 高                  | 高                     |
| **成本**       | 服务器成本          | 免费（有限额）         |

**使用场景:**

- **ElectrumX**: 需要隐私、高可用、大量查询
- **Mempool.space**: 快速开发、小规模应用

### 3. 实现方式对比

#### Bitcoin Core RPC 模式

```go
// 认证方式
client := bitcoindrpc.New(url, user, pass, timeout)

// RPC 调用
err := client.rpcCall(ctx, method, params, &result)
```

#### Mempool.space 模式

```go
// REST API
u := *c.base
u.Path = path.Join(u.Path, "/api/address/", addr)
err := c.getJSON(ctx, u.String(), &result)
```

#### ElectrumX 模式

```go
// JSON-RPC over HTTP
client := electrumx.New(url, timeout)
err := client.rpcCall(ctx, method, params, &result)
```

## 使用 decoders 模块

### 复用现有功能

ElectrumX 实现充分利用了 `decoders` 模块的功能：

```go
import "github.com/crazycloudcc/btcapis/internal/decoders"

// 1. 地址到 scriptPubKey
pkScript, err := decoders.AddressToPkScript(addr)

// 2. 地址类型识别
addrType, err := decoders.AddressToType(addr)

// 3. 完整地址解析
info, err := decoders.DecodeAddress(addr)
```

### 支持的地址类型

| 类型   | 前缀 | 示例       | 支持状态 |
| ------ | ---- | ---------- | -------- |
| P2PKH  | 1    | 1A1zP1...  | ✅       |
| P2SH   | 3    | 3J98t1...  | ✅       |
| P2WPKH | bc1q | bc1qxy2... | ✅       |
| P2WSH  | bc1q | bc1q7cy... | ✅       |
| P2TR   | bc1p | bc1p5c...  | ✅       |

## 错误处理

### 1. 统一的错误格式

```go
if rpcResp.Error != nil {
    return fmt.Errorf("electrumx rpc error %d: %s",
        rpcResp.Error.Code, rpcResp.Error.Message)
}
```

### 2. 常见错误类型

| 错误代码 | 含义       | 处理建议             |
| -------- | ---------- | -------------------- |
| -1       | 未知错误   | 检查参数和服务器状态 |
| -32600   | 无效请求   | 检查 JSON 格式       |
| -32601   | 方法不存在 | 检查方法名           |
| -32602   | 无效参数   | 检查参数类型和数量   |
| -32603   | 内部错误   | 检查服务器日志       |

### 3. 错误处理示例

```go
balance, err := client.AddressGetBalance(ctx, addr)
if err != nil {
    // 检查是否是超时错误
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("请求超时，请稍后重试")
        return
    }

    // 检查是否是地址格式错误
    if strings.Contains(err.Error(), "address to pkscript") {
        log.Println("地址格式错误")
        return
    }

    // 其他错误
    log.Printf("查询失败: %v", err)
    return
}
```

## 性能优化

### 1. 批量查询

对于需要查询多个地址的场景，建议使用 goroutine 并发查询：

```go
type Result struct {
    Addr      string
    Confirmed int64
    Error     error
}

func batchQueryBalance(client *electrumx.Client, ctx context.Context, addrs []string) []Result {
    results := make([]Result, len(addrs))
    var wg sync.WaitGroup

    for i, addr := range addrs {
        wg.Add(1)
        go func(idx int, address string) {
            defer wg.Done()
            confirmed, _, err := client.AddressGetBalance(ctx, address)
            results[idx] = Result{
                Addr:      address,
                Confirmed: confirmed,
                Error:     err,
            }
        }(i, addr)
    }

    wg.Wait()
    return results
}
```

### 2. 缓存策略

对于不经常变化的数据，建议实现缓存：

```go
type CachedClient struct {
    client *electrumx.Client
    cache  *sync.Map
}

func (c *CachedClient) GetBlockHeader(ctx context.Context, height int64) (string, error) {
    // 区块头不会改变，可以永久缓存
    if cached, ok := c.cache.Load(height); ok {
        return cached.(string), nil
    }

    header, err := c.client.BlockchainGetBlockHeader(ctx, height, 0)
    if err != nil {
        return "", err
    }

    c.cache.Store(height, header)
    return header, nil
}
```

## 测试指南

### 1. 单元测试

参考 `examples/tests/test_electrumx.go` 中的测试用例：

```go
func TestElectrumX() {
    client := electrumx.New("http://localhost:50001", 30)
    ctx := context.Background()

    // 测试服务器连接
    err := client.ServerPing(ctx)
    assert.NoError(t, err)

    // 测试余额查询
    confirmed, unconfirmed, err := client.AddressGetBalance(ctx, testAddr)
    assert.NoError(t, err)
    assert.GreaterOrEqual(t, confirmed, int64(0))
}
```

### 2. 集成测试

```bash
# 启动本地 ElectrumX 服务器（Docker）
docker run -d \
  -p 50001:50001 \
  -v ~/.bitcoin:/root/.bitcoin \
  lukechilds/electrumx

# 运行测试
cd examples/tests
go run tests.go
```

## 部署建议

### 1. 公共 ElectrumX 服务器

```go
// Blockstream
client := electrumx.New("https://blockstream.info/electrum", 30)

// 其他公共服务器
// 注意: 公共服务器可能有速率限制
```

### 2. 自建 ElectrumX 服务器

**优点:**

- 无速率限制
- 更好的隐私性
- 完全控制

**要求:**

- Bitcoin Core 节点
- 足够的磁盘空间（索引数据库）
- 稳定的网络连接

**推荐配置:**

```bash
# 使用 Docker Compose
version: '3'
services:
  electrumx:
    image: lukechilds/electrumx
    ports:
      - "50001:50001"
      - "50002:50002"
    environment:
      DAEMON_URL: http://user:pass@bitcoind:8332
      COIN: Bitcoin
    volumes:
      - electrumx-data:/data
```

## 最佳实践

### 1. 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

balance, err := client.AddressGetBalance(ctx, addr)
```

### 2. 重试机制

```go
func retryQuery(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }

        if i < maxRetries-1 {
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }
    return fmt.Errorf("max retries exceeded")
}
```

### 3. 日志记录

代码中已经包含了详细的日志记录：

- 请求开始/结束
- 错误信息
- 响应时间

可以通过日志级别控制输出：

```go
log.SetLevel(log.DebugLevel) // 详细日志
log.SetLevel(log.ErrorLevel) // 仅错误
```

## 常见问题

### Q1: 为什么使用脚本哈希而不是地址？

**A:** ElectrumX 协议设计使用脚本哈希是为了：

- 统一索引不同类型的地址
- 提高查询效率
- 简化数据库设计

### Q2: 如何处理大量地址查询？

**A:** 使用并发查询和连接池：

```go
// 并发控制
sem := make(chan struct{}, 10) // 最多10个并发
for _, addr := range addresses {
    sem <- struct{}{}
    go func(a string) {
        defer func() { <-sem }()
        // 查询逻辑
    }(addr)
}
```

### Q3: 支持哪些网络？

**A:** 支持所有 Bitcoin 网络：

- 主网 (mainnet)
- 测试网 (testnet)
- 回归测试网 (regtest)

通过 `types.CurrentNetwork` 配置。

### Q4: 如何确保查询结果的准确性？

**A:**

1. 使用可信的 ElectrumX 服务器
2. 对关键数据进行多节点验证
3. 检查确认数
4. 验证交易和区块数据

## 未来扩展

### 计划中的功能

1. **WebSocket 支持**: 实现实时订阅功能
2. **批量 RPC**: 支持批量请求提高效率
3. **连接池**: 管理多个连接
4. **自动重连**: 网络故障自动恢复
5. **指标监控**: 请求统计和性能监控

### 扩展示例

```go
// WebSocket 订阅（未来版本）
subscription, err := client.AddressSubscribeWS(ctx, addr, func(status string) {
    fmt.Printf("地址状态更新: %s\n", status)
})

// 批量查询（未来版本）
balances, err := client.BatchGetBalance(ctx, addresses)
```

## 贡献指南

欢迎提交 Issue 和 Pull Request！

### 提交代码前请确保：

1. 遵循现有代码风格
2. 添加适当的注释
3. 包含测试用例
4. 更新相关文档

## 参考资源

- [ElectrumX 协议文档](https://electrumx-spesmilo.readthedocs.io/)
- [Bitcoin Core RPC 文档](https://developer.bitcoin.org/reference/rpc/)
- [btcd 库文档](https://pkg.go.dev/github.com/btcsuite/btcd)

## 许可证

遵循项目主许可证
