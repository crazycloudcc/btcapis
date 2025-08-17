// Package chain 能力探测和特性位管理
package chain

import (
	"context"

	"github.com/crazycloudcc/btcapis/types"
)

// CapabilityDetector 定义能力探测接口
type CapabilityDetector interface {
	// Detect 探测后端能力
	Detect(ctx context.Context) (*types.Capabilities, error)

	// Validate 验证能力声明
	Validate(ctx context.Context, capabilities *types.Capabilities) error
}

// DefaultCapabilityDetector 默认能力探测器
type DefaultCapabilityDetector struct{}

// Detect 执行能力探测
func (d *DefaultCapabilityDetector) Detect(ctx context.Context) (*types.Capabilities, error) {
	// TODO: 实现默认能力探测逻辑
	return &types.Capabilities{}, nil
}

// Validate 验证能力声明
func (d *DefaultCapabilityDetector) Validate(ctx context.Context, capabilities *types.Capabilities) error {
	// TODO: 实现能力验证逻辑
	return nil
}

// CapabilityTest 定义能力测试
type CapabilityTest struct {
	Name        string
	Description string
	Test        func(ctx context.Context) error
	Required    bool
}

// NewCapabilityTest 创建新的能力测试
func NewCapabilityTest(name, description string, test func(ctx context.Context) error, required bool) *CapabilityTest {
	return &CapabilityTest{
		Name:        name,
		Description: description,
		Test:        test,
		Required:    required,
	}
}
