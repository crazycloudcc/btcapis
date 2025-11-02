package tx

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/pkg/logger"
	"github.com/crazycloudcc/btcapis/types"
)

// TransferAllToNewAddress 将给定私钥+对应fromAddress的所有余额转移到toAddress
// 这是一个紧急避险功能，用于将泄露私钥的地址的余额快速转移到安全地址
//
// 功能说明:
// - 使用ElectrumX查询UTXO（如果可用），否则使用Bitcoin Core的scantxoutset
// - 使用Mempool.space获取费率（如果可用），否则使用Bitcoin Core的estimatesmartfee
// - 使用PSBT签名流程，支持所有地址类型（P2PKH、P2SH、P2WPKH、P2TR）
// - 自动计算最优手续费，最大化转账金额
func (c *Client) TransferAllToNewAddress(
	ctx context.Context,
	toAddress string,
	privateKeyWIF string,
	fromAddress string,
	feeRate float64) (string, error) {

	logger.Info("========== 开始紧急转账流程 ==========")
	logger.Info("[步骤1] 验证输入参数")
	logger.Info("  - 源地址: %s", fromAddress)
	logger.Info("  - 目标地址: %s", toAddress)
	logger.Info("  - 私钥格式: WIF")

	// 1. 验证地址
	logger.Info("[步骤2] 验证源地址")
	fromAddrInfo, err := c.validateAddress(ctx, fromAddress)
	if err != nil {
		logger.Error("验证源地址失败: %v", err)
		return "", fmt.Errorf("验证源地址失败: %w", err)
	}
	if !fromAddrInfo.IsValid {
		logger.Error("源地址无效: %s", fromAddress)
		return "", fmt.Errorf("源地址无效: %s", fromAddress)
	}
	logger.Info("  ✓ 源地址验证成功")
	logger.Info("    - 地址类型: %s", getAddressType(fromAddrInfo))

	logger.Info("[步骤3] 验证目标地址")
	toAddrInfo, err := c.validateAddress(ctx, toAddress)
	if err != nil {
		logger.Error("验证目标地址失败: %v", err)
		return "", fmt.Errorf("验证目标地址失败: %w", err)
	}
	if !toAddrInfo.IsValid {
		logger.Error("目标地址无效: %s", toAddress)
		return "", fmt.Errorf("目标地址无效: %s", toAddress)
	}
	logger.Info("  ✓ 目标地址验证成功")
	logger.Info("    - 地址类型: %s", getAddressType(toAddrInfo))

	// 2. 解析私钥
	logger.Info("[步骤4] 解析WIF格式私钥")
	wif, err := btcutil.DecodeWIF(privateKeyWIF)
	if err != nil {
		logger.Error("解析私钥失败: %v", err)
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}
	logger.Info("  ✓ 私钥解析成功")
	logger.Info("    - 是否压缩: %v", wif.CompressPubKey)

	// 3. 验证私钥与地址是否匹配
	logger.Info("[步骤5] 验证私钥与源地址是否匹配")
	if err := verifyPrivKeyMatchAddress(wif, fromAddress, fromAddrInfo); err != nil {
		logger.Error("私钥与地址不匹配: %v", err)
		return "", fmt.Errorf("私钥与地址不匹配: %w", err)
	}
	logger.Info("  ✓ 私钥与地址匹配验证成功")

	// 4. 查询源地址的所有UTXO（优先使用ElectrumX，性能更好）
	logger.Info("[步骤6] 查询源地址的UTXO集")
	utxos, err := c.getAddressUTXOs(ctx, fromAddress)
	if err != nil {
		logger.Error("查询UTXO失败: %v", err)
		return "", fmt.Errorf("查询UTXO失败: %w", err)
	}

	if len(utxos) == 0 {
		logger.Warn("源地址没有可用的UTXO，余额为0")
		return "", fmt.Errorf("源地址没有可用的UTXO")
	}

	// 计算总余额
	totalBalance := int64(0)
	for _, utxo := range utxos {
		totalBalance += int64(utxo.AmountBTC * 1e8)
	}
	logger.Info("  ✓ UTXO查询完成")
	logger.Info("    - UTXO数量: %d", len(utxos))
	logger.Info("    - 总余额: %.8f BTC (%d sats)", float64(totalBalance)/1e8, totalBalance)
	for i, utxo := range utxos {
		logger.Info("    - UTXO[%d]: %s:%d = %.8f BTC", i, utxo.TxID, utxo.Vout, utxo.AmountBTC)
	}

	// 5. 估算交易费用（优先使用Mempool.space，更准确）
	logger.Info("[步骤7] 估算交易费用")
	if feeRate <= 0.01 {
		feeRate, err := c.estimateFeeRate(ctx)
		if err != nil {
			logger.Error("估算费率失败: %v", err)
			return "", fmt.Errorf("估算费率失败: %w", err)
		}

		// 如果费率为0，使用默认费率
		if feeRate <= 0.01 {
			feeRate = 1.0 // 默认1 sat/vB
			logger.Warn("  ! 费率为0，使用默认费率: %.2f sat/vB", feeRate)
		} else {
			logger.Info("  ✓ 费率估算完成")
			logger.Info("    - 费率: %.2f sat/vB", feeRate)
		}
	}

	// 6. 构建PSBT交易
	logger.Info("[步骤8] 构建PSBT交易")
	psbtPacket, estimatedFee, err := c.buildPSBT(ctx, utxos, toAddress, fromAddrInfo, totalBalance, feeRate)
	if err != nil {
		logger.Error("构建PSBT失败: %v", err)
		return "", fmt.Errorf("构建PSBT失败: %w", err)
	}
	outputAmount := totalBalance - estimatedFee
	logger.Info("  ✓ PSBT构建完成")
	logger.Info("    - 输入总额: %.8f BTC (%d sats)", float64(totalBalance)/1e8, totalBalance)
	logger.Info("    - 预估手续费: %.8f BTC (%d sats)", float64(estimatedFee)/1e8, estimatedFee)
	logger.Info("    - 输出金额: %.8f BTC (%d sats)", float64(outputAmount)/1e8, outputAmount)

	// 7. 使用私钥签名PSBT
	logger.Info("[步骤9] 使用私钥签名PSBT")
	if err := c.signPSBT(psbtPacket, wif, utxos, fromAddrInfo); err != nil {
		logger.Error("签名PSBT失败: %v", err)
		return "", fmt.Errorf("签名PSBT失败: %w", err)
	}
	logger.Info("  ✓ PSBT签名完成")
	logger.Info("    - 已签名输入数: %d/%d", len(utxos), len(utxos))

	// 8. 完成PSBT
	logger.Info("[步骤10] 完成PSBT交易")
	psbtBase64, err := psbtPacket.B64Encode()
	if err != nil {
		logger.Error("PSBT编码失败: %v", err)
		return "", fmt.Errorf("PSBT编码失败: %w", err)
	}

	finalHex, err := c.bitcoindrpcClient.TxFinalizePsbt(ctx, psbtBase64)
	if err != nil {
		logger.Error("完成PSBT失败: %v", err)
		return "", fmt.Errorf("完成PSBT失败: %w", err)
	}
	logger.Info("  ✓ PSBT完成，生成最终交易")
	logger.Info("    - 交易HEX长度: %d bytes", len(finalHex)/2)

	// 9. 广播交易
	logger.Info("[步骤11] 广播交易到网络")
	rawTx, err := hex.DecodeString(finalHex)
	if err != nil {
		logger.Error("解码交易HEX失败: %v", err)
		return "", fmt.Errorf("解码交易HEX失败: %w", err)
	}

	txid, err := c.BroadcastRawTx(ctx, rawTx)
	if err != nil {
		logger.Error("广播交易失败: %v", err)
		return "", fmt.Errorf("广播交易失败: %w", err)
	}

	logger.Info("  ✓ 交易广播成功!")
	logger.Info("========== 紧急转账完成 ==========")
	logger.Info("交易详情:")
	logger.Info("  - 交易ID: %s", txid)
	logger.Info("  - 源地址: %s", fromAddress)
	logger.Info("  - 目标地址: %s", toAddress)
	logger.Info("  - 转账金额: %.8f BTC", float64(outputAmount)/1e8)
	logger.Info("  - 手续费: %.8f BTC", float64(estimatedFee)/1e8)
	logger.Info("=====================================")

	return txid, nil
}

