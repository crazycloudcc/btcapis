// Package main 混合后端使用示例
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
)

func main() {
	// 创建混合后端客户端
	client := btcapis.New(
		btcapis.WithBitcoindRPC("http://127.0.0.1:8332", "user", "pass"),
		btcapis.WithMempoolSpace("https://mempool.space/api"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 获取后端能力
	capabilities, err := client.Capabilities(ctx)
	if err != nil {
		log.Fatalf("获取后端能力失败: %v", err)
	}

	fmt.Printf("后端能力: %+v\n", capabilities)

	// 并发查询示例
	results := make(chan string, 3)

	// 查询1: 费率估算
	go func() {
		fee, err := client.EstimateFeeRate(ctx, 6)
		if err != nil {
			results <- fmt.Sprintf("费率估算失败: %v", err)
			return
		}
		results <- fmt.Sprintf("6区块确认费率: %.2f sat/vB", fee)
	}()

	// 查询2: 内存池信息
	go func() {
		mempool, err := client.GetRawMempool(ctx)
		if err != nil {
			results <- fmt.Sprintf("内存池查询失败: %v", err)
			return
		}
		results <- fmt.Sprintf("内存池交易数量: %d", len(mempool))
	}()

	// 查询3: 区块信息
	go func() {
		hash, err := client.GetBlockHash(ctx, 0)
		if err != nil {
			results <- fmt.Sprintf("区块查询失败: %v", err)
			return
		}
		results <- fmt.Sprintf("区块哈希: %s", hash)
	}()

	// 收集结果
	for i := 0; i < 3; i++ {
		select {
		case result := <-results:
			fmt.Println(result)
		case <-ctx.Done():
			fmt.Println("查询超时")
			return
		}
	}
}
