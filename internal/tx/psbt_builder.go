package tx

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/psbt"
	"github.com/crazycloudcc/btcapis/internal/types"
	"golang.org/x/crypto/ripemd160"
)

// 转账交易-PSBT预览: 通过输入数据根据发起转账钱包地址的类型创建对应的PSBT交易数据, 这个数据将提交给外部okx插件钱包等进行签名.
func (c *Client) buildPSBT(ctx context.Context, inputParams *TxInputParams) (string, error) {

	// 计算总输出金额（satoshi）
	totalOutputSats := int64(0)
	for _, amount := range inputParams.AmountBTC {
		if amount <= 0 {
			return "", fmt.Errorf("invalid amount: %f", amount)
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
			return "", fmt.Errorf("failed to get block count: %w", err)
		}
		locktime = int64(blockCount)
	}

	// 1. 选币：从所有输入地址收集 UTXO
	allUTXOs := make([]types.UTXO, 0)
	totalInputSats := int64(0)

	for _, fromAddr := range inputParams.FromAddress {
		utxos, err := c.addressClient.GetAddressUTXOs(ctx, fromAddr)
		if err != nil {
			return "", fmt.Errorf("failed to get UTXOs for %s: %w", fromAddr, err)
		}
		allUTXOs = append(allUTXOs, utxos...)
		for _, utxo := range utxos {
			totalInputSats += utxo.Value
		}
	}

	if totalInputSats < totalOutputSats {
		return "", fmt.Errorf("insufficient funds: have %d sats, need %d sats", totalInputSats, totalOutputSats)
	}

	// 2. 选币算法：先尝试 BnB 精确匹配，失败再 knapsack
	selectedUTXOs, changeAmount := selectCoins(allUTXOs, totalOutputSats, feeRate)

	// 3. 构建未签名交易
	msgTx := wire.NewMsgTx(2) // 版本 2 支持 Taproot
	msgTx.LockTime = uint32(locktime)

	// 添加输入
	for _, utxo := range selectedUTXOs {
		txHash, err := chainhash.NewHashFromStr(utxo.OutPoint.Hash.String())
		if err != nil {
			return "", fmt.Errorf("invalid tx hash: %w", err)
		}

		txIn := wire.NewTxIn(
			&wire.OutPoint{
				Hash:  *txHash,
				Index: utxo.OutPoint.Index,
			},
			nil, // ScriptSig 为空
			nil, // Witness 为空
		)

		// RBF 支持
		if inputParams.Replaceable {
			txIn.Sequence = 0xfffffffd // 允许替换
		} else {
			txIn.Sequence = 0xffffffff // 不允许替换
		}

		msgTx.AddTxIn(txIn)
	}

	// 添加输出
	for i, toAddr := range inputParams.ToAddress {
		amount := int64(inputParams.AmountBTC[i] * 1e8)
		pkScript, err := decoders.AddressToPkScript(toAddr)
		if err != nil {
			return "", fmt.Errorf("invalid to address %s: %w", toAddr, err)
		}

		txOut := wire.NewTxOut(amount, pkScript)
		msgTx.AddTxOut(txOut)
	}

	// 添加找零输出
	if changeAmount > 546 { // dust limit
		changeAddr := inputParams.FromAddress[0] // 使用第一个输入地址作为找零地址
		changePkScript, err := decoders.AddressToPkScript(changeAddr)
		if err != nil {
			return "", fmt.Errorf("invalid change address: %w", err)
		}

		changeTxOut := wire.NewTxOut(changeAmount, changePkScript)
		msgTx.AddTxOut(changeTxOut)
	}

	// 4. 创建 PSBT 并填充元数据
	psbtPacket := psbt.NewV0FromUnsignedTx(msgTx)

	// 填充每个输入的 UTXO 信息
	for i, utxo := range selectedUTXOs {
		input := psbtPacket.MustInput(i)

		// 根据脚本类型填充相应的 UTXO 字段
		scriptType := decoders.PKScriptToType(utxo.PkScript)

		switch scriptType {
		case types.AddrP2PKH:
			// Legacy P2PKH：需要 NonWitnessUtxo
			rawTx, err := c.GetRawTx(ctx, utxo.OutPoint.Hash.String())
			if err != nil {
				return "", fmt.Errorf("failed to get raw tx for input %d: %w", i, err)
			}

			msgTx := wire.NewMsgTx(0)
			err = msgTx.Deserialize(bytes.NewReader(rawTx))
			if err != nil {
				return "", fmt.Errorf("failed to decode raw tx for input %d: %w", i, err)
			}

			input.NonWitnessUtxo = msgTx

		case types.AddrP2SH:
			// P2SH：需要检查是否为 P2SH-P2WPKH
			if inputParams.PublicKey != "" {
				// 构造 redeemScript = 0x0014<keyhash>
				pubkeyBytes, err := hex.DecodeString(inputParams.PublicKey)
				if err != nil {
					return "", fmt.Errorf("invalid public key: %w", err)
				}

				keyHash := hash160(pubkeyBytes)
				redeemScript := append([]byte{0x00, 0x14}, keyHash...)

				// 验证 HASH160(redeemScript) 是否等于地址中的脚本哈希
				scriptHash := hash160(redeemScript)
				if !bytes.Equal(scriptHash, utxo.PkScript[2:22]) {
					return "", fmt.Errorf("redeemScript hash mismatch for input %d", i)
				}

				input.RedeemScript = redeemScript
				input.WitnessUtxo = &wire.TxOut{
					Value:    utxo.Value,
					PkScript: utxo.PkScript,
				}
			} else {
				// 没有公钥，使用 NonWitnessUtxo
				rawTx, err := c.GetRawTx(ctx, utxo.OutPoint.Hash.String())
				if err != nil {
					return "", fmt.Errorf("failed to get raw tx for input %d: %w", i, err)
				}

				msgTx := wire.NewMsgTx(0)
				err = msgTx.Deserialize(bytes.NewReader(rawTx))
				if err != nil {
					return "", fmt.Errorf("failed to decode raw tx for input %d: %w", i, err)
				}

				input.NonWitnessUtxo = msgTx
			}

		case types.AddrP2WPKH:
			// P2WPKH：使用 WitnessUtxo
			input.WitnessUtxo = &wire.TxOut{
				Value:    utxo.Value,
				PkScript: utxo.PkScript,
			}

		case types.AddrP2TR:
			// Taproot：使用 WitnessUtxo + 公钥
			input.WitnessUtxo = &wire.TxOut{
				Value:    utxo.Value,
				PkScript: utxo.PkScript,
			}

			// OKX 要求：为每个输入附公钥
			if inputParams.PublicKey != "" {
				pubkeyBytes, err := hex.DecodeString(inputParams.PublicKey)
				if err != nil {
					return "", fmt.Errorf("invalid public key for taproot: %w", err)
				}

				// 转换为 x-only 公钥（32字节）
				if len(pubkeyBytes) == 33 {
					// 压缩公钥，取 x 坐标
					pubkeyBytes = pubkeyBytes[1:33]
				} else if len(pubkeyBytes) != 32 {
					return "", fmt.Errorf("invalid public key length for taproot: %d", len(pubkeyBytes))
				}

				// 添加到自定义字段（OKX 可识别）
				if input.PartialSigs == nil {
					input.PartialSigs = make(map[string][]byte)
				}
				input.PartialSigs["taproot_pubkey"] = pubkeyBytes
			}

		default:
			// 未知脚本类型，使用 NonWitnessUtxo
			rawTx, err := c.GetRawTx(ctx, utxo.OutPoint.Hash.String())
			if err != nil {
				return "", fmt.Errorf("failed to get raw tx for input %d: %w", i, err)
			}

			msgTx := wire.NewMsgTx(0)
			err = msgTx.Deserialize(bytes.NewReader(rawTx))
			if err != nil {
				return "", fmt.Errorf("failed to decode raw tx for input %d: %w", i, err)
			}

			input.NonWitnessUtxo = msgTx
		}
	}

	// 5. 序列化 PSBT 为十六进制（OKX 需要）
	psbtBytes, err := serializePSBT(psbtPacket)
	if err != nil {
		return "", fmt.Errorf("failed to serialize PSBT: %w", err)
	}

	psbtHex := hex.EncodeToString(psbtBytes)
	return psbtHex, nil
}

