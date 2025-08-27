package tx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"sort"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/psbt"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// 转账交易-PSBT预览: 通过输入数据根据发起转账钱包地址的类型创建对应的PSBT交易数据, 这个数据将提交给外部okx插件钱包等进行签名.
func (c *Client) buildPSBT(ctx context.Context, inputParams *types.TxInputParams) (*psbt.BuildResult, error) {

	// 计算总输出金额（satoshi）
	totalOutputSats := int64(0)
	for _, amount := range inputParams.AmountBTC {
		if amount <= 0 {
			return nil, fmt.Errorf("invalid amount: %f", amount)
		}
		totalOutputSats += int64(amount * 1e8)
	}

	// 获取费率（sat/vB）
	feeRate := inputParams.FeeRate
	if feeRate <= 0 {
		// 默认费率：1 sat/vB
		feeRate = 1.0
	}

	// 获取当前区块高度作为 locktime
	locktime := inputParams.Locktime
	if locktime == 0 {
		blockCount, err := c.bitcoindrpcClient.ChainGetBlockCount(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block count: %w", err)
		}
		locktime = int64(blockCount)
	}

	// 1. 选币：从所有输入地址收集 UTXO
	allUTXOs := make([]types.TxUTXO, 0)
	totalInputSats := int64(0)

	for _, fromAddr := range inputParams.FromAddress {
		utxos, err := c.addressClient.GetAddressUTXOs(ctx, fromAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to get UTXOs for %s: %w", fromAddr, err)
		}
		allUTXOs = append(allUTXOs, utxos...)
		for _, utxo := range utxos {
			totalInputSats += utxo.Value
		}
	}

	if totalInputSats < totalOutputSats {
		return nil, fmt.Errorf("insufficient funds: have %d sats, need %d sats", totalInputSats, totalOutputSats)
	}

	// 2. 选币算法：先尝试 BnB 精确匹配，失败再 knapsack
	selectedUTXOs, changeAmount := selectCoins(allUTXOs, totalOutputSats, feeRate)
	fmt.Printf("changeAmount: %d\n", changeAmount)

	addrScriptInfo, err := decoders.DecodeAddress(inputParams.FromAddress[0])
	fmt.Printf("========== RedeemScriptHashHex: %s\n", hex.EncodeToString(addrScriptInfo.RedeemScriptHashHex))
	fmt.Printf("========== WitnessProgramHex: %s\n", hex.EncodeToString(addrScriptInfo.WitnessProgramHex))
	fmt.Printf("========== TaprootOutputKeyHex: %s\n", hex.EncodeToString(addrScriptInfo.TaprootOutputKeyHex))

	// 3. 将TxUTXO转为PsbtUTXO结构
	psbtUTXOs := make([]psbt.PsbtUTXO, 0, len(selectedUTXOs))
	for _, utxo := range selectedUTXOs {

		pkScript := utxo.PkScript
		if len(pkScript) == 0 {
			// 如果没有 pkScript，尝试通过地址解析
			pkScript = addrScriptInfo.ScriptPubKeyHex
		}

		nonWitnessTxHex := ""
		if decoders.PKScriptToType(utxo.PkScript) == types.AddrP2PKH {
			txRaw, err := c.bitcoindrpcClient.TxGetRaw(ctx, utxo.OutPoint.Hash.String(), false)
			if err != nil {
				return nil, fmt.Errorf("failed to get raw tx for %s: %w", utxo.OutPoint.Hash.String(), err)
			}
			nonWitnessTxHex = hex.EncodeToString(txRaw)
		}

		psbtUTXOs = append(psbtUTXOs, psbt.PsbtUTXO{
			TxID:             utxo.OutPoint.Hash.String(),
			Vout:             utxo.OutPoint.Index,
			ValueSat:         utxo.Value,
			ScriptPubKeyHex:  hex.EncodeToString(pkScript),
			NonWitnessTxHex:  nonWitnessTxHex,
			RedeemScriptHex:  hex.EncodeToString(addrScriptInfo.RedeemScriptHashHex),
			WitnessScriptHex: hex.EncodeToString(addrScriptInfo.WitnessProgramHex),
		})
	}

	result, err := psbt.CreatePSBTForOKX(inputParams, psbtUTXOs, types.CurrentNetworkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create PSBT: %w", err)
	}

	return result, nil
}

