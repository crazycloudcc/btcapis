//go:build integration
// +build integration

// Package test 集成测试
package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/crazycloudcc/btcapis"
)

// 测试配置
var (
	bitcoindURL  = os.Getenv("BITCOIND_URL")
	bitcoindUser = os.Getenv("BITCOIND_USER")
	bitcoindPass = os.Getenv("BITCOIND_PASS")
	mempoolURL   = os.Getenv("MEMPOOL_BASE_URL")
)

func TestMain(m *testing.M) {
	// 检查环境变量
	if bitcoindURL == "" || bitcoindUser == "" || bitcoindPass == "" {
		println("跳过集成测试: 缺少BITCOIND环境变量")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestBitcoindRPCIntegration(t *testing.T) {
	client := btcapis.New(
		btcapis.WithBitcoindRPC(bitcoindURL, bitcoindUser, bitcoindPass),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试费率估算
	fee, err := client.EstimateFeeRate(ctx, 6)
	if err != nil {
		t.Fatalf("费率估算失败: %v", err)
	}

	if fee <= 0 {
		t.Errorf("费率应该大于0，实际: %f", fee)
	}

	t.Logf("6区块确认费率: %.2f sat/vB", fee)
}

func TestMempoolSpaceIntegration(t *testing.T) {
	if mempoolURL == "" {
		t.Skip("跳过mempool.space测试: 缺少MEMPOOL_BASE_URL环境变量")
	}

	client := btcapis.New(
		btcapis.WithMempoolSpace(mempoolURL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试内存池查询
	mempool, err := client.GetRawMempool(ctx)
	if err != nil {
		t.Fatalf("内存池查询失败: %v", err)
	}

	t.Logf("内存池交易数量: %d", len(mempool))
}

func TestMixedBackendsIntegration(t *testing.T) {
	client := btcapis.New(
		btcapis.WithBitcoindRPC(bitcoindURL, bitcoindUser, bitcoindPass),
	)

	if mempoolURL != "" {
		// 如果配置了mempool.space，添加到客户端
		// TODO: 实现WithMempoolSpace选项
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 测试并发查询
	results := make(chan error, 2)

	// 查询1: 费率估算
	go func() {
		_, err := client.EstimateFeeRate(ctx, 6)
		results <- err
	}()

	// 查询2: 内存池信息
	go func() {
		_, err := client.GetRawMempool(ctx)
		results <- err
	}()

	// 收集结果
	for i := 0; i < 2; i++ {
		select {
		case err := <-results:
			if err != nil {
				t.Errorf("并发查询失败: %v", err)
			}
		case <-ctx.Done():
			t.Fatal("并发查询超时")
		}
	}
}
