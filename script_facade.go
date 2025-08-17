// Package btcapis 脚本模块门面
package btcapis

import (
	"github.com/crazycloudcc/btcapis/types"
)

// scriptFacade 提供脚本相关的功能接口
type scriptFacade struct{}

// Build 构建脚本
func (s *scriptFacade) Build(scriptType types.ScriptType, params map[string]interface{}) ([]byte, error) {
	// TODO: 实现脚本构建
	return nil, nil
}

// Decompile 反编译脚本
func (s *scriptFacade) Decompile(script []byte) (*types.ScriptInfo, error) {
	// TODO: 实现脚本反编译
	return nil, nil
}

// Classify 分类脚本类型
func (s *scriptFacade) Classify(script []byte) (types.ScriptType, error) {
	// TODO: 实现脚本分类
	return types.ScriptTypeUnknown, nil
}

// Validate 验证脚本格式
func (s *scriptFacade) Validate(script []byte) error {
	// TODO: 实现脚本验证
	return nil
}