func getRedeemScript(pk string) ([]byte, error) {
	pubkeyBytes, _ := hex.DecodeString(pk)       // 来自 OKX
	pkh := btcutil.Hash160(pubkeyBytes)          // RIPEMD160(SHA256(pubkey))
	redeem := append([]byte{0x00, 0x14}, pkh...) // 0x0014 || pkh
	redeemScriptHex := hex.EncodeToString(redeem)
	fmt.Printf("redeemScriptHex: %s\n", redeemScriptHex)
	return redeem, nil
}

// 选币算法：先尝试 BnB 精确匹配，失败再 knapsack
func selectCoins(utxos []types.TxUTXO, targetAmount int64, feeRate float64) ([]types.TxUTXO, int64) {
	// 按价值排序（降序）
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value > utxos[j].Value
	})

	// 估算交易大小
	estimatedVSize := estimateTransactionVSize(len(utxos), 2) // 假设2个输出
	estimatedFee := int64(math.Ceil(float64(estimatedVSize) * feeRate))

	totalNeeded := targetAmount + estimatedFee

	// 尝试 BnB 精确匹配
	selected := make([]types.TxUTXO, 0)
	currentSum := int64(0)

	for _, utxo := range utxos {
		if currentSum >= totalNeeded {
			break
		}
		selected = append(selected, utxo)
		currentSum += utxo.Value
	}

	// 如果 BnB 成功，计算找零
	if currentSum >= totalNeeded {
		change := currentSum - totalNeeded
		return selected, change
	}

	// BnB 失败，使用 knapsack 贪心算法
	selected = make([]types.TxUTXO, 0)
	currentSum = int64(0)

	for _, utxo := range utxos {
		selected = append(selected, utxo)
		currentSum += utxo.Value

		// 重新估算费用
		newVSize := estimateTransactionVSize(len(selected), 2)
		newFee := int64(math.Ceil(float64(newVSize) * feeRate))

		if currentSum >= targetAmount+newFee {
			change := currentSum - targetAmount - newFee
			return selected, change
		}
	}

	// 如果还是不够，返回所有可用的 UTXO
	return utxos, 0
}

// 估算交易大小（vsize）
func estimateTransactionVSize(inputCount, outputCount int) int64 {
	// 固定开销：版本(4) + 输入计数(1) + 输出计数(1) + locktime(4) + marker/flag(2) = 12 bytes
	baseSize := int64(12)

	// 输入大小估算
	inputSize := int64(0)
	for i := 0; i < inputCount; i++ {
		// 假设平均为 P2WPKH 输入：68 vB
		inputSize += 68
	}

	// 输出大小估算
	outputSize := int64(0)
	for i := 0; i < outputCount; i++ {
		// 假设平均为 P2WPKH 输出：31 vB
		outputSize += 31
	}

	return baseSize + inputSize + outputSize
}

// // 计算 HASH160 (RIPEMD160(SHA256(data)))
// func hash160(data []byte) []byte {
// 	sha256Hash := sha256.Sum256(data)
// 	ripemd160Hash := ripemd160.New()
// 	ripemd160Hash.Write(sha256Hash[:])
// 	return ripemd160Hash.Sum(nil)
// }

// // 序列化 PSBT 为字节
// func serializePSBT(packet *psbt.Packet) ([]byte, error) {
// 	// 这里需要实现 PSBT 序列化
// 	// 由于现有的 psbt 包还没有完整的序列化实现，先返回一个占位实现
// 	// TODO: 实现完整的 PSBT 序列化

// 	// 临时实现：返回一个简单的标识
// 	return []byte("PSBT_PLACEHOLDER"), nil
// }
