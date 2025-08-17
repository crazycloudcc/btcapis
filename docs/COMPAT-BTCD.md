# btcd 兼容性指南

本文档描述了 btcapis 与 btcd 生态系统的兼容性策略和最佳实践。

## 概述

btcapis 设计为与 btcd 生态系统协作，但不强耦合。我们推荐使用 btcd 的底层库，但在公共 API 中避免直接暴露 btcd 类型。

## 兼容性策略

### 1. 推荐使用 btcd 库

btcapis 在内部使用以下 btcd 库：

```go
import (
    "github.com/btcsuite/btcd/chaincfg/chainhash"
    "github.com/btcsuite/btcd/txscript"
    "github.com/btcsuite/btcd/wire"
    "github.com/btcsuite/btcd/btcutil"
)
```

### 2. 类型隔离

在公共 API 中使用自己的类型定义，避免直接暴露 btcd 类型：

```go
// ✅ 推荐：使用内部类型
type Transaction struct {
    TxID   string
    // ... 其他字段
}

// ❌ 避免：直接暴露btcd类型
type Transaction struct {
    *wire.MsgTx  // 不要这样做
}
```

### 3. 转换函数

提供转换函数在内部类型和 btcd 类型之间转换：

```go
// 内部类型转btcd类型
func ToBtcdTx(tx *types.Transaction) *wire.MsgTx {
    // 实现转换逻辑
}

// btcd类型转内部类型
func FromBtcdTx(msgTx *wire.MsgTx) *types.Transaction {
    // 实现转换逻辑
}
```

## 具体实现

### 地址处理

```go
// types/address.go
type AddressInfo struct {
    Address     string
    Network     Network
    Type        AddressType
    ScriptPubKey []byte
}

// address/bech32.go
import (
    "github.com/btcsuite/btcd/bech32"
    "github.com/btcsuite/btcd/chaincfg"
)

func ParseBech32(addr string, network Network) (*AddressInfo, error) {
    // 使用btcd的bech32库解析地址
    hrp, data, err := bech32.Decode(addr)
    if err != nil {
        return nil, err
    }

    // 转换为内部类型
    return &AddressInfo{
        Address: addr,
        Network: network,
        Type:    AddressTypeP2WPKH,
        // ... 其他字段
    }, nil
}
```

### 脚本处理

```go
// types/script.go
type ScriptInfo struct {
    Type      ScriptType
    Hex       string
    ASM       string
    Addresses []string
}

// script/decompile.go
import (
    "github.com/btcsuite/btcd/txscript"
)

func DecompileScript(script []byte) (*ScriptInfo, error) {
    // 使用btcd的txscript库反编译脚本
    tokens, err := txscript.DisasmString(script)
    if err != nil {
        return nil, err
    }

    // 转换为内部类型
    return &ScriptInfo{
        Hex: hex.EncodeToString(script),
        ASM: tokens,
        // ... 其他字段
    }, nil
}
```

### 交易处理

```go
// types/tx.go
type Transaction struct {
    TxID     string
    Version  int32
    Inputs   []TxInput
    Outputs  []TxOutput
}

// tx/tx.go
import (
    "github.com/btcsuite/btcd/wire"
)

func ParseTransaction(rawTx []byte) (*Transaction, error) {
    // 使用btcd的wire库解析交易
    msgTx := wire.NewMsgTx(wire.TxVersion)
    err := msgTx.Deserialize(bytes.NewReader(rawTx))
    if err != nil {
        return nil, err
    }

    // 转换为内部类型
    return FromBtcdTx(msgTx), nil
}

func SerializeTransaction(tx *Transaction) ([]byte, error) {
    // 转换为btcd类型
    msgTx := ToBtcdTx(tx)

    // 序列化
    var buf bytes.Buffer
    err := msgTx.Serialize(&buf)
    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}
```

### 哈希处理

