// ElectrumX API接口实现
package electrumx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/decoders"
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
