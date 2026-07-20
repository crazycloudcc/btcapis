package chain

import (
	"context"
	"errors"
	"fmt"

	"github.com/crazycloudcc/btcapis/internal/adapters/bitcoindrpc"
	"github.com/crazycloudcc/btcapis/types"
)

var ErrBitcoindUnavailable = errors.New("bitcoind rpc unavailable")

func (c *Client) requireBitcoind() (*bitcoindrpc.Client, error) {
	if c == nil || c.bitcoindrpcClient == nil {
		return nil, ErrBitcoindUnavailable
	}
	return c.bitcoindrpcClient, nil
}

// GetNetworkInfo 获取节点网络信息
func (c *Client) GetNetworkInfo(ctx context.Context) (*types.ChainNetworkInfo, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.GetNetworkInfo(ctx)
	if err != nil {
		return nil, err
	}
	return networkInfoFromDTO(dto), nil
}

// GetBlockStats 获取区块统计（无高度参数，与 apis 历史行为一致）
func (c *Client) GetBlockStats(ctx context.Context) (*types.ChainStats, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.GetBlockStatsDefault(ctx)
	if err != nil {
		return nil, err
	}
	return blockStatsFromDTO(dto), nil
}

// GetChainTips 获取链分叉信息
func (c *Client) GetChainTips(ctx context.Context) ([]types.ChainTip, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dtos, err := rpc.GetChainTips(ctx)
	if err != nil {
		return nil, err
	}
	return chainTipsFromDTO(dtos), nil
}

// GetBlockChainInfo 获取链状态
func (c *Client) GetBlockChainInfo(ctx context.Context) (*types.BlockChainInfo, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}
	return blockChainInfoFromDTO(dto), nil
}

// GetBlockHeader 使用区块 hash 查询区块头
func (c *Client) GetBlockHeader(ctx context.Context, blockHash string) (*types.BlockHeader, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.ChainGetBlockHeader(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	return blockHeaderFromDTO(dto), nil
}

// GetBlock 使用区块 hash 查询区块（verbosity=1，仅 txid 列表）
func (c *Client) GetBlock(ctx context.Context, blockHash string) (*types.Block, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.ChainGetBlock(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	return blockFromDTO(dto), nil
}

// GetBlockHashByHeight 使用区块高度查询 hash
func (c *Client) GetBlockHashByHeight(ctx context.Context, height int64) (string, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return "", err
	}
	return rpc.ChainGetBlockHash(ctx, height)
}

// ValidateAddress 校验地址
func (c *Client) ValidateAddress(ctx context.Context, addr string) (*types.AddressValidation, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.AddressValidate(ctx, addr)
	if err != nil {
		return nil, err
	}
	return addressValidationFromDTO(dto), nil
}

// TestMempoolAccept 预检查交易是否可被内存池接受
func (c *Client) TestMempoolAccept(ctx context.Context, rawtx []byte) (string, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return "", err
	}
	return rpc.TxTestMempoolAccept(ctx, rawtx)
}

// GetMempoolTxIds 拉取内存池交易 txid 列表
func (c *Client) GetMempoolTxIds(ctx context.Context) ([]string, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	return rpc.MempoolGetTxs(ctx)
}

// GetMempoolTxEntry 获取内存池交易详情
func (c *Client) GetMempoolTxEntry(ctx context.Context, txid string) (*types.MempoolTxEntry, error) {
	rpc, err := c.requireBitcoind()
	if err != nil {
		return nil, err
	}
	dto, err := rpc.MempoolGetTx(ctx, txid)
	if err != nil {
		return nil, err
	}
	return mempoolTxFromDTO(dto), nil
}

