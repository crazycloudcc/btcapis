package tx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"sort"

	"github.com/btcsuite/btcd/btcutil"
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

	// 3. 将TxUTXO转为PsbtUTXO结构
	redeemBytes, err := getRedeemScript(inputParams.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get redeem script: %w", err)
	}
	psbtUTXOs := make([]psbt.PsbtUTXO, 0, len(selectedUTXOs))
	for _, utxo := range selectedUTXOs {
		nonWitnessTxHex := ""
		redeemScriptHex := ""

		if isP2SH(utxo.PkScript) && len(redeemBytes) > 0 && isSegwitProgram(redeemBytes) {
			// P2SH-P2WPKH 或 P2SH-P2WSH
			// 需要提供 redeemScript
			redeemScriptHex = hex.EncodeToString(redeemBytes)
		} else if isLegacy(utxo.PkScript) { // legacy 类型需要提供前序交易
			raw, err := c.GetRawTx(ctx, utxo.OutPoint.Hash.String())
			if err != nil {
				return nil, fmt.Errorf("failed to get raw tx for UTXO %s:%d: %w", utxo.OutPoint.Hash.String(), utxo.OutPoint.Index, err)
			}
			nonWitnessTxHex = hex.EncodeToString(raw)
		}

		psbtUTXOs = append(psbtUTXOs, psbt.PsbtUTXO{
			TxID:            utxo.OutPoint.Hash.String(),
			Vout:            utxo.OutPoint.Index,
			ValueSat:        utxo.Value,
			ScriptPubKeyHex: fmt.Sprintf("%x", utxo.PkScript),
			NonWitnessTxHex: nonWitnessTxHex,
			RedeemScriptHex: redeemScriptHex,
		})
	}

	result, err := psbt.CreatePSBTForOKX(inputParams, psbtUTXOs, types.CurrentNetworkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create PSBT: %w", err)
	}

	return result, nil
}

// legacy 的定义：P2PKH 或 P2SH（注意：是否嵌套 SegWit 由外层逻辑结合 redeemScript 再判定）
func isLegacy(spk []byte) bool {
	return isP2PKH(spk) || isP2SH(spk)
}

// 判断是否为传统 P2PKH: OP_DUP OP_HASH160 0x14 <20-byte> OP_EQUALVERIFY OP_CHECKSIG
func isP2PKH(pkScript []byte) bool {
	if len(pkScript) != 25 {
		return false
	}
	return pkScript[0] == 0x76 && // OP_DUP
		pkScript[1] == 0xa9 && // OP_HASH160
		pkScript[2] == 0x14 && // PUSH_20
		pkScript[23] == 0x88 && // OP_EQUALVERIFY
		pkScript[24] == 0xac // OP_CHECKSIG
}

// 匹配传统P2SH
func isP2SH(pkScript []byte) bool {
	// 最典型 P2SH 长度 23 字节
	if len(pkScript) != 23 {
		return false
	}
	// OP_HASH160 (0xa9), PUSH_20 (0x14), ...20B..., OP_EQUAL (0x87)
	return pkScript[0] == 0xa9 && pkScript[1] == 0x14 && pkScript[22] == 0x87
}

// 判断给定脚本是否为“SegWit 程序”本体（用于 P2SH redeemScript 判定是否 P2WPKH/P2WSH/Taproot）
// 规则：首字节为 OP_0(0x00) 或 OP_1..OP_16(0x51..0x60)，随后紧跟一个 pushlen（最常见 20 或 32）
// 注意：这里严格限制为 20/32 字节数据长度，满足 P2WPKH / P2WSH / P2TR 的常见情况。
func isSegwitProgram(pkScript []byte) bool {
	n := len(pkScript)
	if n < 4 || n > 42 { // 2字节头 + 至少2字节数据，一般不超过 42
		return false
	}
	ver := pkScript[0]
	pushLen := int(pkScript[1])

	// 版本校验：OP_0 或 OP_1..OP_16
	if ver != 0x00 && (ver < 0x51 || ver > 0x60) {
		return false
	}

	// 长度一致性：pushLen 必须等于余下数据长度
	if 2+pushLen != n {
		return false
	}

	// 只接受标准长度：20（v0-pkh）或 32（v0-wsh / v1-taproot）
	return pushLen == 20 || pushLen == 32
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