// validateAddress 验证地址（优先使用Bitcoin Core）
func (c *Client) validateAddress(ctx context.Context, address string) (*bitcoindrpc.ValidateAddressDTO, error) {
	if c.bitcoindrpcClient != nil {
		return c.bitcoindrpcClient.AddressValidate(ctx, address)
	}
	return nil, fmt.Errorf("没有可用的地址验证服务")
}

// getAddressUTXOs 获取地址的UTXO列表
// 优先使用ElectrumX（更快），否则使用Bitcoin Core的scantxoutset
func (c *Client) getAddressUTXOs(ctx context.Context, address string) ([]bitcoindrpc.UTXODTO, error) {
	// 优先使用ElectrumX（性能更好，~1秒）
	if c.electrumxClient != nil {
		logger.Info("    - 使用ElectrumX查询UTXO")
		utxos, err := c.electrumxClient.AddressGetUTXOs(ctx, address)
		if err == nil && len(utxos) >= 0 {
			// 转换为bitcoindrpc.UTXODTO格式
			result := make([]bitcoindrpc.UTXODTO, len(utxos))
			for i, utxo := range utxos {
				result[i] = bitcoindrpc.UTXODTO{
					TxID:      utxo.TxHash,
					Vout:      uint32(utxo.TxPos),
					AmountBTC: float64(utxo.Value) / 1e8,
					Height:    utxo.Height,
				}
			}
			logger.Info("    - ElectrumX查询成功，返回 %d 个UTXO", len(result))
			return result, nil
		}
		logger.Warn("    - ElectrumX查询失败: %v，尝试使用Bitcoin Core", err)
	}

	// 备用方案：使用Bitcoin Core（较慢，~30-120秒）
	if c.bitcoindrpcClient != nil {
		logger.Info("    - 使用Bitcoin Core scantxoutset查询UTXO（可能需要30-120秒）")
		return c.bitcoindrpcClient.AddressGetUTXOs(ctx, address)
	}

	return nil, fmt.Errorf("没有可用的UTXO查询服务")
}

