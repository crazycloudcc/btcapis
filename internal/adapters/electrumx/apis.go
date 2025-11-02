// ElectrumX API接口实现
package electrumx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// ===== 地址相关接口 =====

// AddressGetBalance 获取地址余额
// 参数: addr - 比特币地址
// 返回: 已确认余额（聪）、未确认余额（聪）、错误
func (c *Client) AddressGetBalance(ctx context.Context, addr string) (int64, int64, error) {
	// 将地址转换为脚本哈希
	scriptHash, err := addressToScriptHash(addr)
	if err != nil {
		return 0, 0, fmt.Errorf("address to scripthash: %w", err)
	}

	var balance BalanceDTO
	if err := c.rpcCall(ctx, "blockchain.scripthash.get_balance", []interface{}{scriptHash}, &balance); err != nil {
		return 0, 0, err
	}

	return balance.Confirmed, balance.Unconfirmed, nil
}

// AddressGetHistory 获取地址交易历史
// 参数: addr - 比特币地址
// 返回: 交易历史列表、错误
func (c *Client) AddressGetHistory(ctx context.Context, addr string) ([]HistoryDTO, error) {
	// 将地址转换为脚本哈希
	scriptHash, err := addressToScriptHash(addr)
	if err != nil {
		return nil, fmt.Errorf("address to scripthash: %w", err)
	}

	var history []HistoryDTO
	if err := c.rpcCall(ctx, "blockchain.scripthash.get_history", []interface{}{scriptHash}, &history); err != nil {
		return nil, err
	}

	return history, nil
}

