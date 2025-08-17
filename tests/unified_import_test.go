// Package tests 测试统一导入功能
package tests

import (
	"testing"

	"github.com/yourusername/btcapis"
)

// TestUnifiedImport 测试统一导入功能
func TestUnifiedImport(t *testing.T) {
	t.Run("测试BTC功能导入", func(t *testing.T) {
		// 测试地址生成
		address, err := btcapis.GenerateAddress()
		if err != nil {
			t.Fatalf("GenerateAddress失败: %v", err)
		}
		if address.Address == "" {
			t.Error("生成的地址为空")
		}

		// 测试地址验证
		if !btcapis.ValidateAddress(address.Address) {
			t.Error("地址验证失败")
		}

		// 测试私钥验证
		if !btcapis.ValidatePrivateKey(address.PrivateKey) {
			t.Error("私钥验证失败")
		}
	})

	t.Run("测试工具函数导入", func(t *testing.T) {
		// 测试随机数生成
		randomHex, err := btcapis.GenerateRandomHex(16)
		if err != nil {
			t.Fatalf("GenerateRandomHex失败: %v", err)
		}
		if len(randomHex) != 16 {
			t.Errorf("随机数长度不正确: %d", len(randomHex))
		}

		// 测试哈希计算
		hash, err := btcapis.CalculateHash([]byte("test"), btcapis.SHA256)
		if err != nil {
			t.Fatalf("CalculateHash失败: %v", err)
		}
		if hash == "" {
			t.Error("哈希计算结果为空")
		}
	})

	t.Run("测试常量导入", func(t *testing.T) {
		// 测试网络类型常量
		if btcapis.Mainnet != "mainnet" {
			t.Errorf("Mainnet常量值错误: %s", btcapis.Mainnet)
		}
		if btcapis.Testnet != "testnet" {
			t.Errorf("Testnet常量值错误: %s", btcapis.Testnet)
		}

		// 测试交易状态常量
		if btcapis.Pending != "pending" {
			t.Errorf("Pending常量值错误: %s", btcapis.Pending)
		}
		if btcapis.Confirmed != "confirmed" {
			t.Errorf("Confirmed常量值错误: %s", btcapis.Confirmed)
		}
	})

	t.Run("测试配置功能导入", func(t *testing.T) {
		// 测试默认配置加载
		config := btcapis.LoadDefaultConfig()
		if config == nil {
			t.Fatal("LoadDefaultConfig返回nil")
		}

		// 测试网络配置
		if config.Network.BTC.NetworkType == "" {
			t.Error("网络类型为空")
		}
		if config.BTC.DefaultFeeRate <= 0 {
			t.Error("默认费率不正确")
		}
	})

	t.Run("测试类型导入", func(t *testing.T) {
		// 测试BTC地址类型
		var _ btcapis.BTCAddress
		var _ btcapis.BTCTransaction
		var _ btcapis.BTCBlock

		// 测试通用类型
		var _ btcapis.Transaction
		var _ btcapis.Block
		var _ btcapis.APIResponse

		// 测试工具类型
		var _ btcapis.HashType
	})
}
