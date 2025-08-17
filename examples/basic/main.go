// Package main 基础使用示例
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
)

func main() {
	// 创建客户端
	client := btcapis.New(
		btcapis.WithBitcoindRPC("http://127.0.0.1:8332", "user", "pass"),
		btcapis.WithMempoolSpace("https://mempool.space/api"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 地址解析（纯计算，无需后端）
	info, err := btcapis.Address.Parse("bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", btcapis.Mainnet)
	if err != nil {
		log.Fatalf("地址解析失败: %v", err)
	}
	fmt.Printf("地址信息: %+v\n", info)

	// 费率估算（自动降级）
	fee, err := client.EstimateFeeRate(ctx, 6)
	if err != nil {
		log.Fatalf("费率估算失败: %v", err)
	}
	fmt.Printf("6区块确认费率: %.2f sat/vB\n", fee)

	// 获取当前区块高度
	height, err := client.GetBlockHash(ctx, 0) // 临时使用GetBlockHash作为示例
	if err != nil {
		log.Fatalf("获取区块哈希失败: %v", err)
	}
	fmt.Printf("区块哈希: %s\n", height)

	// 获取内存池信息
	mempool, err := client.GetRawMempool(ctx)
	if err != nil {
		log.Fatalf("获取内存池失败: %v", err)
	}
	fmt.Printf("内存池交易数量: %d\n", len(mempool))
}