func networkInfoFromDTO(dto *bitcoindrpc.NetworkInfoDTO) *types.ChainNetworkInfo {
	if dto == nil {
		return nil
	}
	networks := make([]types.ChainNetworkMeta, 0, len(dto.Networks))
	for _, n := range dto.Networks {
		networks = append(networks, types.ChainNetworkMeta{
			Name:                      n.Name,
			Limited:                   n.Limited,
			Reachable:                 n.Reachable,
			Proxy:                     n.Proxy,
			ProxyRandomizeCredentials: n.ProxyRandomizeCredentials,
		})
	}
	localAddrs := make([]string, 0, len(dto.Localaddresses))
	for _, la := range dto.Localaddresses {
		localAddrs = append(localAddrs, fmt.Sprintf("%s:%d", la.Address, la.Port))
	}
	var warnings any
	if dto.Warnings != nil {
		warnings = dto.Warnings.AsAny()
	}
	return &types.ChainNetworkInfo{
		Version:            int64(dto.Version),
		Subversion:         dto.Subversion,
		ProtocolVersion:    int64(dto.ProtocolVersion),
		LocalServices:      dto.LocalServices,
		LocalServicesNames: dto.LocalServicesNames,
		LocalRelay:         dto.LocalRelay,
		TimeOffset:         int64(dto.TimeOffset),
		NetworkActive:      dto.NetworkActive,
		Connections:        int64(dto.Connections),
		ConnectionsIn:      int64(dto.Connectionsin),
		ConnectionsOut:     int64(dto.Connectionsout),
		Networks:           networks,
		RelayFee:           dto.Relayfee,
		IncrementalFee:     dto.Incrementalfee,
		LocalAddresses:     localAddrs,
		Warnings:           warnings,
	}
}

func blockStatsFromDTO(dto *bitcoindrpc.BlockStatsDTO) *types.ChainStats {
	if dto == nil {
		return nil
	}
	percentiles := make([]float64, len(dto.Feeratepercentiles))
	for i, v := range dto.Feeratepercentiles {
		percentiles[i] = float64(v)
	}
	return &types.ChainStats{
		AvgFee:             float64(dto.Avgfee),
		AvgFeeRate:         float64(dto.Avgfeerate),
		AvgTxSize:          float64(dto.Avgtxsize),
		BlockHash:          dto.Blockhash,
		FeeratePercentiles: percentiles,
		Height:             int64(dto.Height),
		Ins:                int64(dto.Ins),
		MaxFee:             float64(dto.Maxfee),
		MaxFeeRate:         float64(dto.Maxfeerate),
		MaxTxSize:          int64(dto.Maxtxsize),
		MedianFee:          float64(dto.Medianfee),
		MedianTime:         int64(dto.Mediantime),
		MedianTxSize:       float64(dto.Mediantxsize),
		MinFee:             float64(dto.Minfee),
		MinFeeRate:         float64(dto.Minfeerate),
		MinTxSize:          int64(dto.Mintxsize),
		Outs:               int64(dto.Outs),
		Subsidy:            float64(dto.Subsidy),
		SwtotalSize:        int64(dto.SwtotalSize),
		SwtotalWeight:      int64(dto.SwtotalWeight),
		Swtxs:              int64(dto.Swtxs),
		Time:               int64(dto.Time),
		TotalOut:           float64(dto.TotalOut),
		TotalSize:          int64(dto.TotalSize),
		TotalWeight:        int64(dto.TotalWeight),
		Totalfee:           float64(dto.Totalfee),
		Txs:                int64(dto.Txs),
		UTXOIncrease:       int64(dto.UtxoIncrease),
		UTXOSizeInc:        int64(dto.UtxoSizeInc),
		UTXOIncreaseActual: int64(dto.UtxoIncreaseActual),
		UTXOSizeIncActual:  int64(dto.UtxoSizeIncActual),
	}
}

func chainTipsFromDTO(dtos []bitcoindrpc.ChainTipDTO) []types.ChainTip {
	out := make([]types.ChainTip, len(dtos))
	for i, dto := range dtos {
		out[i] = types.ChainTip{
			Height:    dto.Height,
			Hash:      dto.Hash,
			BranchLen: int64(dto.BranchLen),
			Status:    dto.Status,
		}
	}
	return out
}

