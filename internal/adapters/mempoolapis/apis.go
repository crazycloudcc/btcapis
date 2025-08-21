package mempoolapis

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/crazycloudcc/btcapis/internal/types"
)

func GetRawTx(ctx context.Context, txid string) ([]byte, error) {
	u := *config.base
	u.Path = path.Join(u.Path, "/api/tx/", txid, "hex")
	b, err := getBytes(ctx, u.String())
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(b)))
}

func Broadcast(ctx context.Context, rawtx []byte) (string, error) {
	// mempool.space 支持 POST /api/tx，body 为 hex
	u := *config.base
	u.Path = path.Join(u.Path, "/api/tx")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(hex.EncodeToString(rawtx)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := config.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("mempool POST /api/tx status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	txid, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(txid)), nil
}

// 获取地址余额
func GetAddressBalance(ctx context.Context, addr string) (int64, int64, error) {
	u := *config.base
	u.Path = path.Join(u.Path, "/api/address/", addr)
	var dto struct {
		ChainStats struct {
			Funded int64 `json:"funded_txo_sum"`
			Spent  int64 `json:"spent_txo_sum"`
		} `json:"chain_stats"`
		MempoolStats struct {
			Funded int64 `json:"funded_txo_sum"`
			Spent  int64 `json:"spent_txo_sum"`
		} `json:"mempool_stats"`
	}
	if err := getJSON(ctx, u.String(), &dto); err != nil {
		return 0, 0, err
	}
	confirmed := dto.ChainStats.Funded - dto.ChainStats.Spent
	mempool := dto.MempoolStats.Funded - dto.MempoolStats.Spent
	return confirmed, mempool, nil
}

// 获取地址 UTXO
func GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	u := *config.base
	u.Path = path.Join(u.Path, "/api/address/", addr, "/utxo")
	var dtos []UTXODTO
	if err := getJSON(ctx, u.String(), &dtos); err != nil {
		return nil, err
	}
	utxos := make([]types.UTXO, 0, len(dtos))
	for _, d := range dtos {
		txidBytes, _ := hex.DecodeString(d.Txid)
		u := types.UTXO{
			OutPoint: types.OutPoint{Hash: types.Hash32(txidBytes), Index: d.Vout},
			Value:    d.Value,
		}
		if d.Status.Confirmed {
			u.Height = uint32(d.Status.BlockHeight)
		}
		utxos = append(utxos, u)
	}
	return utxos, nil
}
