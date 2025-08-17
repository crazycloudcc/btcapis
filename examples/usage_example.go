// Package examples 提供API使用示例
// 展示如何在其他项目中使用这个BTC API合集
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/btcapis"
)

func main() {
	fmt.Println("=== BTC API 使用示例（简化导入） ===")

	// 生成比特币地址
	btcAddress, err := btcapis.GenerateAddress()
	if err != nil {
		log.Fatalf("生成BTC地址失败: %v", err)
	}

	fmt.Printf("生成的BTC地址: %s\n", btcAddress.Address)
	fmt.Printf("私钥: %s\n", btcAddress.PrivateKey)
	fmt.Printf("公钥: %s\n", btcAddress.PublicKey)
	fmt.Printf("地址类型: %s\n", btcAddress.AddressType)

	// 验证BTC地址
	isValid := btcapis.ValidateAddress(btcAddress.Address)
	fmt.Printf("地址验证结果: %t\n", isValid)

	// 获取地址类型
	addressType := btcapis.GetAddressType(btcAddress.Address)
	fmt.Printf("地址类型检测: %s\n", addressType)

	// 生成不同类型的地址
	legacyAddress, err := btcapis.GenerateAddressWithType("legacy")
	if err != nil {
		log.Fatalf("生成Legacy地址失败: %v", err)
	}
	fmt.Printf("Legacy地址: %s\n", legacyAddress.Address)

	// 计算交易费用
	fee := btcapis.CalculateTransactionFee(2, 2, 10.0) // 2输入2输出，费率10 sat/byte
	fmt.Printf("估算交易费用: %.8f BTC\n", fee)

	// 验证私钥
	isValidPrivateKey := btcapis.ValidatePrivateKey(btcAddress.PrivateKey)
	fmt.Printf("私钥验证结果: %t\n", isValidPrivateKey)

	fmt.Println("\n=== 通用工具函数示例 ===")

	// 生成随机十六进制字符串
	randomHex, err := btcapis.GenerateRandomHex(32)
	if err != nil {
		log.Fatalf("生成随机十六进制失败: %v", err)
	}
	fmt.Printf("随机十六进制: %s\n", randomHex)

	// 计算哈希值
	hash, err := btcapis.CalculateHash([]byte("Hello, Bitcoin!"), btcapis.SHA256)
	if err != nil {
		log.Fatalf("计算哈希失败: %v", err)
	}
	fmt.Printf("SHA256哈希: %s\n", hash)

	// 验证十六进制字符串
	isValidHex := btcapis.ValidateHexString(randomHex)
	fmt.Printf("十六进制验证结果: %t\n", isValidHex)

	fmt.Println("\n=== 通用类型示例 ===")

	// 创建成功响应
	successResp := btcapis.NewSuccessResponse(map[string]interface{}{
		"message": "操作成功",
		"data":    btcAddress,
	})
	fmt.Printf("成功响应: %+v\n", successResp)

	// 创建错误响应
	errorResp := btcapis.NewErrorResponse("操作失败", 500)
	fmt.Printf("错误响应: %+v\n", errorResp)

	// 创建BTC交易结构
	transaction := &btcapis.Transaction{
		TxHash:      "0x1234567890abcdef",
		From:        btcAddress.Address,
		To:          "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		Amount:      "0.001",
		Fee:         "0.0001",
		Status:      btcapis.Pending,
		BlockNumber: 12345,
		Timestamp:   time.Now(),
		Data:        []byte("Hello Bitcoin"),
	}
	fmt.Printf("交易信息: %+v\n", transaction)

	fmt.Println("\n=== 配置管理示例 ===")

	// 加载默认配置
	defaultConfig := btcapis.LoadDefaultConfig()
	fmt.Printf("默认网络类型: %s\n", defaultConfig.Network.BTC.NetworkType)
	fmt.Printf("默认RPC地址: %s\n", defaultConfig.GetCurrentNetworkRPC())
	fmt.Printf("默认费率: %.1f sat/byte\n", defaultConfig.BTC.DefaultFeeRate)
	fmt.Printf("默认确认数: %d\n", defaultConfig.BTC.Confirmations)

	// 获取特定网络RPC地址
	testnetRPC := defaultConfig.GetBTCNetworkRPC("testnet")
	fmt.Printf("测试网RPC: %s\n", testnetRPC)

	fmt.Println("\n=== 网络类型常量示例 ===")
	fmt.Printf("主网: %s\n", btcapis.Mainnet)
	fmt.Printf("测试网: %s\n", btcapis.Testnet)
	fmt.Printf("回归测试网: %s\n", btcapis.Regtest)

	fmt.Println("\n=== 交易状态常量示例 ===")
	fmt.Printf("待确认: %s\n", btcapis.Pending)
	fmt.Printf("已确认: %s\n", btcapis.Confirmed)
	fmt.Printf("失败: %s\n", btcapis.Failed)
}