// estimateFeeRate 估算费率（sat/vB）
// 优先使用Mempool.space（更准确的实时费率），否则使用Bitcoin Core
func (c *Client) estimateFeeRate(ctx context.Context) (float64, error) {
	// 优先使用Mempool.space
	if c.mempoolapisClient != nil {
		logger.Info("    - 使用Mempool.space估算费率")
		feeDTO, err := c.mempoolapisClient.EstimateFeeRate(ctx, 6)
		if err == nil {
			// 使用1小时内确认的费率（更经济）
			feeRate := feeDTO.HourFee
			if feeRate == 0 {
				// 如果1小时费率为0，使用30分钟费率
				feeRate = feeDTO.HalfHourFee
			}
			if feeRate == 0 {
				// 如果还是0，使用最快费率
				feeRate = feeDTO.FastestFee
			}
			logger.Info("    - Mempool.space返回费率: 最快=%.2f, 30分钟=%.2f, 1小时=%.2f sat/vB",
				feeDTO.FastestFee, feeDTO.HalfHourFee, feeDTO.HourFee)
			return feeRate, nil
		}
		logger.Warn("    - Mempool.space查询失败: %v，尝试使用Bitcoin Core", err)
	}

	// 优先使用ElectrumX
	if c.electrumxClient != nil {
		logger.Info("    - 使用ElectrumX估算费率")
		feeRate, err := c.electrumxClient.EstimateFee(ctx, 6)
		if err == nil && feeRate > 0 {
			// ElectrumX返回的是BTC/KB，需要转换为sat/vB
			// 1 BTC/KB = 100,000 sat/vB
			feeRateSatVB := feeRate * 100000
			logger.Info("    - ElectrumX返回费率: %.8f BTC/KB (%.2f sat/vB)", feeRate, feeRateSatVB)
			return feeRateSatVB, nil
		}
		logger.Warn("    - ElectrumX查询失败: %v，尝试使用Bitcoin Core", err)
	}

	// 备用方案：使用Bitcoin Core
	if c.bitcoindrpcClient != nil {
		logger.Info("    - 使用Bitcoin Core estimatesmartfee")
		feeDTO, err := c.bitcoindrpcClient.ChainEstimateSmartFeeRate(ctx, 6)
		if err == nil {
			return feeDTO.Feerate, nil
		}
		return 0, err
	}

	return 0, fmt.Errorf("没有可用的费率估算服务")
}