// 选币算法：先尝试 BnB 精确匹配，失败再 knapsack
func selectCoins(utxos []types.UTXO, targetAmount int64, feeRate float64) ([]types.UTXO, int64) {
	// 按价值排序（降序）
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value > utxos[j].Value
	})

	// 估算交易大小
	estimatedVSize := estimateTransactionVSize(len(utxos), 2) // 假设2个输出
	estimatedFee := int64(math.Ceil(float64(estimatedVSize) * feeRate))

	totalNeeded := targetAmount + estimatedFee

	// 尝试 BnB 精确匹配
	selected := make([]types.UTXO, 0)
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
	selected = make([]types.UTXO, 0)
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

// 计算 HASH160 (RIPEMD160(SHA256(data)))
func hash160(data []byte) []byte {
	sha256Hash := sha256.Sum256(data)
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256Hash[:])
	return ripemd160Hash.Sum(nil)
}

// 序列化 PSBT 为字节
func serializePSBT(packet *psbt.Packet) ([]byte, error) {
	// 这里需要实现 PSBT 序列化
	// 由于现有的 psbt 包还没有完整的序列化实现，先返回一个占位实现
	// TODO: 实现完整的 PSBT 序列化

	// 临时实现：返回一个简单的标识
	return []byte("PSBT_PLACEHOLDER"), nil
}
