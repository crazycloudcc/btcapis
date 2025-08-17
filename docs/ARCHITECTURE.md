# 架构设计

本文档描述了 btcapis 的架构设计，包括端口/适配器/门面模式、线程安全性和错误模型。

## 架构概览

btcapis 采用三层架构设计：

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

## 门面层 (Facade)

门面层提供统一的 API 接口，隐藏内部复杂性：

- **btcapis.go**: 主门面，提供 Client 和配置选项
- **address_facade.go**: 地址相关功能
- **script_facade.go**: 脚本相关功能
- **tx_facade.go**: 交易相关功能
- **psbt_facade.go**: PSBT 相关功能
- **chain_facade.go**: 链上数据查询功能

### 设计原则

1. **简单性**: 提供简洁的 API 接口
2. **一致性**: 所有模块使用相同的设计模式
3. **可组合性**: 支持多种后端组合使用

## 端口层 (Ports)

端口层定义抽象接口，不依赖具体实现：

### chain.Backend 接口

```go
type Backend interface {
    ChainReader
    Broadcaster
    FeeEstimator
    MempoolView

    Capabilities(ctx context.Context) (Capabilities, error)
    Name() string
    IsHealthy(ctx context.Context) bool
}
```

### 能力探测

通过`Capabilities()`方法动态探测后端能力：

```go
type Capabilities struct {
    HasChainReader    bool
    HasBroadcaster    bool
    HasFeeEstimator   bool
    HasMempoolView    bool
    Network           types.Network
    SupportsSegWit    bool
    SupportsTaproot   bool
    // ... 其他能力
}
```

## 适配器层 (Adapters)

适配器层实现具体的后端服务：

### Bitcoin Core RPC

- **providers/bitcoindrpc/client.go**: RPC 客户端
- **providers/bitcoindrpc/methods.go**: RPC 方法封装
- **providers/bitcoindrpc/mapper.go**: 数据映射
- **providers/bitcoindrpc/options.go**: 配置选项

### mempool.space

- **providers/mempoolspace/client.go**: REST 客户端
- **providers/mempoolspace/schema.go**: 响应模式
- **providers/mempoolspace/mapper.go**: 数据映射

## 路由策略

### 智能路由

路由器根据操作类型和后端能力选择最佳后端：

1. **首选策略**: 优先使用支持该操作的后端
2. **降级策略**: 首选失败时自动降级到其他后端
3. **并发策略**: 对查询操作可并发执行，选择最快响应

### 故障转移

- 健康检查: 定期检查后端状态
- 自动切换: 检测到故障时自动切换到健康后端
- 重试机制: 支持指数退避重试

## 线程安全性

### 并发模型

- **只读操作**: 支持并发访问
- **写操作**: 串行化执行，避免竞态条件
- **状态更新**: 使用互斥锁保护共享状态

### 资源管理

- **连接池**: HTTP 客户端复用连接
- **速率限制**: 防止后端过载
- **超时控制**: 避免长时间阻塞

## 错误模型

### 错误类型

1. **业务错误**: 如资源不存在、参数无效
2. **网络错误**: 如连接超时、后端不可用
3. **系统错误**: 如内存不足、配置错误

### 错误处理

```go
// 使用errors.Is检查错误类型
if errors.Is(err, types.ErrNotFound) {
    // 处理资源不存在
}

// 使用errors.As获取具体错误信息
var backendErr *BackendError
if errors.As(err, &backendErr) {
    log.Printf("后端 %s 错误: %v", backendErr.Backend, backendErr.Err)
}
```

## 性能优化

### 缓存策略

- **静态数据**: 区块哈希、网络信息等
- **动态数据**: 费率、内存池状态等（短 TTL）
- **分层缓存**: 内存 L1 + 持久化 L2

### 并发优化

- **连接复用**: 避免频繁建立连接
- **批量操作**: 支持批量查询减少网络开销
- **异步处理**: 非阻塞的并发查询

## 扩展性

### 新后端接入

1. 实现`chain.Backend`接口
2. 添加能力探测
3. 实现数据映射
4. 编写测试用例

### 新功能模块

1. 在门面层添加新接口
2. 在端口层定义抽象
3. 在适配器层实现具体逻辑

## 监控与可观测性

### 指标收集

- 请求延迟
- 成功率
- 后端健康状态
- 缓存命中率

### 日志记录

- 结构化日志
- 可配置的日志级别
- 请求追踪 ID

## 总结

btcapis 的架构设计遵循以下原则：

1. **关注点分离**: 每层职责明确，便于维护
2. **接口抽象**: 通过接口实现解耦
3. **可测试性**: 每层都可以独立测试
4. **可扩展性**: 易于添加新功能和后端
5. **性能优先**: 支持并发、缓存、连接复用等优化