```go
// types/hash.go
type Hash [32]byte

// 转换函数
func FromChainHash(chainHash chainhash.Hash) Hash {
    var hash Hash
    copy(hash[:], chainHash[:])
    return hash
}

func ToChainHash(hash Hash) chainhash.Hash {
    return chainhash.Hash(hash)
}
```

## 版本兼容性

### Go 模块版本

```go
// go.mod
require (
    github.com/btcsuite/btcd v0.24.0
    github.com/btcsuite/btcd/btcutil v1.1.5
    github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2
    github.com/btcsuite/btcd/txscript/v4 v4.0.0
    github.com/btcsuite/btcd/wire v0.0.0-20231230182452-3c6c6c2fbf91
)
```

### 版本约束

- 使用语义化版本控制
- 支持最新的稳定版本
- 避免使用预发布版本

## 迁移策略

### 从 btcd 直接使用迁移

如果用户之前直接使用 btcd 类型：

```go
// 之前：直接使用btcd类型
import "github.com/btcsuite/btcd/wire"

func processTx(msgTx *wire.MsgTx) {
    // 处理交易
}

// 现在：使用btcapis类型
import "github.com/crazycloudcc/btcapis/types"

func processTx(tx *types.Transaction) {
    // 处理交易
}

// 如果需要btcd类型，使用转换函数
msgTx := btcapis.Tx.ToBtcdTx(tx)
```

### 渐进式迁移

1. **第一阶段**: 使用 btcapis 的公共 API
2. **第二阶段**: 逐步替换内部 btcd 类型使用
3. **第三阶段**: 完全迁移到 btcapis 类型系统

## 性能考虑

### 类型转换开销

类型转换会带来一定的性能开销，但通常可以忽略：

```go
// 性能测试示例
func BenchmarkTypeConversion(b *testing.B) {
    tx := createTestTransaction()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        msgTx := ToBtcdTx(tx)
        _ = FromBtcdTx(msgTx)
    }
}
```

### 优化建议

1. **缓存转换结果**: 避免重复转换
2. **批量转换**: 一次转换多个对象
3. **延迟转换**: 只在必要时进行转换

## 测试策略

### 兼容性测试

```go
// test/btcd_compat_test.go
func TestBtcdCompatibility(t *testing.T) {
    // 测试btcd类型转换
    msgTx := createBtcdMsgTx()

    // 转换为内部类型
    tx := FromBtcdTx(msgTx)

    // 再转换回btcd类型
    msgTx2 := ToBtcdTx(tx)

    // 验证一致性
    if !reflect.DeepEqual(msgTx, msgTx2) {
        t.Error("类型转换不一致")
    }
}
```

### 基准测试

```go
func BenchmarkBtcdCompatibility(b *testing.B) {
    msgTx := createBtcdMsgTx()

    b.Run("ToInternal", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = FromBtcdTx(msgTx)
        }
    })

    b.Run("ToBtcd", func(b *testing.B) {
        tx := FromBtcdMsgTx(msgTx)
        for i := 0; i < b.N; i++ {
            _ = ToBtcdTx(tx)
        }
    })
}
```

## 最佳实践

### 1. 类型使用

- 在公共 API 中使用内部类型
- 在内部实现中使用 btcd 类型
- 提供清晰的转换函数

### 2. 错误处理

- 保持错误语义一致
- 包装 btcd 错误为内部错误
- 提供有意义的错误信息

### 3. 性能优化

- 避免不必要的类型转换
- 使用对象池减少内存分配
- 实现高效的序列化/反序列化

### 4. 文档

- 明确说明兼容性策略
- 提供迁移指南
- 记录已知限制

## 总结

btcapis 与 btcd 生态系统的兼容性策略：

1. **推荐使用**: 在内部使用 btcd 库
2. **类型隔离**: 公共 API 使用内部类型
3. **转换函数**: 提供类型转换功能
4. **渐进迁移**: 支持平滑迁移路径
5. **性能优化**: 最小化转换开销

通过这种策略，btcapis 可以：

- 利用 btcd 生态系统的成熟性
- 保持 API 的稳定性
- 支持未来的技术演进
- 提供良好的用户体验