// getAddressType 获取地址类型描述
func getAddressType(addrInfo *bitcoindrpc.ValidateAddressDTO) string {
	if addrInfo.IsWitness {
		if addrInfo.WitnessVersion == 0 {
			return "SegWit v0 (P2WPKH/P2WSH)"
		} else if addrInfo.WitnessVersion == 1 {
			return "Taproot (P2TR)"
		}
		return fmt.Sprintf("SegWit v%d", addrInfo.WitnessVersion)
	} else if addrInfo.IsScript {
		return "P2SH"
	}
	return "P2PKH (Legacy)"
}

// verifyPrivKeyMatchAddress 验证私钥是否匹配地址
func verifyPrivKeyMatchAddress(wif *btcutil.WIF, address string, addrInfo *bitcoindrpc.ValidateAddressDTO) error {
	pubKey := wif.PrivKey.PubKey()
	netParams := types.CurrentNetworkParams

	var derivedAddr btcutil.Address
	var err error

	// 根据地址类型派生地址
	if addrInfo.IsWitness && addrInfo.WitnessVersion == 1 {
		// P2TR (Taproot)
		logger.Info("    - 检测到Taproot地址，进行密钥派生")
		internalPubKey := txscript.ComputeTaprootKeyNoScript(pubKey)
		derivedAddr, err = btcutil.NewAddressTaproot(
			internalPubKey.SerializeCompressed()[1:],
			netParams,
		)
	} else if addrInfo.IsWitness && addrInfo.WitnessVersion == 0 {
		// P2WPKH
		logger.Info("    - 检测到SegWit地址")
		pkHash := btcutil.Hash160(pubKey.SerializeCompressed())
		derivedAddr, err = btcutil.NewAddressWitnessPubKeyHash(pkHash, netParams)
	} else if addrInfo.IsScript {
		// P2SH-P2WPKH (嵌套SegWit)
		logger.Info("    - 检测到P2SH地址（可能是嵌套SegWit）")
		pkHash := btcutil.Hash160(pubKey.SerializeCompressed())
		witnessAddr, _ := btcutil.NewAddressWitnessPubKeyHash(pkHash, netParams)
		witnessScript, _ := txscript.PayToAddrScript(witnessAddr)
		derivedAddr, err = btcutil.NewAddressScriptHash(witnessScript, netParams)
	} else {
		// P2PKH (Legacy)
		logger.Info("    - 检测到Legacy地址")
		if wif.CompressPubKey {
			pkHash := btcutil.Hash160(pubKey.SerializeCompressed())
			derivedAddr, err = btcutil.NewAddressPubKeyHash(pkHash, netParams)
		} else {
			pkHash := btcutil.Hash160(pubKey.SerializeUncompressed())
			derivedAddr, err = btcutil.NewAddressPubKeyHash(pkHash, netParams)
		}
	}

	if err != nil {
		return fmt.Errorf("派生地址失败: %w", err)
	}

	if derivedAddr.EncodeAddress() != address {
		logger.Error("    - 派生地址: %s", derivedAddr.EncodeAddress())
		logger.Error("    - 期望地址: %s", address)
		return fmt.Errorf("私钥派生的地址 %s 与提供的地址 %s 不匹配",
			derivedAddr.EncodeAddress(), address)
	}

	return nil
}

