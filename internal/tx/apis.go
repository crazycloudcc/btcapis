package tx

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"

	"bytes"
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/psbt"
	"github.com/crazycloudcc/btcapis/internal/types"
	"golang.org/x/crypto/ripemd160"
)

// 转账交易-PSBT预览: 通过输入数据根据发起转账钱包地址的类型创建对应的PSBT交易数据, 这个数据将提交给外部okx插件钱包等进行签名.
func (c *Client) SendBTCByPSBTPreview(ctx context.Context, inputParams *TxInputParams) (string, error) {
	// 参数验证
	if len(inputParams.FromAddress) == 0 {
		return "", fmt.Errorf("from address is required")
	}
	if len(inputParams.ToAddress) == 0 {
		return "", fmt.Errorf("to address is required")
	}
	if len(inputParams.ToAddress) != len(inputParams.AmountBTC) {
		return "", fmt.Errorf("to address and amount count mismatch")
	}

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
		// 默认费率：6 sat/vB
		feeRate = 6.0
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
		pkScript, err := addressToPkScript(toAddr)
		if err != nil {
			return "", fmt.Errorf("invalid to address %s: %w", toAddr, err)
		}

		txOut := wire.NewTxOut(amount, pkScript)
		msgTx.AddTxOut(txOut)
	}

	// 添加找零输出
	if changeAmount > 546 { // dust limit
		changeAddr := inputParams.FromAddress[0] // 使用第一个输入地址作为找零地址
		changePkScript, err := addressToPkScript(changeAddr)
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
		scriptType := classifyScript(utxo.PkScript)

		switch scriptType {
		case "p2pkh":
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

		case "p2sh":
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

		case "p2wpkh":
			// P2WPKH：使用 WitnessUtxo
			input.WitnessUtxo = &wire.TxOut{
				Value:    utxo.Value,
				PkScript: utxo.PkScript,
			}

		case "p2tr":
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

// 地址转换为 PkScript
func addressToPkScript(addr string) ([]byte, error) {
	address, err := btcutil.DecodeAddress(addr, types.CurrentNetworkParams)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(address)
}

// 脚本类型识别
func classifyScript(pkScript []byte) string {
	if len(pkScript) == 0 {
		return "unknown"
	}

	// P2PKH: OP_DUP OP_HASH160 PUSH20 <20> OP_EQUALVERIFY OP_CHECKSIG
	if len(pkScript) == 25 && pkScript[0] == 0x76 && pkScript[1] == 0xa9 && pkScript[2] == 0x14 && pkScript[23] == 0x88 && pkScript[24] == 0xac {
		return "p2pkh"
	}
	// P2SH: OP_HASH160 PUSH20 <20> OP_EQUAL
	if len(pkScript) == 23 && pkScript[0] == 0xa9 && pkScript[1] == 0x14 && pkScript[22] == 0x87 {
		return "p2sh"
	}
	// P2WPKH v0: OP_0 PUSH20 <20>
	if len(pkScript) == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14 {
		return "p2wpkh"
	}
	// P2WSH v0: OP_0 PUSH32 <32>
	if len(pkScript) == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20 {
		return "p2wsh"
	}
	// P2TR v1: OP_1 PUSH32 <32>
	if len(pkScript) == 34 && pkScript[0] == 0x51 && pkScript[1] == 0x20 {
		return "p2tr"
	}

	return "unknown"
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

// 接收OKX签名后的交易数据并广播
func (c *Client) SendBTCByPSBT(ctx context.Context, psbt string) (string, error) {
	// 兼容 OKX psbtHex 与 base64 两种输入
	normalized := strings.TrimSpace(psbt)
	var psbtBase64 string
	// 判定十六进制
	isHex := func(s string) bool {
		if len(s)%2 != 0 || len(s) == 0 {
			return false
		}
		for i := 0; i < len(s); i++ {
			ch := s[i]
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				return false
			}
		}
		return true
	}
	if isHex(normalized) {
		bin, err := hex.DecodeString(normalized)
		if err != nil {
			return "", err
		}
		psbtBase64 = base64.StdEncoding.EncodeToString(bin)
	} else {
		psbtBase64 = normalized
	}

	// finalizepsbt -> 原始交易hex
	rawHex, err := c.bitcoindrpcClient.TxFinalizePsbt(ctx, psbtBase64)
	if err != nil {
		return "", err
	}
	bin, err := hex.DecodeString(rawHex)
	if err != nil {
		return "", err
	}
	return c.bitcoindrpcClient.TxBroadcast(ctx, bin)
}

// // 普通转账交易预览 =>
// func (c *Client) SendBTCByPSBTPreview(ctx context.Context, inputParams *TxInputParams) ([]byte, error) {

// 	dto := bitcoindrpc.TxCreateRawDTO{
// 		Inputs: []bitcoindrpc.TxInputCreateRawDTO{
// 			{
// 				TxID: inputParams.FromAddress[0],
// 				Vout: 0,
// 			},
// 		},
// 		Outputs: []bitcoindrpc.TxOutputCreateRawDTO{},
// 	}

// 	for i, to := range inputParams.ToAddress {
// 		dto.Outputs = append(dto.Outputs, bitcoindrpc.TxOutputCreateRawDTO{
// 			Address: to,
// 			Amount:  inputParams.AmountBTC[i],
// 		})
// 	}

// 	if inputParams.Data != "" {
// 		dto.Outputs = append(dto.Outputs, bitcoindrpc.TxOutputCreateRawDTO{
// 			DataHex: inputParams.Data,
// 		})
// 	}

// 	// step.1 构建交易
// 	raw, err := c.createTx(ctx, dto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// step.2 填充交易费用
// 	funded, err := c.fundTx(ctx, string(raw), bitcoindrpc.TxFundOptionsDTO{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// step.3 签名交易
// 	signed, err := c.signTx(ctx, funded.Hex)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// step.4 广播交易
// 	broadcasted, err := c.broadcastTx(ctx, signed)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return broadcasted, nil
// }

// 查询交易元数据: 优先使用bitcoindrpc, 如果没有再使用mempoolapis
func (c *Client) GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	ret, err := c.bitcoindrpcClient.TxGetRaw(ctx, txid, false)
	if err == nil {
		return ret, nil
	}

	if c.mempoolapisClient != nil {
		ret, err = c.mempoolapisClient.TxGetRaw(ctx, txid)
		if err == nil {
			return ret, nil
		}
	}

	return nil, err
}

// 按照types.Tx格式返回交易数据
func (c *Client) GetTx(ctx context.Context, txid string) (*types.Tx, error) {
	raw, err := c.GetRawTx(ctx, txid)
	if err != nil {
		return nil, err
	}

	ret, err := decoders.DecodeRawTx(raw)
	if err != nil {
		return nil, err
	}

	return ret, err
}

// 构建交易 createrawtransaction
func (c *Client) createTx(ctx context.Context, dto bitcoindrpc.TxCreateRawDTO) ([]byte, error) {
	return c.bitcoindrpcClient.TxCreateRaw(ctx, dto)
}

// 填充交易费用
func (c *Client) fundTx(ctx context.Context, rawtx string, options bitcoindrpc.TxFundOptionsDTO) (bitcoindrpc.TxFundRawResultDTO, error) {
	return c.bitcoindrpcClient.TxFundRaw(ctx, rawtx, options)
}

// 签名交易
func (c *Client) signTx(ctx context.Context, rawtx string) (string, error) {
	return c.bitcoindrpcClient.TxSignRawWithKey(ctx, rawtx)
}

// 广播交易
func (c *Client) broadcastTx(ctx context.Context, rawtx string) (string, error) {
	bin, err := hex.DecodeString(rawtx)
	if err != nil {
		return "", err
	}
	return c.bitcoindrpcClient.TxBroadcast(ctx, bin)
}
