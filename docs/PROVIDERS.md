# 后端接入指南

本文档描述了如何为 btcapis 添加新的后端服务支持。

## 概述

添加新后端需要实现以下组件：

1. **客户端**: 实现`chain.Backend`接口
2. **数据映射**: 将后端响应转换为内部类型
3. **能力探测**: 动态检测后端支持的功能
4. **配置选项**: 提供灵活的配置方式

## 实现步骤

### 1. 创建后端目录

在`providers/`下创建新的目录：

```
providers/
├─ yourbackend/
│  ├─ client.go      # 主客户端
│  ├─ schema.go      # 响应模式定义
│  ├─ mapper.go      # 数据映射器
│  ├─ options.go     # 配置选项
│  └─ README.md      # 后端说明文档
```

### 2. 实现 Backend 接口

```go
// providers/yourbackend/client.go
package yourbackend

import (
    "context"
    "github.com/crazycloudcc/btcapis/chain"
    "github.com/crazycloudcc/btcapis/types"
)

type Client struct {
    // 客户端字段
}

// 实现chain.Backend接口
func (c *Client) GetRawTransaction(ctx context.Context, txid string) ([]byte, error) {
    // TODO: 实现交易查询
    return nil, nil
}

func (c *Client) GetBlockHash(ctx context.Context, height int64) (string, error) {
    // TODO: 实现区块哈希查询
    return "", nil
}

// ... 实现其他接口方法

func (c *Client) Capabilities(ctx context.Context) (*types.Capabilities, error) {
    // TODO: 实现能力探测
    return &types.Capabilities{}, nil
}

func (c *Client) Name() string {
    return "your-backend"
}

func (c *Client) IsHealthy(ctx context.Context) bool {
    // TODO: 实现健康检查
    return true
}
```

### 3. 定义响应模式

```go
// providers/yourbackend/schema.go
package yourbackend

// 定义后端原始响应结构
type TransactionResponse struct {
    ID      string `json:"id"`
    Amount  int64  `json:"amount"`
    // ... 其他字段
}

type BlockResponse struct {
    Hash    string `json:"hash"`
    Height  int64  `json:"height"`
    // ... 其他字段
}
```

### 4. 实现数据映射

```go
// providers/yourbackend/mapper.go
package yourbackend

import (
    "github.com/crazycloudcc/btcapis/types"
)

type Mapper struct{}

func (m *Mapper) MapTransaction(resp *TransactionResponse) *types.Transaction {
    if resp == nil {
        return nil
    }

    return &types.Transaction{
        TxID:    resp.ID,
        // ... 映射其他字段
    }
}

func (m *Mapper) MapBlock(resp *BlockResponse) *types.BlockHeader {
    if resp == nil {
        return nil
    }

    return &types.BlockHeader{
        Hash:   resp.Hash,
        Height: resp.Height,
        // ... 映射其他字段
    }
}
```

### 5. 配置选项

```go
// providers/yourbackend/options.go
package yourbackend

import (
    "time"
)

type Option func(*Config)

type Config struct {
    BaseURL    string
    Timeout    time.Duration
    MaxRetries int
    // ... 其他配置字段
}

func WithTimeout(timeout time.Duration) Option {
    return func(c *Config) {
        c.Timeout = timeout
    }
}

func WithMaxRetries(maxRetries int) Option {
    return func(c *Config) {
        c.MaxRetries = maxRetries
    }
}

func DefaultConfig() *Config {
    return &Config{
        Timeout:    30 * time.Second,
        MaxRetries: 3,
        // ... 设置默认值
    }
}
```

### 6. 工厂函数

```go
// providers/yourbackend/client.go

// NewClient 创建新的客户端
func NewClient(config *Config) (*Client, error) {
    if config == nil {
        config = DefaultConfig()
    }

    // 验证配置
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }

    // 创建客户端
    client := &Client{
        config: config,
        // ... 初始化其他字段
    }

    // 探测能力
    if err := client.detectCapabilities(context.Background()); err != nil {
        return nil, fmt.Errorf("failed to detect capabilities: %w", err)
    }

    return client, nil
}
```