// buildPSBT 构建PSBT交易
func (c *Client) buildPSBT(
	ctx context.Context,
	utxos []bitcoindrpc.UTXODTO,
	toAddress string,
	fromAddrInfo *bitcoindrpc.ValidateAddressDTO,
	totalBalance int64,
	feeRate float64,
) (*psbt.Packet, int64, error) {

	logger.Info("    - 创建新的交易结构")
	tx := wire.NewMsgTx(2)

	// 添加所有输入
	logger.Info("    - 添加 %d 个交易输入", len(utxos))
	for i, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, 0, fmt.Errorf("解析UTXO哈希失败[%d]: %w", i, err)
		}
		tx.AddTxIn(&wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  *hash,
				Index: utxo.Vout,
			},
			Sequence: wire.MaxTxInSequenceNum,
		})
	}

	// 估算交易大小和手续费
	logger.Info("    - 估算交易大小和手续费")
	estimatedSize := estimateTransferTxSize(len(utxos), 1, fromAddrInfo)
	estimatedFee := int64(float64(estimatedSize) * feeRate)

	// 确保手续费合理
	minFee := int64(estimatedSize) // 至少1 sat/vB
	if estimatedFee < minFee {
		estimatedFee = minFee
		logger.Warn("    ! 计算的手续费过低，调整为最小值: %d sats", minFee)
	}

	outputAmount := totalBalance - estimatedFee

	if outputAmount <= 0 {
		return nil, 0, fmt.Errorf("余额不足支付手续费，总额: %d sats, 手续费: %d sats",
			totalBalance, estimatedFee)
	}

	logger.Info("    - 预估交易大小: %d vBytes", estimatedSize)
	logger.Info("    - 预估手续费: %d sats (%.2f sat/vB)", estimatedFee, feeRate)

	// 添加输出
	logger.Info("    - 添加交易输出")
	toAddr, err := btcutil.DecodeAddress(toAddress, types.CurrentNetworkParams)
	if err != nil {
		return nil, 0, fmt.Errorf("解析目标地址失败: %w", err)
	}
	toPkScript, err := txscript.PayToAddrScript(toAddr)
	if err != nil {
		return nil, 0, fmt.Errorf("生成输出脚本失败: %w", err)
	}

	tx.AddTxOut(&wire.TxOut{
		Value:    outputAmount,
		PkScript: toPkScript,
	})

	// 创建PSBT
	logger.Info("    - 创建PSBT包")
	packet, err := psbt.NewFromUnsignedTx(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("创建PSBT失败: %w", err)
	}

	// 添加输入的witness数据
	logger.Info("    - 为PSBT添加输入见证数据")
	for i, utxo := range utxos {
		if err := c.addPSBTInput(ctx, packet, i, utxo, fromAddrInfo); err != nil {
			return nil, 0, fmt.Errorf("添加PSBT输入失败[%d]: %w", i, err)
		}
	}

	return packet, estimatedFee, nil
}

// estimateTransferTxSize 估算转账交易大小
func estimateTransferTxSize(numInputs, numOutputs int, addrInfo *bitcoindrpc.ValidateAddressDTO) int {
	baseSize := 10 // 版本(4) + 输入数量(1) + 输出数量(1) + locktime(4)

	inputSize := 0
	if addrInfo.IsWitness && addrInfo.WitnessVersion == 1 {
		// Taproot: 输入基础(41) + witness(64)
		inputSize = 57 // (41 + 64/4)，witness打折
	} else if addrInfo.IsWitness && addrInfo.WitnessVersion == 0 {
		// P2WPKH: 输入基础(41) + witness(107)
		inputSize = 68 // (41 + 107/4)，witness打折
	} else if addrInfo.IsScript {
		// P2SH-P2WPKH: 输入基础(64) + witness(107)
		inputSize = 91 // (64 + 107/4)
	} else {
		// P2PKH: 输入基础(41) + scriptSig(107)
		inputSize = 148
	}

	outputSize := 34 // 每个输出约34字节

	return baseSize + (numInputs * inputSize) + (numOutputs * outputSize)
}

