// Package psbt PSBT处理
package psbt

import (
	"fmt"

	"github.com/crazycloudcc/btcapis/types"
)

// CreatePSBT 创建新的PSBT
func CreatePSBT() (*types.PSBT, error) {
	// TODO: 实现PSBT创建
	return &types.PSBT{
		Version: 0,
		Inputs:  []types.PSBTInput{},
		Outputs: []types.PSBTOutput{},
		Unknown: []types.PSBTUnknown{},
	}, nil
}

// ParsePSBT 解析PSBT
func ParsePSBT(rawPSBT []byte) (*types.PSBT, error) {
	if len(rawPSBT) == 0 {
		return nil, fmt.Errorf("raw PSBT cannot be empty")
	}

	// TODO: 实现PSBT解析
	// 使用BIP174规范进行解析

	return nil, fmt.Errorf("not implemented")
}

// SerializePSBT 序列化PSBT
func SerializePSBT(psbt *types.PSBT) ([]byte, error) {
	if psbt == nil {
		return nil, fmt.Errorf("PSBT cannot be nil")
	}

	// TODO: 实现PSBT序列化
	// 使用BIP174规范进行序列化

	return nil, fmt.Errorf("not implemented")
}

// AddInput 添加输入
func AddInput(psbt *types.PSBT, input types.PSBTInput) error {
	if psbt == nil {
		return fmt.Errorf("PSBT cannot be nil")
	}

	psbt.Inputs = append(psbt.Inputs, input)
	return nil
}

// AddOutput 添加输出
func AddOutput(psbt *types.PSBT, output types.PSBTOutput) error {
	if psbt == nil {
		return fmt.Errorf("PSBT cannot be nil")
	}

	psbt.Outputs = append(psbt.Outputs, output)
	return nil
}

// FinalizePSBT 完成PSBT
func FinalizePSBT(psbt *types.PSBT) (*types.Transaction, error) {
	if psbt == nil {
		return nil, fmt.Errorf("PSBT cannot be nil")
	}

	// TODO: 实现PSBT完成
	// 验证所有输入都已签名
	// 构建最终交易

	return nil, fmt.Errorf("not implemented")
}

// ValidatePSBT 验证PSBT
func ValidatePSBT(psbt *types.PSBT) error {
	if psbt == nil {
		return fmt.Errorf("PSBT cannot be nil")
	}

	// TODO: 实现PSBT验证
	// 检查基本格式、输入输出数量等

	return nil
}