## 能力探测

### 实现能力探测

```go
func (c *Client) detectCapabilities(ctx context.Context) error {
    capabilities := &types.Capabilities{
        Network: c.config.Network,
        // ... 设置默认能力
    }

    // 探测链上数据读取能力
    if _, err := c.GetBlockHeight(ctx); err == nil {
        capabilities.HasChainReader = true
    }

    // 探测交易广播能力
    if c.supportsBroadcasting() {
        capabilities.HasBroadcaster = true
    }

    // 探测费率估算能力
    if c.supportsFeeEstimation() {
        capabilities.HasFeeEstimator = true
    }

    // 探测内存池视图能力
    if c.supportsMempoolView() {
        capabilities.HasMempoolView = true
    }

    c.capabilities = capabilities
    return nil
}
```

### 能力测试

```go
func (c *Client) supportsBroadcasting() bool {
    // 检查后端是否支持交易广播
    // 可以通过配置、API文档或测试调用确定
    return true
}

func (c *Client) supportsFeeEstimation() bool {
    // 检查后端是否支持费率估算
    return true
}

func (c *Client) supportsMempoolView() bool {
    // 检查后端是否支持内存池查询
    return true
}
```

## 错误处理

### 错误包装

```go
import "github.com/crazycloudcc/btcapis/errors"

func (c *Client) GetRawTransaction(ctx context.Context, txid string) ([]byte, error) {
    resp, err := c.apiCall(ctx, "get_transaction", map[string]interface{}{
        "txid": txid,
    })
    if err != nil {
        return nil, errors.NewBackendError(c.Name(), err)
    }

    // ... 处理响应
    return nil, nil
}
```

### 错误分类

```go
func (c *Client) classifyError(err error) error {
    if isNotFoundError(err) {
        return errors.ErrNotFound
    }

    if isTimeoutError(err) {
        return errors.ErrTimeout
    }

    if isBackendUnavailableError(err) {
        return errors.ErrBackendUnavailable
    }

    return err
}
```

## 测试

### 单元测试

```go
// providers/yourbackend/client_test.go
package yourbackend

import (
    "testing"
    "context"
)

func TestClient_GetRawTransaction(t *testing.T) {
    client := &Client{}

    ctx := context.Background()
    _, err := client.GetRawTransaction(ctx, "test_txid")

    // 添加测试断言
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
}
```

### 集成测试

```go
// test/integration_test.go

func TestYourBackendIntegration(t *testing.T) {
    config := &Config{
        BaseURL: os.Getenv("YOUR_BACKEND_URL"),
        // ... 其他配置
    }

    client, err := NewClient(config)
    if err != nil {
        t.Fatalf("failed to create client: %v", err)
    }

    // 测试各种功能
    // ...
}
```

## 配置示例

### 环境变量

```bash
export YOUR_BACKEND_URL="https://api.yourbackend.com"
export YOUR_BACKEND_API_KEY="your_api_key"
export YOUR_BACKEND_TIMEOUT="30s"
```

### 配置文件

```yaml
backends:
  yourbackend:
    url: "https://api.yourbackend.com"
    api_key: "your_api_key"
    timeout: "30s"
    max_retries: 3
    rate_limit: 100
```

## 最佳实践

### 1. 接口实现

- 完整实现`chain.Backend`接口
- 提供有意义的错误信息
- 支持上下文取消

### 2. 性能优化

- 实现连接池
- 支持批量操作
- 添加适当的缓存

### 3. 可靠性

- 实现重试机制
- 添加超时控制
- 支持健康检查

### 4. 可观测性

- 记录详细的日志
- 提供性能指标
- 支持请求追踪

### 5. 文档

- 提供使用示例
- 说明配置选项
- 记录已知限制

## 总结

添加新后端需要：

1. 实现`chain.Backend`接口
2. 提供数据映射功能
3. 实现能力探测
4. 编写完整的测试
5. 提供配置选项
6. 遵循错误处理规范

通过遵循这些指南，可以确保新后端与现有系统良好集成，并提供一致的用户体验。
