// Package main Bitcoin Core费率估算示例
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/btcapis"
)

func main() {
	// 创建仅使用Bitcoin Core的客户端
	client := btcapis.New(
		btcapis.WithBitcoindRPC("http://127.0.0.1:8332", "user", "pass"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 不同确认时间的费率估算
	targets := []int{1, 3, 6, 10, 25, 144}

	for _, target := range targets {
		fee, err := client.EstimateFeeRate(ctx, target)
		if err != nil {
			log.Printf("估算%d区块确认费率失败: %v", target, err)
			continue
		}
		fmt.Printf("%d区块确认费率: %.2f sat/vB\n", target, fee)
	}

	// 获取内存池信息
	mempool, err := client.GetRawMempool(ctx)
	if err != nil {
		log.Fatalf("获取内存池失败: %v", err)
	}
	fmt.Printf("内存池交易数量: %d\n", len(mempool))
}
