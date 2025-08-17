// Package btcapis 交易模块门面
package btcapis

import (
	"github.com/crazycloudcc/btcapis/types"
)

// txFacade 提供交易相关的功能接口
type txFacade struct{}

// Parse 解析原始交易
func (t *txFacade) Parse(rawTx []byte) (*types.Transaction, error) {
	// TODO: 实现交易解析
	return nil, nil
}

// Serialize 序列化交易
func (t *txFacade) Serialize(tx *types.Transaction) ([]byte, error) {
	// TODO: 实现交易序列化
	return nil, nil
}

// CalculateTxID 计算交易ID
func (t *txFacade) CalculateTxID(rawTx []byte) (string, error) {
	// TODO: 实现交易ID计算
	return "", nil
}

// Validate 验证交易格式
func (t *txFacade) Validate(rawTx []byte) error {
	// TODO: 实现交易验证
	return nil
}

// GetSignatureHash 获取签名哈希
func (t *txFacade) GetSignatureHash(tx *types.Transaction, inputIndex int, scriptPubKey []byte, hashType types.SigHashType) ([]byte, error) {
	// TODO: 实现签名哈希计算
	return nil, nil
}