func blockChainInfoFromDTO(dto *bitcoindrpc.ChainInfoDTO) *types.BlockChainInfo {
	if dto == nil {
		return nil
	}
	var warnings any
	if dto.Warnings != nil {
		warnings = dto.Warnings.AsAny()
	}
	return &types.BlockChainInfo{
		Chain:                dto.Chain,
		Blocks:               int64(dto.Blocks),
		Headers:              int64(dto.Headers),
		BestBlockHash:        dto.Bestblockhash,
		Bits:                 dto.Bits,
		Target:               dto.Target,
		Difficulty:           dto.Difficulty,
		Time:                 int64(dto.Time),
		MedianTime:           int64(dto.MedianTime),
		VerificationProgress: dto.Verificationprogress,
		InitialBlockDownload: dto.Initialblockdownload,
		ChainWork:            dto.Chainwork,
		SizeOnDisk:           int64(dto.Sizeondisk),
		Pruned:               dto.Pruned,
		Warnings:             warnings,
	}
}

func blockHeaderFromDTO(dto *bitcoindrpc.BlockDTO) *types.BlockHeader {
	if dto == nil {
		return nil
	}
	return &types.BlockHeader{
		Hash:              dto.Hash,
		Confirmations:     int64(dto.Confirmations),
		Height:            int64(dto.Height),
		Version:           int64(dto.Version),
		VersionHex:        dto.VersionHex,
		MerkleRoot:        dto.MerkleRoot,
		Time:              int64(dto.Time),
		MedianTime:        int64(dto.MedianTime),
		Nonce:             int64(dto.Nonce),
		Bits:              dto.Bits,
		Difficulty:        dto.Difficulty,
		ChainWork:         dto.Chainwork,
		NTx:               int64(dto.NTx),
		PreviousBlockHash: dto.PreviousBlockHash,
		NextBlockHash:     dto.NextBlockHash,
	}
}

func blockFromDTO(dto *bitcoindrpc.BlockDTO) *types.Block {
	if dto == nil {
		return nil
	}
	header := blockHeaderFromDTO(dto)
	return &types.Block{
		BlockHeader:  *header,
		StrippedSize: int64(dto.StrippedSize),
		Size:         int64(dto.Size),
		Weight:       int64(dto.Weight),
		Tx:           dto.Tx,
	}
}

func addressValidationFromDTO(dto *bitcoindrpc.ValidateAddressDTO) *types.AddressValidation {
	if dto == nil {
		return nil
	}
	return &types.AddressValidation{
		IsValid:        dto.IsValid,
		Address:        dto.Address,
		ScriptPubKey:   dto.ScriptPubKey,
		IsScript:       dto.IsScript,
		IsWitness:      dto.IsWitness,
		WitnessVersion: int64(dto.WitnessVersion),
		WitnessProgram: dto.WitnessProgram,
	}
}

func mempoolTxFromDTO(dto *bitcoindrpc.MempoolTxDTO) *types.MempoolTxEntry {
	if dto == nil {
		return nil
	}
	return &types.MempoolTxEntry{
		Vsize:             int64(dto.Vsize),
		Weight:            int64(dto.Weight),
		Time:              int64(dto.Time),
		Height:            int64(dto.Height),
		DescendantCount:   int64(dto.DescendantCount),
		DescendantSize:    int64(dto.DescendantSize),
		AncestorCount:     int64(dto.AncestorCount),
		AncestorSize:      int64(dto.AncestorSize),
		Wtxid:             dto.Wtxid,
		Fees:              types.MempoolTxFees{Base: dto.Fees.Base, Modified: dto.Fees.Modified, Ancestor: dto.Fees.Ancestor, Descendant: dto.Fees.Descendant},
		Depends:           dto.Depends,
		SpentBy:           dto.Spentby,
		Bip125Replaceable: dto.Bip125Replaceable,
		Unbroadcast:       dto.Unbroadcast,
	}
}