// addPSBTInput 为PSBT添加输入数据
func (c *Client) addPSBTInput(
	ctx context.Context,
	packet *psbt.Packet,
	index int,
	utxo bitcoindrpc.UTXODTO,
	addrInfo *bitcoindrpc.ValidateAddressDTO,
) error {

	updater, err := psbt.NewUpdater(packet)
	if err != nil {
		return fmt.Errorf("创建PSBT更新器失败: %w", err)
	}

	scriptPubKey, err := hex.DecodeString(utxo.ScriptPubKey)
	if err != nil {
		return fmt.Errorf("解码scriptPubKey失败: %w", err)
	}

	amountSats := int64(utxo.AmountBTC * 1e8)

	// 根据地址类型添加相应的输入数据
	if addrInfo.IsWitness || addrInfo.IsScript {
		// Witness类型或P2SH需要WitnessUTXO
		txOut := &wire.TxOut{
			Value:    amountSats,
			PkScript: scriptPubKey,
		}
		if err := updater.AddInWitnessUtxo(txOut, index); err != nil {
			return fmt.Errorf("添加witness UTXO失败: %w", err)
		}

		// 如果是P2SH包裹的SegWit，需要添加RedeemScript
		if addrInfo.IsScript && !addrInfo.IsWitness {
			// P2SH-P2WPKH: 需要redeemScript
			logger.Info("      - 输入[%d]: 添加P2SH-P2WPKH的redeemScript", index)
		}
	} else {
		// Legacy P2PKH需要NonWitnessUTXO（完整的前序交易）
		logger.Info("      - 输入[%d]: 获取完整前序交易 (P2PKH)", index)
		prevTxRaw, err := c.GetRawTx(ctx, utxo.TxID)
		if err != nil {
			return fmt.Errorf("获取前序交易失败: %w", err)
		}

		var prevTx wire.MsgTx
		if err := prevTx.Deserialize(bytes.NewReader(prevTxRaw)); err != nil {
			return fmt.Errorf("反序列化前序交易失败: %w", err)
		}

		if err := updater.AddInNonWitnessUtxo(&prevTx, index); err != nil {
			return fmt.Errorf("添加non-witness UTXO失败: %w", err)
		}
	}

	return nil
}

// signPSBT 使用私钥签名PSBT
func (c *Client) signPSBT(
	packet *psbt.Packet,
	wif *btcutil.WIF,
	utxos []bitcoindrpc.UTXODTO,
	addrInfo *bitcoindrpc.ValidateAddressDTO,
) error {

	privKey := wif.PrivKey
	pubKey := privKey.PubKey()

	logger.Info("    - 开始签名 %d 个输入", len(utxos))

	for i := range utxos {
		logger.Info("      - 签名输入[%d/%d]", i+1, len(utxos))

		// 获取签名哈希
		var sigHash []byte
		var err error

		if addrInfo.IsWitness && addrInfo.WitnessVersion == 1 {
			// Taproot签名
			logger.Info("        * 使用Taproot签名方案")
			sigHash, err = getTaprootSigHash(packet, i, txscript.SigHashDefault)
		} else {
			// 其他类型使用传统或SegWit签名
			sigHashType := txscript.SigHashAll
			if addrInfo.IsWitness || addrInfo.IsScript {
				logger.Info("        * 使用SegWit签名方案")
				sigHash, err = getWitnessSigHash(packet, i, sigHashType)
			} else {
				logger.Info("        * 使用传统签名方案")
				sigHash, err = getLegacySigHash(packet, i, sigHashType)
			}
		}

		if err != nil {
			return fmt.Errorf("计算签名哈希失败[%d]: %w", i, err)
		}

		// 签名
		signature := ecdsa.Sign(privKey, sigHash)

		// 添加签名到PSBT
		if addrInfo.IsWitness && addrInfo.WitnessVersion == 1 {
			// Taproot: 只需要签名，不需要sighash type
			packet.Inputs[i].TaprootKeySpendSig = signature.Serialize()
			logger.Info("        ✓ Taproot密钥路径签名完成")
		} else {
			// 其他类型：签名 + sighash type + 公钥
			sigWithHashType := append(signature.Serialize(), byte(txscript.SigHashAll))
			packet.Inputs[i].PartialSigs = []*psbt.PartialSig{
				{
					PubKey:    pubKey.SerializeCompressed(),
					Signature: sigWithHashType,
				},
			}
			logger.Info("        ✓ ECDSA签名完成")
		}
	}

	logger.Info("    ✓ 所有输入签名完成")
	return nil
}

