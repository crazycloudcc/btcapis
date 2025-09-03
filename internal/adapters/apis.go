// 后续版本需要将以下两个Client合并为一个Client
package adapters

import (
	"context"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
)

// bitcoindrpc的接口集合
type CoreAdapterClient interface {

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// 地址相关接口 client_address.go

	// 查询钱包UTXO集（基于描述符/地址扫描全链 UTXO） - 会导致节点进行全量查询, 慎用, 等待时间很长
	AddressGetUTXOs(ctx context.Context, addr string) ([]bitcoindrpc.UTXODTO, error)

	// 查询钱包详细信息: 根据是否导入到本地节点, 返回数据不同
	AddressGetInfo(ctx context.Context, addr string) (*bitcoindrpc.AddressInfoDTO, error)

	// 校验钱包
	AddressValidate(ctx context.Context, addr string) (*bitcoindrpc.ValidateAddressDTO, error)

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// 链相关接口 client_chain.go

	// 估算交易费率
	ChainEstimateSmartFeeRate(ctx context.Context, targetBlocks int) (*bitcoindrpc.FeeRateSmartDTO, error)

	// 查询 UTXO
	ChainGetUTXO(ctx context.Context, hash [32]byte, index uint32) ([]byte, int64, error)

	// 获取节点区块数量
	ChainGetBlockCount(ctx context.Context) (int, error)

	// 获取最新区块的hash
	ChainGetBestBlockHash(ctx context.Context) (string, error)

	// 使用区块高度 查询区块哈希
	ChainGetBlockHash(ctx context.Context, height int64) (string, error)

	// 使用区块block hash 查询区块头
	ChainGetBlockHeader(ctx context.Context, hash string) (*bitcoindrpc.BlockHeaderDTO, error)

	// 使用区块block hash 查询区块
	ChainGetBlock(ctx context.Context, hash string) (*bitcoindrpc.BlockDTO, error)

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// 内存池相关接口 client_mempool.go

	// 获取内存池信息
	MempoolGetInfo(ctx context.Context) (*bitcoindrpc.MempoolInfoDTO, error)

	// 获取内存池交易信息
	MempoolGetTxs(ctx context.Context) ([]string, error)

	// 获取内存池交易信息
	MempoolGetTx(ctx context.Context, txid string) (*bitcoindrpc.MempoolTxDTO, error)

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// 节点状态相关接口 client_status.go

	// 获取节点网络信息
	GetNetworkInfo(ctx context.Context) (*bitcoindrpc.NetworkInfoDTO, error)

	// 获取链信息
	GetChainInfo(ctx context.Context) (*bitcoindrpc.ChainInfoDTO, error)

	// 获取区块统计信息
	GetBlockStats(ctx context.Context, height int64) (*bitcoindrpc.BlockStatsDTO, error)

	// 获取链顶信息
	GetChainTip(ctx context.Context) (*bitcoindrpc.ChainTipDTO, error)

	/////////////////////////////////////////////////////////////////////////////////////////////////
	// 交易相关接口 client_tx.go

	// 获取交易元数据
	TxGetRaw(ctx context.Context, txid string, decodeFlag bool) ([]byte, error)

	// 构建交易(taproot需要使用psbt)
	TxCreateRaw(ctx context.Context, dto *bitcoindrpc.TxCreateRawDTO) ([]byte, error)

	// 填充交易费用(taproot需要使用psbt)
	TxFundRaw(ctx context.Context, rawtx string, options *bitcoindrpc.TxFundOptionsDTO) (*bitcoindrpc.TxFundRawResultDTO, error)

	// 签名交易(taproot需要使用psbt)
	TxSignRawWithKey(ctx context.Context, rawtx string) (string, error)

	// 完成psbt交易
	TxFinalizePsbt(ctx context.Context, psbt string) (*bitcoindrpc.SignedTxDTO, error)

	// 广播交易
	TxBroadcast(ctx context.Context, rawtx []byte) (string, error)

	// 预检查交易 testmempoolaccept: 需要组装交易数据后生成hex字符串再测试
	TxTestMempoolAccept(ctx context.Context, rawtx []byte) (string, error)

	/////////////////////////////////////////////////////////////////////////////////////////////////
}
