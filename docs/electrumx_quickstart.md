# ElectrumX 模块快速上手指南

## 目录

1. [环境准备](#环境准备)
2. [安装配置](#安装配置)
3. [基础使用](#基础使用)
4. [高级用法](#高级用法)
5. [常见问题](#常见问题)
6. [完整示例](#完整示例)

## 环境准备

### 系统要求

- Go 1.23.0 或更高版本
- 可访问的 ElectrumX 服务器

### ElectrumX 服务器选择

#### 选项 1: 使用公共服务器（快速开始）

```go
// Blockstream 公共服务器
electrumxURL := "https://blockstream.info/electrum"
```

**注意**: 公共服务器可能有速率限制，不适合生产环境。

#### 选项 2: 自建服务器（推荐生产环境）

使用 Docker 快速部署：

```bash
# 下载并运行 ElectrumX
docker run -d \
  --name electrumx \
  -p 50001:50001 \
  -p 50002:50002 \
  -e DAEMON_URL=http://user:pass@bitcoind:8332 \
  -e COIN=Bitcoin \
  -v electrumx-data:/data \
  lukechilds/electrumx
```

本地 URL:

```go
electrumxURL := "http://localhost:50001"
```

## 安装配置

### 1. 导入包

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/crazycloudcc/btcapis/internal/adapters/electrumx"
    "github.com/crazycloudcc/btcapis/internal/utils"
)
```

### 2. 创建客户端

```go
// 创建 ElectrumX 客户端
// 参数: ElectrumX服务器地址, 超时时间(秒)
client := electrumx.New("http://localhost:50001", 30)
```

### 3. 设置网络

```go
import "github.com/crazycloudcc/btcapis/types"

// 设置网络类型
types.SetCurrentNetwork("mainnet")  // 主网
// types.SetCurrentNetwork("testnet")  // 测试网
// types.SetCurrentNetwork("signet")   // 签名测试网
// types.SetCurrentNetwork("regtest")  // 回归测试网
```

## 基础使用

### 示例 1: 查询地址余额

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/crazycloudcc/btcapis/internal/adapters/electrumx"
    "github.com/crazycloudcc/btcapis/internal/utils"
    "github.com/crazycloudcc/btcapis/types"
)

func main() {
    // 设置网络
    types.SetCurrentNetwork("mainnet")

    // 创建客户端
    client := electrumx.New("http://localhost:50001", 30)
    ctx := context.Background()

    // 查询地址余额
    address := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"
    confirmed, unconfirmed, err := client.AddressGetBalance(ctx, address)
    if err != nil {
        log.Fatal(err)
    }

    // 转换为 BTC
    confirmedBTC := utils.SatsToBTC(float64(confirmed))
    unconfirmedBTC := utils.SatsToBTC(float64(unconfirmed))

    fmt.Printf("地址: %s\n", address)
    fmt.Printf("已确认: %d sats (%.8f BTC)\n", confirmed, confirmedBTC)
    fmt.Printf("未确认: %d sats (%.8f BTC)\n", unconfirmed, unconfirmedBTC)
}
```

### 示例 2: 查询 UTXO 列表

```go
func queryUTXOs(client *electrumx.Client, address string) {
    ctx := context.Background()

    // 获取 UTXO 列表
    utxos, err := client.AddressGetUTXOs(ctx, address)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("UTXO 总数: %d\n\n", len(utxos))

    var totalValue int64
    for i, utxo := range utxos {
        totalValue += utxo.Value
        valueBTC := utils.SatsToBTC(float64(utxo.Value))

        status := "已确认"
        if utxo.Height == 0 {
            status = "未确认"
        }

        fmt.Printf("%d. TxHash: %s:%d\n", i+1, utxo.TxHash, utxo.TxPos)
        fmt.Printf("   金额: %d sats (%.8f BTC)\n", utxo.Value, valueBTC)
        fmt.Printf("   高度: %d (%s)\n\n", utxo.Height, status)
    }

    totalBTC := utils.SatsToBTC(float64(totalValue))
    fmt.Printf("总价值: %d sats (%.8f BTC)\n", totalValue, totalBTC)
}
```

### 示例 3: 查询交易历史

```go
func queryHistory(client *electrumx.Client, address string) {
    ctx := context.Background()

    // 获取交易历史
    history, err := client.AddressGetHistory(ctx, address)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("交易总数: %d\n\n", len(history))

    // 显示最近 10 笔交易
    count := 10
    if len(history) < count {
        count = len(history)
    }

    for i := len(history) - count; i < len(history); i++ {
        tx := history[i]

        status := "已确认"
        confirmations := "多次确认"
        if tx.Height == 0 {
            status = "未确认"
            confirmations = "0"
        } else if tx.Height == -1 {
            status = "未广播"
            confirmations = "N/A"
        }

        fmt.Printf("%d. TXID: %s\n", i+1, tx.TxHash)
        fmt.Printf("   高度: %d\n", tx.Height)
        fmt.Printf("   状态: %s (%s)\n", status, confirmations)
        fmt.Printf("   手续费: %d sats\n\n", tx.Fee)
    }
}
```

### 示例 4: 估算交易手续费

```go
func estimateFees(client *electrumx.Client) {
    ctx := context.Background()

    // 定义不同的确认目标
    targets := map[int]string{
        1:  "高优先级（下一个区块）",
        3:  "中高优先级（约30分钟）",
        6:  "中等优先级（约1小时）",
        12: "中低优先级（约2小时）",
        24: "低优先级（约4小时）",
    }

    fmt.Println("手续费估算:")
    fmt.Println(strings.Repeat("-", 60))

    for _, blocks := range []int{1, 3, 6, 12, 24} {
        feeRate, err := client.EstimateFee(ctx, blocks)
        if err != nil {
            log.Printf("估算 %d 区块失败: %v\n", blocks, err)
            continue
        }

        // 转换为 sat/vB
        satPerVB := feeRate * 100000 // BTC/KB -> sat/vB

        fmt.Printf("%2d 区块 - %-20s: %.2f sat/vB\n",
            blocks, targets[blocks], satPerVB)
    }

    // 获取中继手续费
    relayFee, err := client.RelayFee(ctx)
    if err == nil {
        relayFeeSatPerVB := relayFee * 100000
        fmt.Printf("\n网络最小中继费率: %.2f sat/vB\n", relayFeeSatPerVB)
    }
}
```

### 示例 5: 广播交易

```go
func broadcastTransaction(client *electrumx.Client, rawTxHex string) {
    ctx := context.Background()

    // 广播交易
    txid, err := client.TransactionBroadcast(ctx, rawTxHex)
    if err != nil {
        log.Fatalf("广播交易失败: %v", err)
    }

    fmt.Printf("交易已成功广播!\n")
    fmt.Printf("TXID: %s\n", txid)
    fmt.Printf("查看交易: https://mempool.space/tx/%s\n", txid)
}
```

## 高级用法

### 并发查询多个地址

```go
import "sync"

type BalanceResult struct {
    Address    string
    Confirmed  int64
    Unconfirmed int64
    Error      error
}

func batchQueryBalances(client *electrumx.Client, addresses []string) []BalanceResult {
    ctx := context.Background()
    results := make([]BalanceResult, len(addresses))

    var wg sync.WaitGroup
    for i, addr := range addresses {
        wg.Add(1)
        go func(index int, address string) {
            defer wg.Done()

            confirmed, unconfirmed, err := client.AddressGetBalance(ctx, address)
            results[index] = BalanceResult{
                Address:    address,
                Confirmed:  confirmed,
                Unconfirmed: unconfirmed,
                Error:      err,
            }
        }(i, addr)
    }

    wg.Wait()
    return results
}

// 使用示例
func main() {
    client := electrumx.New("http://localhost:50001", 30)

    addresses := []string{
        "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
        "bc1qgnmdx4pyaxrkhtgeqgh0g93cvar7achq8kjtnm",
        // ... 更多地址
    }

    results := batchQueryBalances(client, addresses)

    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: %v\n", result.Address, result.Error)
        } else {
            fmt.Printf("✓ %s: %d sats\n", result.Address, result.Confirmed)
        }
    }
}
```

### 实现重试机制

```go
import "time"

func queryWithRetry(fn func() error, maxRetries int) error {
    var err error

    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }

        // 指数退避
        if i < maxRetries-1 {
            waitTime := time.Duration(1<<uint(i)) * time.Second
            fmt.Printf("请求失败，%v 后重试... (尝试 %d/%d)\n",
                waitTime, i+1, maxRetries)
            time.Sleep(waitTime)
        }
    }

    return fmt.Errorf("达到最大重试次数 (%d): %v", maxRetries, err)
}

// 使用示例
func main() {
    client := electrumx.New("http://localhost:50001", 30)
    address := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"

    var confirmed, unconfirmed int64
    err := queryWithRetry(func() error {
        var err error
        confirmed, unconfirmed, err = client.AddressGetBalance(
            context.Background(), address)
        return err
    }, 3)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("余额: %d sats\n", confirmed)
}
```

### 超时控制

```go
import "time"

func queryWithTimeout(client *electrumx.Client, address string) {
    // 创建带超时的 context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 查询余额
    confirmed, unconfirmed, err := client.AddressGetBalance(ctx, address)

    if err != nil {
        if err == context.DeadlineExceeded {
            fmt.Println("请求超时")
            return
        }
        log.Fatal(err)
    }

    fmt.Printf("余额: %d sats\n", confirmed)
}
```

### 监控地址变化

```go
func monitorAddress(client *electrumx.Client, address string, interval time.Duration) {
    ctx := context.Background()
    var lastBalance int64

    // 获取初始余额
    confirmed, _, err := client.AddressGetBalance(ctx, address)
    if err != nil {
        log.Fatal(err)
    }
    lastBalance = confirmed
    fmt.Printf("初始余额: %d sats\n", lastBalance)

    // 定期检查
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        confirmed, _, err := client.AddressGetBalance(ctx, address)
        if err != nil {
            log.Printf("查询失败: %v\n", err)
            continue
        }

        if confirmed != lastBalance {
            change := confirmed - lastBalance
            fmt.Printf("余额变化: %+d sats (新余额: %d sats)\n",
                change, confirmed)
            lastBalance = confirmed
        }
    }
}

// 使用示例
func main() {
    client := electrumx.New("http://localhost:50001", 30)
    address := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"

    // 每 30 秒检查一次
    monitorAddress(client, address, 30*time.Second)
}
```

## 常见问题

### Q1: 如何处理连接错误？

```go
confirmed, unconfirmed, err := client.AddressGetBalance(ctx, addr)
if err != nil {
    if strings.Contains(err.Error(), "connection refused") {
        fmt.Println("无法连接到 ElectrumX 服务器，请检查:")
        fmt.Println("1. 服务器是否正在运行")
        fmt.Println("2. 服务器地址是否正确")
        fmt.Println("3. 防火墙设置")
    } else if strings.Contains(err.Error(), "timeout") {
        fmt.Println("请求超时，请尝试:")
        fmt.Println("1. 增加超时时间")
        fmt.Println("2. 检查网络连接")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
    return
}
```

### Q2: 如何验证地址格式？

```go
import "github.com/crazycloudcc/btcapis/internal/decoders"

func validateAddress(address string) bool {
    _, err := decoders.AddressToPkScript(address)
    return err == nil
}

// 使用示例
if !validateAddress(address) {
    fmt.Println("无效的比特币地址")
    return
}
```

### Q3: 如何在不同网络间切换？

```go
import "github.com/crazycloudcc/btcapis/types"

func switchNetwork(network string, client **electrumx.Client) {
    types.SetCurrentNetwork(network)

    // 根据网络选择不同的服务器
    var electrumxURL string
    switch network {
    case "mainnet":
        electrumxURL = "http://mainnet-server:50001"
    case "testnet":
        electrumxURL = "http://testnet-server:50001"
    case "signet":
        electrumxURL = "http://signet-server:50001"
    }

    *client = electrumx.New(electrumxURL, 30)
    fmt.Printf("已切换到 %s 网络\n", network)
}
```

### Q4: 如何处理大量地址查询？

```go
func queryLargeAddressList(client *electrumx.Client, addresses []string) {
    // 控制并发数
    maxConcurrent := 10
    sem := make(chan struct{}, maxConcurrent)

    var wg sync.WaitGroup
    for _, addr := range addresses {
        wg.Add(1)
        sem <- struct{}{} // 获取信号量

        go func(address string) {
            defer wg.Done()
            defer func() { <-sem }() // 释放信号量

            ctx := context.Background()
            confirmed, _, err := client.AddressGetBalance(ctx, address)
            if err != nil {
                log.Printf("查询 %s 失败: %v\n", address, err)
                return
            }

            if confirmed > 0 {
                fmt.Printf("%s: %d sats\n", address, confirmed)
            }
        }(addr)
    }

    wg.Wait()
}
```

## 完整示例

### 完整的钱包余额查询工具

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/crazycloudcc/btcapis/internal/adapters/electrumx"
    "github.com/crazycloudcc/btcapis/internal/decoders"
    "github.com/crazycloudcc/btcapis/internal/utils"
    "github.com/crazycloudcc/btcapis/types"
)

func main() {
    // 命令行参数
    var (
        network     = flag.String("network", "mainnet", "网络类型 (mainnet/testnet/signet)")
        serverURL   = flag.String("server", "http://localhost:50001", "ElectrumX 服务器地址")
        address     = flag.String("address", "", "要查询的比特币地址")
        showUTXOs   = flag.Bool("utxos", false, "显示 UTXO 列表")
        showHistory = flag.Bool("history", false, "显示交易历史")
        timeout     = flag.Int("timeout", 30, "请求超时时间（秒）")
    )
    flag.Parse()

    // 验证参数
    if *address == "" {
        log.Fatal("请指定要查询的地址: -address <address>")
    }

    // 设置网络
    types.SetCurrentNetwork(*network)
    fmt.Printf("网络: %s\n", *network)
    fmt.Printf("服务器: %s\n", *serverURL)
    fmt.Println(strings.Repeat("-", 80))

    // 验证地址
    if _, err := decoders.AddressToPkScript(*address); err != nil {
        log.Fatalf("无效的地址: %v", err)
    }

    // 创建客户端
    client := electrumx.New(*serverURL, *timeout)
    ctx := context.Background()

    // 测试服务器连接
    fmt.Print("测试服务器连接... ")
    if err := client.ServerPing(ctx); err != nil {
        fmt.Println("❌ 失败")
        log.Fatalf("无法连接到服务器: %v", err)
    }
    fmt.Println("✓ 成功")

    // 查询余额
    fmt.Println("\n=== 地址余额 ===")
    queryBalance(client, *address)

    // 查询 UTXO
    if *showUTXOs {
        fmt.Println("\n=== UTXO 列表 ===")
        queryUTXOs(client, *address)
    }

    // 查询历史
    if *showHistory {
        fmt.Println("\n=== 交易历史 ===")
        queryHistory(client, *address)
    }

    fmt.Println("\n查询完成!")
}

func queryBalance(client *electrumx.Client, address string) {
    ctx := context.Background()

    confirmed, unconfirmed, err := client.AddressGetBalance(ctx, address)
    if err != nil {
        log.Fatalf("查询余额失败: %v", err)
    }

    confirmedBTC := utils.SatsToBTC(float64(confirmed))
    unconfirmedBTC := utils.SatsToBTC(float64(unconfirmed))
    totalBTC := confirmedBTC + unconfirmedBTC

    fmt.Printf("地址: %s\n", address)
    fmt.Printf("已确认:   %12d sats (%.8f BTC)\n", confirmed, confirmedBTC)
    fmt.Printf("未确认:   %12d sats (%.8f BTC)\n", unconfirmed, unconfirmedBTC)
    fmt.Printf("总计:     %12d sats (%.8f BTC)\n", confirmed+unconfirmed, totalBTC)
}

func queryUTXOs(client *electrumx.Client, address string) {
    ctx := context.Background()

    utxos, err := client.AddressGetUTXOs(ctx, address)
    if err != nil {
        log.Fatalf("查询 UTXO 失败: %v", err)
    }

    fmt.Printf("UTXO 总数: %d\n\n", len(utxos))

    if len(utxos) == 0 {
        fmt.Println("没有 UTXO")
        return
    }

    var totalValue int64
    for i, utxo := range utxos {
        totalValue += utxo.Value
        valueBTC := utils.SatsToBTC(float64(utxo.Value))

        status := "✓ 已确认"
        if utxo.Height == 0 {
            status = "⏳ 未确认"
        }

        fmt.Printf("%3d. %s:%d\n", i+1, utxo.TxHash, utxo.TxPos)
        fmt.Printf("     金额: %12d sats (%.8f BTC)\n", utxo.Value, valueBTC)
        fmt.Printf("     高度: %d %s\n\n", utxo.Height, status)
    }

    totalBTC := utils.SatsToBTC(float64(totalValue))
    fmt.Printf("总价值: %d sats (%.8f BTC)\n", totalValue, totalBTC)
}

func queryHistory(client *electrumx.Client, address string) {
    ctx := context.Background()

    history, err := client.AddressGetHistory(ctx, address)
    if err != nil {
        log.Fatalf("查询交易历史失败: %v", err)
    }

    fmt.Printf("交易总数: %d\n\n", len(history))

    if len(history) == 0 {
        fmt.Println("没有交易历史")
        return
    }

    // 显示最近 20 笔交易
    count := 20
    if len(history) < count {
        count = len(history)
    }

    for i := len(history) - count; i < len(history); i++ {
        tx := history[i]

        status := "✓ 已确认"
        if tx.Height == 0 {
            status = "⏳ 未确认"
        } else if tx.Height == -1 {
            status = "❌ 未广播"
        }

        fmt.Printf("%3d. %s\n", i+1, tx.TxHash)
        fmt.Printf("     高度: %d %s\n", tx.Height, status)
        fmt.Printf("     手续费: %d sats\n\n", tx.Fee)
    }
}
```

### 使用方法

```bash
# 基本查询
go run wallet_tool.go -address bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh

# 显示 UTXO
go run wallet_tool.go -address bc1q... -utxos

# 显示交易历史
go run wallet_tool.go -address bc1q... -history

# 指定网络和服务器
go run wallet_tool.go \
    -network testnet \
    -server http://testnet-server:50001 \
    -address tb1q... \
    -utxos -history

# 自定义超时
go run wallet_tool.go -address bc1q... -timeout 60
```

## 总结

本指南涵盖了 ElectrumX 模块的所有基本和高级用法。通过这些示例，您应该能够:

✅ 创建和配置 ElectrumX 客户端  
✅ 查询地址余额和 UTXO  
✅ 查询交易历史  
✅ 估算手续费  
✅ 广播交易  
✅ 实现并发查询  
✅ 处理错误和超时  
✅ 构建完整的应用程序

如需更多帮助，请参考:

- [ElectrumX README](../internal/adapters/electrumx/README.md)
- [ElectrumX 实现详解](./electrumx_implementation.md)
- [示例代码](../examples/electrumx_usage_demo.go)

---

**文档版本**: 1.0  
**最后更新**: 2025 年 11 月 2 日
