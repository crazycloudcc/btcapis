package tx

import (
	"context"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// 创建普通交易 wire.Tx, 可以转为PSBT.
// TxInputParams的参数需要由调用者校验.
func (c *Client) createNormalTx(ctx context.Context, inputParams *types.TxInputParams) (*wire.MsgTx, []*types.TxUTXO, error) {
	tx := wire.NewMsgTx(2)

	// 1. 计算总输出金额（satoshi）
	totalOutAmountSats := int64(0)
	for _, amount := range inputParams.AmountBTC {
		totalOutAmountSats += int64(amount * 1e8)
	}

	// 2. 获取当前区块高度作为 locktime
	tx.LockTime = uint32(inputParams.Locktime)

	// 3. 选币：从第一个输入地址收集 UTXO
	arrUTXOs, err := c.addressClient.GetAddressUTXOs(ctx, inputParams.FromAddress[0])
	if err != nil {
		return nil, nil, err
	}

	// 4. 筛选 UTXO
	selectedUTXOs, totalInputSats, err := selectUTXOs(arrUTXOs, totalOutAmountSats)
	if err != nil {
		return nil, nil, err
	}

	// 5. 初步估算交易大小（vsize）
	estVsize := estimateTxSize(len(selectedUTXOs), inputParams.FromAddress[0], inputParams.ToAddress, len(inputParams.Data))
	estFee := int64(float64(estVsize) * inputParams.FeeRate) // satoshi

	// 5.1 检查输入金额是否足够
	if totalInputSats < totalOutAmountSats+estFee {
		fmt.Printf("totalInputSats=%d, totalOutAmountSats=%d, estFee=%d\n", totalInputSats, totalOutAmountSats, estFee)
		return nil, nil, errors.New("insufficient funds after fee calculation")
	}

	// 6. 构建输入
	// 6.1 设置输入的 sequence 以支持 RBF
	seq := uint32(wire.MaxTxInSequenceNum)
	if inputParams.Replaceable {
		seq -= 2
	}

	// 6.2 添加输入
	// 注意：这里的输入顺序会影响最终的签名顺序.
	// 因为本接口只支持单个from address, 所以不需要复杂的排序逻辑.
	for i := 0; i < len(selectedUTXOs); i++ {
		h, err := chainhash.NewHashFromStr(selectedUTXOs[i].OutPoint.Hash.String())
		if err != nil {
			return nil, nil, fmt.Errorf("TxID(utxo.OutPoint.Hash) 解析失败: %v", err)
		}

		tx.AddTxIn(&wire.TxIn{
			PreviousOutPoint: wire.OutPoint{Hash: *h, Index: selectedUTXOs[i].OutPoint.Index},
			Sequence:         seq,
		})
	}

	// 7. 构建输出
	// 7.1 添加普通输出
	for i, toAddr := range inputParams.ToAddress {
		pkScript, err := decoders.AddressToPkScript(toAddr)
		if err != nil {
			return nil, nil, fmt.Errorf("解析收款地址失败: %v", err)
		}

		tx.AddTxOut(&wire.TxOut{
			Value:    int64(inputParams.AmountBTC[i] * 1e8),
			PkScript: pkScript,
		})
	}

	// 7.2 添加OP_RETURN输出（如果有）
	if inputParams.Data != "" {
		script, _ := txscript.NewScriptBuilder().
			AddOp(txscript.OP_RETURN).
			AddData([]byte(inputParams.Data)).
			Script()
		tx.AddTxOut(&wire.TxOut{Value: 0, PkScript: script})
	}

	// 7.3 添加找零输出（如果有）
	changeSats := totalInputSats - totalOutAmountSats - estFee
	if changeSats > 0 {
		// 计算找零和新增输出的差值: 如果
		changeAddrType, err := decoders.AddressToType(inputParams.ChangeAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("解析找零地址失败: %v", err)
		}
		_, addOutNeedSats := types.GetOutSize(changeAddrType)

		// 找零金额需要再减去新增输出的费用, 如果还有剩余才进行找零
		if changeSats > int64(addOutNeedSats) {
			changePkScript, err := decoders.AddressToPkScript(inputParams.ChangeAddress)
			if err != nil {
				return nil, nil, fmt.Errorf("解析找零地址失败: %v", err)
			}
			tx.AddTxOut(&wire.TxOut{
				Value:    changeSats,
				PkScript: changePkScript,
			})
		}
	}

	return tx, selectedUTXOs, nil
}

// helper: 筛选utxo
func selectUTXOs(utxos []types.TxUTXO, targetSats int64) ([]*types.TxUTXO, int64, error) {
	selected := make([]*types.TxUTXO, 0)

	total := int64(0)
	for _, utxo := range utxos {
		selected = append(selected, &utxo)
		total += utxo.Value
		if total >= targetSats {
			break
		}
	}

	// 余额不足
	if total < targetSats {
		return nil, 0, errors.New("insufficient funds")
	}

	return selected, total, nil
}

// helper: 初步估算交易大小: 目前对P2SH和P2WSH的支持较弱, 仅做参考.
// inCount: 输入只需要数量, 因为都是归属于from address的utxo.
func estimateTxSize(inCount int, fromAddr string, toAddrs []string, opReturnDataLen int) int {
	vsize := 0

	fromAddrType, err := decoders.AddressToType(fromAddr)
	if err != nil {
		fromAddrType = types.AddrUnknown
	}

	// 输入大小
	vsize += types.GetInSize(fromAddrType) * inCount

	// 输出大小
	for _, toAddr := range toAddrs {
		addrType, err := decoders.AddressToType(toAddr)
		if err != nil {
			addrType = types.AddrUnknown
		}
		_, outVsize := types.GetOutSize(addrType)
		vsize += outVsize
	}

	// OP_RETURN 输出大小
	vsize += types.GetOpReturnSize(opReturnDataLen)

	return vsize
}
