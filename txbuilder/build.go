// 组 MsgTx、找零、RBF/locktime
package txbuilder

// import (
// 	"bytes"
// 	"context"
// 	"encoding/hex"
// 	"errors"
// 	"math"

// 	"github.com/btcsuite/btcd/btcutil"
// 	"github.com/btcsuite/btcd/chaincfg"
// 	"github.com/btcsuite/btcd/chaincfg/chainhash"
// 	"github.com/btcsuite/btcd/txscript"
// 	"github.com/btcsuite/btcd/wire"
// 	"github.com/crazycloudcc/btcapis/types"
// )

// const (
// 	DefaultDust = int64(546)
// )

// // Build - 最小可编译骨架：
// //   - 选择输入（目前直接使用 req.Inputs）
// //   - 估算费率与 vsize
// //   - 组装 wire.MsgTx（未签名）
// //   - 生成未签名原始交易 hex
// //   - PSBT 留空（后续在 psbt.go 实现）
// func Build(ctx context.Context, be interface{}, req *types.BuildTxRequest) (*types.BuildTxResult, error) {
// 	if req == nil {
// 		return nil, errors.New("nil request")
// 	}
// 	if len(req.Outputs) == 0 {
// 		return nil, errors.New("no outputs")
// 	}
// 	if len(req.Inputs) == 0 {
// 		return nil, errors.New("no inputs")
// 	}

// 	// 兜底参数
// 	if req.FeeRateSatPerVb <= 0 {
// 		req.FeeRateSatPerVb = 5.0
// 	}
// 	if req.DustLimit <= 0 {
// 		req.DustLimit = DefaultDust
// 	}
// 	sighash := req.SighashType
// 	if sighash == 0 {
// 		sighash = 0x01 // SIGHASH_ALL
// 	}

// 	// 选择输入
// 	selected, inSum := SelectInputsGreedy(req)

// 	// 估算：先假设有找零
// 	var (
// 		inputV  int64
// 		outputV int64
// 	)
// 	for _, in := range selected {
// 		st := in.ScriptType
// 		if st == "" && len(in.ScriptPubKey) > 0 {
// 			st = ClassifyScript(in.ScriptPubKey)
// 		}
// 		inputV += EstimateInputVSize(st)
// 	}
// 	// 已有输出
// 	for range req.Outputs {
// 		// 这里只能假设为 p2wpkh（估算）；真实脚本在构建时使用地址类型决定
// 		outputV += EstimateOutputVSize("p2wpkh")
// 	}
// 	// 假设找零为 p2wpkh
// 	outputVWithChange := outputV + EstimateOutputVSize("p2wpkh")
// 	vsizeWithChange := int64(10) + inputV + outputVWithChange
// 	feeWithChange := int64(math.Ceil(float64(vsizeWithChange) * req.FeeRateSatPerVb))

// 	// 计算找零
// 	outSum := int64(0)
// 	for _, o := range req.Outputs {
// 		outSum += o.Value
// 	}
// 	change := inSum - outSum - feeWithChange

// 	useChange := change >= req.DustLimit && (req.MinChange == 0 || change >= req.MinChange)
// 	if !useChange {
// 		// 无找零重算
// 		vsize := int64(10) + inputV + outputV
// 		fee := int64(math.Ceil(float64(vsize) * req.FeeRateSatPerVb))
// 		change = inSum - outSum - fee
// 		if change < 0 {
// 			return nil, errors.New("insufficient funds (no change)")
// 		}
// 	}

// 	// 组装未签名交易
// 	msg := &wire.MsgTx{Version: 2, LockTime: req.LockTime}
// 	for _, in := range selected {
// 		var h chainhash.Hash
// 		if err := chainhash.Decode(&h, in.TxID); err != nil {
// 			return nil, err
// 		}
// 		ti := wire.NewTxIn(&wire.OutPoint{Hash: h, Index: in.Vout}, nil, nil)
// 		if req.EnableRBF {
// 			ti.Sequence = 0xfffffffd
// 		} else {
// 			ti.Sequence = 0xffffffff
// 		}
// 		msg.AddTxIn(ti)
// 	}
// 	// 输出脚本：按地址生成 pkScript
// 	params := req.Network.ToParams()
// 	for _, o := range req.Outputs {
// 		pk, err := addrToPkScript(o.Address, params)
// 		if err != nil {
// 			return nil, err
// 		}
// 		msg.AddTxOut(&wire.TxOut{Value: o.Value, PkScript: pk})
// 	}
// 	// 找零
// 	if useChange {
// 		if req.ChangeAddress == "" {
// 			return nil, errors.New("change address required when change is used")
// 		}
// 		pk, err := addrToPkScript(req.ChangeAddress, params)
// 		if err != nil {
// 			return nil, err
// 		}
// 		msg.AddTxOut(&wire.TxOut{Value: change, PkScript: pk})
// 	}

// 	// 未签名原始交易（不含见证）
// 	var nobuf bytes.Buffer
// 	if err := msg.SerializeNoWitness(&nobuf); err != nil {
// 		return nil, err
// 	}
// 	unsignedHex := hex.EncodeToString(nobuf.Bytes())

// 	// PSBT（占位，后续实现）
// 	psbtB64, _ := MakePSBT(msg, selected, sighash, be)

// 	// 估算 vsize（以最终是否带找零为准）
// 	var vsizeEst int64
// 	if useChange {
// 		vsizeEst = int64(10) + inputV + outputVWithChange
// 	} else {
// 		vsizeEst = int64(10) + inputV + outputV
// 	}
// 	feePaid := inSum - outSum - func() int64 {
// 		if useChange {
// 			return change
// 		}
// 		return 0
// 	}()

// 	return &types.BuildTxResult{
// 		UnsignedTxHex: unsignedHex,
// 		PSBTBase64:    psbtB64,
// 		SelectedUTXO:  selected,
// 		VSizeEstimate: vsizeEst,
// 		FeePaid:       feePaid,
// 		ChangeValue: func() int64 {
// 			if useChange {
// 				return change
// 			}
// 			return 0
// 		}(),
// 	}, nil
// }

// func addrToPkScript(addr string, params *chaincfg.Params) ([]byte, error) {
// 	a, err := btcutil.DecodeAddress(addr, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return txscript.PayToAddrScript(a)
// }
