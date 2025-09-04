package tx

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// 转账交易-PSBT预览: 通过输入数据根据发起转账钱包地址的类型创建对应的PSBT交易数据, 这个数据将提交给外部okx插件钱包等进行签名.
func (c *Client) CreatePSBT(ctx context.Context, inputParams *types.TxInputParams) (string, error) {
	ret, err := c.buildPSBT(ctx, inputParams)
	if err != nil {
		return "", err
	}
	return ret.PSBTBase64, nil
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
	fmt.Printf("psbt input len: %d, isHex: %v\n", len(normalized), isHex(normalized))

	// 十六进制转base64
	if isHex(normalized) {
		bin, err := hex.DecodeString(normalized)
		if err != nil {
			return "", err
		}
		psbtBase64 = base64.StdEncoding.EncodeToString(bin)
	} else {
		psbtBase64 = normalized
	}
	fmt.Printf("psbt base64 len: %d\n", len(psbtBase64))
	fmt.Printf("psbt base64: %s\n", psbtBase64)

	// finalizepsbt -> 原始交易hex
	finalizeData, err := c.bitcoindrpcClient.TxFinalizePsbt(ctx, psbtBase64)
	if err != nil || !finalizeData.Complete {
		fmt.Printf("finalizepsbt error: %v, complete: %v\n", err, finalizeData.Complete)
		return "", err
	}
	fmt.Println("finalize rawHex: ", finalizeData.Hex)

	bin, err := hex.DecodeString(finalizeData.Hex)
	if err != nil {
		return "", err
	}

	fmt.Println("broadcast tx len: ", len(bin))

	// 广播交易
	return c.bitcoindrpcClient.TxBroadcast(ctx, bin)
}

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