// AddressGetUTXOs 获取地址UTXO列表
// 参数: addr - 比特币地址
// 返回: UTXO列表、错误
func (c *Client) AddressGetUTXOs(ctx context.Context, addr string) ([]UTXODTO, error) {
	// 将地址转换为脚本哈希
	scriptHash, err := addressToScriptHash(addr)
	if err != nil {
		return nil, fmt.Errorf("address to scripthash: %w", err)
	}

	var utxos []UTXODTO
	if err := c.rpcCall(ctx, "blockchain.scripthash.listunspent", []interface{}{scriptHash}, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}

// AddressGetMempool 获取地址在内存池中的交易
// 参数: addr - 比特币地址
// 返回: 内存池交易列表、错误
func (c *Client) AddressGetMempool(ctx context.Context, addr string) ([]MempoolDTO, error) {
	// 将地址转换为脚本哈希
	scriptHash, err := addressToScriptHash(addr)
	if err != nil {
		return nil, fmt.Errorf("address to scripthash: %w", err)
	}

	var mempool []MempoolDTO
	if err := c.rpcCall(ctx, "blockchain.scripthash.get_mempool", []interface{}{scriptHash}, &mempool); err != nil {
		return nil, err
	}

	return mempool, nil
}

// AddressSubscribe 订阅地址变更通知
// 参数: addr - 比特币地址
// 返回: 当前状态哈希、错误
// 注意: 此方法需要WebSocket连接支持，HTTP连接仅返回当前状态
func (c *Client) AddressSubscribe(ctx context.Context, addr string) (string, error) {
	// 将地址转换为脚本哈希
	scriptHash, err := addressToScriptHash(addr)
	if err != nil {
		return "", fmt.Errorf("address to scripthash: %w", err)
	}

	var status string
	if err := c.rpcCall(ctx, "blockchain.scripthash.subscribe", []interface{}{scriptHash}, &status); err != nil {
		return "", err
	}

	return status, nil
}

// ===== 交易相关接口 =====

// TransactionGet 获取交易详情
// 参数: txid - 交易ID, verbose - 是否返回详细信息
// 返回: 交易原始十六进制或详细信息、错误
func (c *Client) TransactionGet(ctx context.Context, txid string, verbose bool) (interface{}, error) {
	var result interface{}
	if err := c.rpcCall(ctx, "blockchain.transaction.get", []interface{}{txid, verbose}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TransactionGetRaw 获取交易原始十六进制数据
// 参数: txid - 交易ID
// 返回: 交易原始十六进制、错误
func (c *Client) TransactionGetRaw(ctx context.Context, txid string) (string, error) {
	var rawTx string
	if err := c.rpcCall(ctx, "blockchain.transaction.get", []interface{}{txid, false}, &rawTx); err != nil {
		return "", err
	}
	return rawTx, nil
}

// TransactionBroadcast 广播交易
// 参数: rawTxHex - 交易原始十六进制数据
// 返回: 交易ID、错误
func (c *Client) TransactionBroadcast(ctx context.Context, rawTxHex string) (string, error) {
	var txid string
	if err := c.rpcCall(ctx, "blockchain.transaction.broadcast", []interface{}{rawTxHex}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// TransactionGetMerkle 获取交易的Merkle证明
// 参数: txid - 交易ID, height - 区块高度
// 返回: Merkle证明数据、错误
func (c *Client) TransactionGetMerkle(ctx context.Context, txid string, height int64) (interface{}, error) {
	var result interface{}
	if err := c.rpcCall(ctx, "blockchain.transaction.get_merkle", []interface{}{txid, height}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TransactionIDFromPos 根据区块高度和位置获取交易ID
// 参数: height - 区块高度, txPos - 交易在区块中的位置, merkle - 是否包含Merkle证明
// 返回: 交易ID或包含Merkle证明的数据、错误
func (c *Client) TransactionIDFromPos(ctx context.Context, height int64, txPos int64, merkle bool) (interface{}, error) {
	var result interface{}
	if err := c.rpcCall(ctx, "blockchain.transaction.id_from_pos", []interface{}{height, txPos, merkle}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ===== 区块相关接口 =====

// BlockchainGetBlockHeader 获取区块头信息
// 参数: height - 区块高度, cpHeight - checkpoint高度（可选）
// 返回: 区块头十六进制数据、错误
func (c *Client) BlockchainGetBlockHeader(ctx context.Context, height int64, cpHeight int64) (string, error) {
	var params []interface{}
	if cpHeight > 0 {
		params = []interface{}{height, cpHeight}
	} else {
		params = []interface{}{height}
	}

	var header string
	if err := c.rpcCall(ctx, "blockchain.block.header", params, &header); err != nil {
		return "", err
	}
	return header, nil
}

// BlockchainGetBlockHeaders 批量获取区块头信息
// 参数: startHeight - 起始区块高度, count - 获取数量, cpHeight - checkpoint高度（可选）
// 返回: 区块头数据、错误
func (c *Client) BlockchainGetBlockHeaders(ctx context.Context, startHeight int64, count int64, cpHeight int64) (interface{}, error) {
	var params []interface{}
	if cpHeight > 0 {
		params = []interface{}{startHeight, count, cpHeight}
	} else {
		params = []interface{}{startHeight, count}
	}

	var result interface{}
	if err := c.rpcCall(ctx, "blockchain.block.headers", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ===== 手续费相关接口 =====

// EstimateFee 估算交易手续费
// 参数: blocks - 目标确认区块数
// 返回: 手续费率（BTC/KB）、错误
func (c *Client) EstimateFee(ctx context.Context, blocks int) (float64, error) {
	var feeRate float64
	if err := c.rpcCall(ctx, "blockchain.estimatefee", []interface{}{blocks}, &feeRate); err != nil {
		return 0, err
	}
	return feeRate, nil
}

// RelayFee 获取中继手续费
// 返回: 最小中继手续费率（BTC/KB）、错误
func (c *Client) RelayFee(ctx context.Context) (float64, error) {
	var feeRate float64
	if err := c.rpcCall(ctx, "blockchain.relayfee", []interface{}{}, &feeRate); err != nil {
		return 0, err
	}
	return feeRate, nil
}

// ===== 服务器相关接口 =====

// ServerVersion 获取服务器版本信息
// 参数: clientName - 客户端名称, protocolVersion - 协议版本
// 返回: 服务器版本信息、错误
func (c *Client) ServerVersion(ctx context.Context, clientName string, protocolVersion string) (*ServerVersionDTO, error) {
	var result []string
	if err := c.rpcCall(ctx, "server.version", []interface{}{clientName, protocolVersion}, &result); err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid server version response")
	}

	return &ServerVersionDTO{
		ServerVersion:   result[0],
		ProtocolVersion: result[1],
	}, nil
}

// ServerFeatures 获取服务器功能信息
// 返回: 服务器功能信息、错误
func (c *Client) ServerFeatures(ctx context.Context) (*ServerFeaturesDTO, error) {
	var features ServerFeaturesDTO
	if err := c.rpcCall(ctx, "server.features", []interface{}{}, &features); err != nil {
		return nil, err
	}
	return &features, nil
}

// ServerPing 心跳检测
// 返回: 错误
func (c *Client) ServerPing(ctx context.Context) error {
	var result interface{}
	return c.rpcCall(ctx, "server.ping", []interface{}{}, &result)
}

// ServerBanner 获取服务器横幅信息
// 返回: 横幅文本、错误
func (c *Client) ServerBanner(ctx context.Context) (string, error) {
	var banner string
	if err := c.rpcCall(ctx, "server.banner", []interface{}{}, &banner); err != nil {
		return "", err
	}
	return banner, nil
}

// ===== 其他接口 =====

// GetBlockchainTip 获取当前区块链最新高度
// 返回: 最新区块高度、错误
func (c *Client) GetBlockchainTip(ctx context.Context) (int64, error) {
	// 使用 server.features 获取当前区块高度信息
	// 或使用 blockchain.headers.subscribe 获取最新区块头
	var result struct {
		Height int64  `json:"height"`
		Hex    string `json:"hex"`
	}
	if err := c.rpcCall(ctx, "blockchain.headers.subscribe", []interface{}{}, &result); err != nil {
		return 0, err
	}
	return result.Height, nil
}

// ===== 工具函数 =====

// addressToScriptHash 将比特币地址转换为ElectrumX使用的脚本哈希
// ElectrumX使用的是脚本的SHA256哈希的反序
func addressToScriptHash(addr string) (string, error) {
	// 使用decoders模块将地址转换为scriptPubKey
	pkScript, err := decoders.AddressToPkScript(addr)
	if err != nil {
		return "", fmt.Errorf("address to pkscript: %w", err)
	}

	// 计算SHA256哈希
	scriptHash := computeScriptHash(pkScript)
	return scriptHash, nil
}

// computeScriptHash 计算脚本的哈希值（ElectrumX格式）
// ElectrumX使用SHA256哈希的反序（little-endian）
func computeScriptHash(script []byte) string {
	// 引入crypto/sha256进行哈希计算
	hash := sha256Sum(script)

	// 反序字节序（ElectrumX使用little-endian）
	for i := 0; i < len(hash)/2; i++ {
		hash[i], hash[len(hash)-1-i] = hash[len(hash)-1-i], hash[i]
	}

	return hex.EncodeToString(hash)
}

// sha256Sum 计算SHA256哈希
func sha256Sum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// ===== 扩展功能 =====

// GetBalancesByXPRV 通过扩展私钥查询所有派生地址的余额
// 参数:
//
//	ctx - 上下文
//	xprv - 扩展私钥（HD钱包主密钥）
//	derivationPaths - 派生路径列表，例如 [][]uint32{{44', 0', 0', 0, 0}, {44', 0', 0', 0, 1}}
//	scanCount - 每种类型地址扫描的数量（如果为0，默认扫描20个）
//
// 返回: 地址余额信息列表、错误
func (c *Client) GetBalancesByXPRV(ctx context.Context, xprv string, scanCount uint32) ([]types.AddressBalanceInfo, error) {
	if scanCount == 0 {
		scanCount = 20 // 默认扫描20个地址
	}

	// 导入必要的包
	var results []types.AddressBalanceInfo

	// 解析扩展私钥
	master, err := parseExtendedKey(xprv)
	if err != nil {
		return nil, fmt.Errorf("parse extended key: %w", err)
	}

	// 派生四种类型的地址：P2PKH, P2SH-P2WPKH, P2WPKH, P2TR
	addressTypes := []struct {
		name string
		path []uint32
	}{
		{"P2PKH", []uint32{hardenedKey(44), hardenedKey(0), hardenedKey(0), 0}},
		{"P2SH-P2WPKH", []uint32{hardenedKey(49), hardenedKey(0), hardenedKey(0), 0}},
		{"P2WPKH", []uint32{hardenedKey(84), hardenedKey(0), hardenedKey(0), 0}},
		{"P2TR", []uint32{hardenedKey(86), hardenedKey(0), hardenedKey(0), 0}},
	}

	// 遍历每种地址类型
	for _, addrType := range addressTypes {
		// 遍历索引
		for i := uint32(0); i < scanCount; i++ {
			// 构建完整路径
			fullPath := append(addrType.path, i)

			// 派生地址
			address, err := deriveAddressFromPath(master, fullPath)
			if err != nil {
				results = append(results, types.AddressBalanceInfo{
					Address: fmt.Sprintf("%s[%d]", addrType.name, i),
					Error:   fmt.Errorf("derive address: %w", err),
				})
				continue
			}

			// 查询余额
			confirmed, unconfirmed, err := c.AddressGetBalance(ctx, address)
			results = append(results, types.AddressBalanceInfo{
				Address:     address,
				Confirmed:   confirmed,
				Unconfirmed: unconfirmed,
				Total:       confirmed + unconfirmed,
				Error:       err,
			})

			// 如果连续3个地址都没有余额，可以提前停止扫描该类型
			// （这是常见的"gap limit"策略）
			if i >= 3 {
				noBalanceCount := 0
				if results[len(results)-1].Total == 0 {
					noBalanceCount++
				}
				if noBalanceCount >= 3 {
					break
				}
			}
		}
	}

	return results, nil
}

// GetBalancesByPrivateKey 通过单个私钥查询对应地址的余额
// 参数:
//
//	ctx - 上下文
//	privateKeyWIF - WIF格式的私钥
//
// 返回: 地址余额信息、错误
func (c *Client) GetBalancesByPrivateKey(ctx context.Context, privateKeyWIF string) (*types.AddressBalanceInfo, error) {
	// 解析WIF私钥
	address, err := privateKeyToAddress(privateKeyWIF)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	// 查询余额
	confirmed, unconfirmed, err := c.AddressGetBalance(ctx, address)
	if err != nil {
		return &types.AddressBalanceInfo{
			Address: address,
			Error:   err,
		}, err
	}

	return &types.AddressBalanceInfo{
		Address:     address,
		Confirmed:   confirmed,
		Unconfirmed: unconfirmed,
		Total:       confirmed + unconfirmed,
	}, nil
}

// FilterAddressesWithBalance 过滤出余额大于0的地址
// 参数:
//
//	ctx - 上下文
//	addresses - 地址列表
//	minBalance - 最小余额（聪），默认为0表示只要有余额即可
//	concurrent - 并发查询数，建议设置为10-50
//
// 返回: 余额大于指定值的地址列表、错误
func (c *Client) FilterAddressesWithBalance(ctx context.Context, addresses []string, minBalance int64, concurrent int) ([]types.AddressBalanceInfo, error) {
	if concurrent <= 0 {
		concurrent = 10 // 默认并发数
	}

	results := make([]types.AddressBalanceInfo, 0)
	resultsChan := make(chan types.AddressBalanceInfo, len(addresses))
	semaphore := make(chan struct{}, concurrent)
	errorsChan := make(chan error, len(addresses))

	// 并发查询所有地址
	for _, addr := range addresses {
		semaphore <- struct{}{} // 获取信号量
		go func(address string) {
			defer func() { <-semaphore }() // 释放信号量

			confirmed, unconfirmed, err := c.AddressGetBalance(ctx, address)
			total := confirmed + unconfirmed

			if err != nil {
				errorsChan <- fmt.Errorf("query %s: %w", address, err)
				return
			}

			// 只返回余额大于等于最小值的地址
			if total >= minBalance {
				resultsChan <- types.AddressBalanceInfo{
					Address:     address,
					Confirmed:   confirmed,
					Unconfirmed: unconfirmed,
					Total:       total,
				}
			}
		}(addr)
	}

	// 等待所有goroutine完成
	for i := 0; i < concurrent; i++ {
		semaphore <- struct{}{}
	}

	close(resultsChan)
	close(errorsChan)

	// 收集结果
	for result := range resultsChan {
		results = append(results, result)
	}

	// 收集错误（如果有）
	var errs []error
	for err := range errorsChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		// 返回第一个错误，但仍然返回已查询到的结果
		return results, errs[0]
	}

	return results, nil
}

// BatchGetBalances 批量查询地址余额（返回所有地址的余额，包括0余额）
// 参数:
//
//	ctx - 上下文
//	addresses - 地址列表
//	concurrent - 并发查询数
//
// 返回: 所有地址的余额信息列表、错误
func (c *Client) BatchGetBalances(ctx context.Context, addresses []string, concurrent int) ([]types.AddressBalanceInfo, error) {
	if concurrent <= 0 {
		concurrent = 10
	}

	results := make([]types.AddressBalanceInfo, len(addresses))
	semaphore := make(chan struct{}, concurrent)

	// 使用通道来保持顺序
	type indexedResult struct {
		index int
		info  types.AddressBalanceInfo
	}
	resultsChan := make(chan indexedResult, len(addresses))

	// 并发查询
	for i, addr := range addresses {
		semaphore <- struct{}{}
		go func(index int, address string) {
			defer func() { <-semaphore }()

			confirmed, unconfirmed, err := c.AddressGetBalance(ctx, address)
			resultsChan <- indexedResult{
				index: index,
				info: types.AddressBalanceInfo{
					Address:     address,
					Confirmed:   confirmed,
					Unconfirmed: unconfirmed,
					Total:       confirmed + unconfirmed,
					Error:       err,
				},
			}
		}(i, addr)
	}

	// 等待所有查询完成
	for i := 0; i < concurrent; i++ {
		semaphore <- struct{}{}
	}
	close(resultsChan)

	// 按原始顺序组织结果
	for result := range resultsChan {
		results[result.index] = result.info
	}

	return results, nil
}

// ===== 辅助函数 =====

// hardenedKey 返回硬化密钥索引
func hardenedKey(index uint32) uint32 {
	return hdkeychain.HardenedKeyStart + index
}

// parseExtendedKey 解析扩展密钥（xprv或xpub）
func parseExtendedKey(key string) (*hdkeychain.ExtendedKey, error) {
	return hdkeychain.NewKeyFromString(key)
}

// deriveAddressFromPath 从主密钥和路径派生地址
func deriveAddressFromPath(master *hdkeychain.ExtendedKey, path []uint32) (string, error) {
	// 派生到指定路径
	key := master
	for _, index := range path {
		var err error
		key, err = key.Derive(index)
		if err != nil {
			return "", fmt.Errorf("derive index %d: %w", index, err)
		}
	}

	// 获取网络参数
	params := types.CurrentNetworkParams

	// 根据路径判断地址类型
	if len(path) >= 1 {
		// 获取purpose（第一个索引去掉hardened bit）
		purpose := path[0] - hdkeychain.HardenedKeyStart

		switch purpose {
		case 44: // P2PKH
			return deriveP2PKHAddress(key, params)
		case 49: // P2SH-P2WPKH
			return deriveP2SHAddress(key, params)
		case 84: // P2WPKH
			return deriveP2WPKHAddress(key, params)
		case 86: // P2TR (Taproot)
			return deriveP2TRAddress(key, params)
		}
	}

	// 默认使用P2WPKH
	return deriveP2WPKHAddress(key, params)
}

// deriveP2PKHAddress 派生P2PKH地址（1...）
func deriveP2PKHAddress(key *hdkeychain.ExtendedKey, params *chaincfg.Params) (string, error) {
	pub, err := key.Neuter()
	if err != nil {
		return "", err
	}
	pk, err := pub.ECPubKey()
	if err != nil {
		return "", err
	}
	pkHash := btcutil.Hash160(pk.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// deriveP2SHAddress 派生P2SH-P2WPKH地址（3...）
func deriveP2SHAddress(key *hdkeychain.ExtendedKey, params *chaincfg.Params) (string, error) {
	pub, err := key.Neuter()
	if err != nil {
		return "", err
	}
	pk, err := pub.ECPubKey()
	if err != nil {
		return "", err
	}

	// 构建witness程序：OP_0 <20-byte-pubkey-hash>
	pkHash := btcutil.Hash160(pk.SerializeCompressed())
	witnessProgram := make([]byte, 0, 22)
	witnessProgram = append(witnessProgram, txscript.OP_0, 0x14)
	witnessProgram = append(witnessProgram, pkHash...)

	// 计算脚本哈希
	scriptHash := btcutil.Hash160(witnessProgram)
	addr, err := btcutil.NewAddressScriptHashFromHash(scriptHash, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// deriveP2WPKHAddress 派生P2WPKH地址（bc1q...）
func deriveP2WPKHAddress(key *hdkeychain.ExtendedKey, params *chaincfg.Params) (string, error) {
	pub, err := key.Neuter()
	if err != nil {
		return "", err
	}
	pk, err := pub.ECPubKey()
	if err != nil {
		return "", err
	}
	pkHash := btcutil.Hash160(pk.SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// deriveP2TRAddress 派生P2TR地址（bc1p...）
func deriveP2TRAddress(key *hdkeychain.ExtendedKey, params *chaincfg.Params) (string, error) {
	// 获取私钥
	privKey, err := key.ECPrivKey()
	if err != nil {
		return "", err
	}

	// 获取内部公钥
	internalPubKey := privKey.PubKey()

	// 计算tweaked公钥（Taproot）
	tweakedKey := txscript.ComputeTaprootKeyNoScript(internalPubKey)

	// 获取x-only公钥（32字节）
	xOnlyPubKey := tweakedKey.X().Bytes()

	// 创建Taproot地址
	addr, err := btcutil.NewAddressTaproot(xOnlyPubKey, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// privateKeyToAddress 从WIF格式私钥派生地址
func privateKeyToAddress(wif string) (string, error) {
	// 解析WIF私钥
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return "", fmt.Errorf("decode WIF: %w", err)
	}

	// 获取公钥
	pubKey := w.PrivKey.PubKey()

	// 根据是否压缩判断地址类型
	var addr btcutil.Address
	if w.CompressPubKey {
		// 压缩公钥 -> P2WPKH (bc1q...)
		pkHash := btcutil.Hash160(pubKey.SerializeCompressed())
		addr, err = btcutil.NewAddressWitnessPubKeyHash(pkHash, types.CurrentNetworkParams)
	} else {
		// 非压缩公钥 -> P2PKH (1...)
		pkHash := btcutil.Hash160(pubKey.SerializeUncompressed())
		addr, err = btcutil.NewAddressPubKeyHash(pkHash, types.CurrentNetworkParams)
	}

	if err != nil {
		return "", fmt.Errorf("create address: %w", err)
	}

	return addr.EncodeAddress(), nil
}
