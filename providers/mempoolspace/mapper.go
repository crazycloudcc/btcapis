// Package mempoolspace 数据映射器
package mempoolspace

import (
	"time"

	"github.com/yourusername/btcapis/types"
)

// Mapper mempool.space数据映射器
type Mapper struct{}

// NewMapper 创建新的映射器
func NewMapper() *Mapper {
	return &Mapper{}
}

// MapFeeResponse 映射手续费响应
func (m *Mapper) MapFeeResponse(resp *FeeResponse) *types.FeeEstimate {
	if resp == nil {
		return nil
	}

	return &types.FeeEstimate{
		TargetBlocks: 1, // 最快确认
		FeeRate:      float64(resp.FastestFee),
		FeeRateBTC:   float64(resp.FastestFee) / 100000000.0,
		Confidence:   1.0,
		EstimatedAt:  time.Now(),
		Backend:      "mempool.space",
	}
}

// MapMempoolResponse 映射内存池响应
func (m *Mapper) MapMempoolResponse(resp *MempoolResponse) *types.MempoolInfo {
	if resp == nil {
		return nil
	}

	return &types.MempoolInfo{
		Size:             resp.Count,
		Bytes:            resp.VSize,
		Usage:            resp.VSize,
		MaxMempool:       0, // mempool.space不提供此信息
		MempoolMinFee:    0, // mempool.space不提供此信息
		MinRelayFee:      0, // mempool.space不提供此信息
		UnbroadcastCount: 0, // mempool.space不提供此信息
		LastUpdated:      time.Now(),
		Backend:          "mempool.space",
	}
}

// MapBlockResponse 映射区块响应
func (m *Mapper) MapBlockResponse(resp *BlockResponse) *types.BlockHeader {
	if resp == nil {
		return nil
	}

	return &types.BlockHeader{
		Hash:             resp.ID,
		Version:          int32(resp.Version),
		PreviousHash:     resp.PreviousBlockHash,
		MerkleRoot:       resp.MerkleRoot,
		Timestamp:        resp.Timestamp,
		Bits:             0, // mempool.space不提供此信息
		Nonce:            0, // mempool.space不提供此信息
		Height:           resp.Height,
		Size:             resp.Size,
		Weight:           resp.Weight,
		TransactionCount: resp.TxCount,
	}
}

// MapTransactionResponse 映射交易响应
func (m *Mapper) MapTransactionResponse(resp *TransactionResponse) *types.Transaction {
	if resp == nil {
		return nil
	}

	inputs := make([]types.TxInput, len(resp.Inputs))
	for i, input := range resp.Inputs {
		inputs[i] = types.TxInput{
			OutPoint: types.OutPoint{
				TxID: input.TxID,
				Vout: uint32(input.Vout),
			},
			ScriptSig: []byte(input.ScriptSig),
			Witness:   make([][]byte, len(input.Witness)),
			Sequence:  uint32(input.Sequence),
		}

		// 转换witness
		for j, w := range input.Witness {
			inputs[i].Witness[j] = []byte(w)
		}
	}

	outputs := make([]types.TxOutput, len(resp.Outputs))
	for i, output := range resp.Outputs {
		outputs[i] = types.TxOutput{
			Value:        output.Value,
			ScriptPubKey: []byte(output.ScriptPubKey),
			Address:      output.ScriptPubKeyAddress,
			Spent:        false, // mempool.space不提供此信息
		}
	}

	return &types.Transaction{
		TxID:        resp.TxID,
		Version:     int32(resp.Version),
		LockTime:    uint32(resp.LockTime),
		Inputs:      inputs,
		Outputs:     outputs,
		Size:        resp.Size,
		Weight:      resp.Weight,
		Fee:         resp.Fee,
		BlockHeight: resp.Status.BlockHeight,
		BlockHash:   resp.Status.BlockHash,
		BlockTime:   resp.Status.BlockTime,
		Confirmed:   resp.Status.Confirmed,
	}
}
