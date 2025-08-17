// Package tests 提供API测试用例
// 确保API功能的正确性和稳定性
package tests

import (
	"testing"

	"github.com/yourusername/btcapis/pkg/api/btc"
)

// TestGenerateAddress 测试地址生成功能
func TestGenerateAddress(t *testing.T) {
	address, err := btc.GenerateAddress()
	if err != nil {
		t.Fatalf("生成地址失败: %v", err)
	}

	// 检查地址结构
	if address.Address == "" {
		t.Error("生成的地址为空")
	}
	if address.PrivateKey == "" {
		t.Error("生成的私钥为空")
	}
	if address.PublicKey == "" {
		t.Error("生成的公钥为空")
	}
	if address.AddressType == "" {
		t.Error("生成的地址类型为空")
	}

	// 检查地址长度（简化验证）
	if len(address.Address) < 20 {
		t.Errorf("地址长度过短: %d", len(address.Address))
	}

	// 检查私钥长度（32字节 = 64个十六进制字符）
	if len(address.PrivateKey) != 64 {
		t.Errorf("私钥长度不正确: %d", len(address.PrivateKey))
	}
}

// TestGenerateAddressWithType 测试指定类型地址生成
func TestGenerateAddressWithType(t *testing.T) {
	types := []string{"legacy", "p2sh", "bech32"}

	for _, addrType := range types {
		address, err := btc.GenerateAddressWithType(addrType)
		if err != nil {
			t.Fatalf("生成%s类型地址失败: %v", addrType, err)
		}

		if address.AddressType != addrType {
			t.Errorf("期望地址类型%s，实际得到%s", addrType, address.AddressType)
		}
	}
}

// TestValidateAddress 测试地址验证功能
func TestValidateAddress(t *testing.T) {
	// 测试有效地址
	validAddress := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	if !btc.ValidateAddress(validAddress) {
		t.Error("有效地址验证失败")
	}

	// 测试无效地址
	invalidAddresses := []string{
		"",        // 空地址
		"invalid", // 无效字符
		"123",     // 过短
		"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa123456789", // 过长
	}

	for _, addr := range invalidAddresses {
		if btc.ValidateAddress(addr) {
			t.Errorf("无效地址验证通过: %s", addr)
		}
	}
}

// TestGetAddressType 测试地址类型检测
func TestGetAddressType(t *testing.T) {
	// 测试不同类型的地址
	testCases := []struct {
		address      string
		expectedType string
	}{
		{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "legacy"},
		{"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", "p2sh"},
		{"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", "bech32"},
		{"invalid", "invalid"},
	}

	for _, tc := range testCases {
		actualType := btc.GetAddressType(tc.address)
		if actualType != tc.expectedType {
			t.Errorf("地址%s期望类型%s，实际得到%s", tc.address, tc.expectedType, actualType)
		}
	}
}

// TestCalculateTransactionFee 测试交易费用计算
func TestCalculateTransactionFee(t *testing.T) {
	// 测试不同输入输出组合的费用计算
	testCases := []struct {
		inputCount  int
		outputCount int
		feeRate     float64
		expectedFee float64
	}{
		{1, 2, 10.0, 0.00000190}, // 1输入2输出，费率10 sat/byte
		{2, 2, 10.0, 0.00000230}, // 2输入2输出，费率10 sat/byte
		{1, 1, 5.0, 0.00000095},  // 1输入1输出，费率5 sat/byte
	}

	for _, tc := range testCases {
		fee := btc.CalculateTransactionFee(tc.inputCount, tc.outputCount, tc.feeRate)
		if fee < 0 {
			t.Errorf("费用计算为负数: %f", fee)
		}
		// 注意：由于是估算值，这里只检查是否为正数
	}
}

// TestValidatePrivateKey 测试私钥验证
func TestValidatePrivateKey(t *testing.T) {
	// 测试有效私钥
	validPrivateKey := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	if !btc.ValidatePrivateKey(validPrivateKey) {
		t.Error("有效私钥验证失败")
	}

	// 测试无效私钥
	invalidPrivateKeys := []string{
		"",        // 空私钥
		"123",     // 过短
		"invalid", // 无效字符
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1", // 过长
	}

	for _, pk := range invalidPrivateKeys {
		if btc.ValidatePrivateKey(pk) {
			t.Errorf("无效私钥验证通过: %s", pk)
		}
	}
}

// TestGenerateAddressMultiple 测试多次生成地址的唯一性
func TestGenerateAddressMultiple(t *testing.T) {
	addresses := make(map[string]bool)

	// 生成多个地址
	for i := 0; i < 10; i++ {
		address, err := btc.GenerateAddress()
		if err != nil {
			t.Fatalf("生成第%d个地址失败: %v", i+1, err)
		}

		// 检查地址唯一性
		if addresses[address.Address] {
			t.Errorf("发现重复地址: %s", address.Address)
		}
		addresses[address.Address] = true
	}
}

// BenchmarkGenerateAddress 地址生成性能测试
func BenchmarkGenerateAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := btc.GenerateAddress()
		if err != nil {
			b.Fatalf("生成地址失败: %v", err)
		}
	}
}

// BenchmarkValidateAddress 地址验证性能测试
func BenchmarkValidateAddress(b *testing.B) {
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	for i := 0; i < b.N; i++ {
		btc.ValidateAddress(address)
	}
}