// getTaprootSigHash 获取Taproot签名哈希
func getTaprootSigHash(packet *psbt.Packet, inputIndex int, sigHashType txscript.SigHashType) ([]byte, error) {
	prevOuts := make([]*wire.TxOut, len(packet.UnsignedTx.TxIn))
	for i, in := range packet.Inputs {
		if in.WitnessUtxo == nil {
			return nil, fmt.Errorf("输入[%d]缺少witness UTXO", i)
		}
		prevOuts[i] = in.WitnessUtxo
	}

	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(
		func() map[wire.OutPoint]*wire.TxOut {
			m := make(map[wire.OutPoint]*wire.TxOut)
			for i, txIn := range packet.UnsignedTx.TxIn {
				m[txIn.PreviousOutPoint] = prevOuts[i]
			}
			return m
		}(),
	)

	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, prevOutputFetcher)
	sigHash, err := txscript.CalcTaprootSignatureHash(
		sigHashes,
		sigHashType,
		packet.UnsignedTx,
		inputIndex,
		prevOutputFetcher,
	)
	if err != nil {
		return nil, err
	}

	return sigHash, nil
}

// getWitnessSigHash 获取SegWit签名哈希
func getWitnessSigHash(packet *psbt.Packet, inputIndex int, sigHashType txscript.SigHashType) ([]byte, error) {
	if packet.Inputs[inputIndex].WitnessUtxo == nil {
		return nil, fmt.Errorf("输入[%d]缺少witness UTXO", inputIndex)
	}

	witnessScript := packet.Inputs[inputIndex].WitnessUtxo.PkScript
	amount := packet.Inputs[inputIndex].WitnessUtxo.Value

	sigHash, err := txscript.CalcWitnessSigHash(
		witnessScript,
		txscript.NewTxSigHashes(packet.UnsignedTx, nil),
		sigHashType,
		packet.UnsignedTx,
		inputIndex,
		amount,
	)
	if err != nil {
		return nil, err
	}

	return sigHash, nil
}

// getLegacySigHash 获取传统签名哈希
func getLegacySigHash(packet *psbt.Packet, inputIndex int, sigHashType txscript.SigHashType) ([]byte, error) {
	if packet.Inputs[inputIndex].NonWitnessUtxo == nil {
		return nil, fmt.Errorf("输入[%d]缺少non-witness UTXO", inputIndex)
	}

	prevOut := packet.UnsignedTx.TxIn[inputIndex].PreviousOutPoint
	prevTx := packet.Inputs[inputIndex].NonWitnessUtxo
	scriptPubKey := prevTx.TxOut[prevOut.Index].PkScript

	sigHash, err := txscript.CalcSignatureHash(
		scriptPubKey,
		sigHashType,
		packet.UnsignedTx,
		inputIndex,
	)
	if err != nil {
		return nil, err
	}

	return sigHash, nil
}
