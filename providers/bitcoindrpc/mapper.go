// Package bitcoindrpc 数据映射器
package bitcoindrpc

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/crazycloudcc/btcapis/types"
)

// Mapper Bitcoin Core RPC数据映射器
type Mapper struct{}

// NewMapper 创建新的映射器
func NewMapper() *Mapper {
	return &Mapper{}
}

// MapBlockHeader 映射区块头
func (m *Mapper) MapBlockHeader(raw interface{}) (*types.BlockHeader, error) {
	// TODO: 实现区块头映射
	return &types.BlockHeader{}, nil
}

// MapBlock 映射区块
func (m *Mapper) MapBlock(raw interface{}) (*types.Block, error) {
	// TODO: 实现区块映射
	return &types.Block{}, nil
}

// MapTransaction 映射交易
func (m *Mapper) MapTransaction(raw interface{}) (*types.Transaction, error) {
	// TODO: 实现交易映射
	return &types.Transaction{}, nil
}

// MapUTXO 映射UTXO
func (m *Mapper) MapUTXO(raw interface{}) (*types.UTXO, error) {
	// TODO: 实现UTXO映射
	return &types.UTXO{}, nil
}

// MapMempoolInfo 映射内存池信息
func (m *Mapper) MapMempoolInfo(raw interface{}) (*types.MempoolInfo, error) {
	// TODO: 实现内存池信息映射
	return &types.MempoolInfo{}, nil
}

// MapMempoolEntry 映射内存池条目
func (m *Mapper) MapMempoolEntry(raw interface{}) (*types.MempoolEntry, error) {
	// TODO: 实现内存池条目映射
	return &types.MempoolEntry{}, nil
}

// MapFeeEstimate 映射手续费估算
func (m *Mapper) MapFeeEstimate(raw interface{}) (*types.FeeEstimate, error) {
	// TODO: 实现手续费估算映射
	return &types.FeeEstimate{}, nil
}

// MapFeeHistory 映射手续费历史
func (m *Mapper) MapFeeHistory(raw interface{}) (*types.FeeHistory, error) {
	// TODO: 实现手续费历史映射
	return &types.FeeHistory{}, nil
}

// 辅助函数
func (m *Mapper) parseHexString(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

func (m *Mapper) parseTime(timestamp interface{}) (time.Time, error) {
	switch v := timestamp.(type) {
	case int64:
		return time.Unix(v, 0), nil
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return time.Unix(i, 0), nil
		}
		return time.Parse(time.RFC3339, v)
	default:
		return time.Time{}, nil
	}
}

func (m *Mapper) parseAmount(amount interface{}) (int64, error) {
	switch v := amount.(type) {
	case float64:
		return int64(v * 100000000), nil // 转换为聪
	case int64:
		return v, nil
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f * 100000000), nil
		}
		return 0, nil
	default:
		return 0, nil
	}
}
