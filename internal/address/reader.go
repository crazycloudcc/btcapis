package address

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/crazycloudcc/btcapis/internal/adapters/mempoolapis"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/types"
)

// GetAddressScriptInfo 通过地址, 获取地址的锁定脚本信息.
func GetAddressScriptInfo(ctx context.Context, addr string) (*types.AddressScriptInfo, error) {
	scriptInfo, err := decoders.DecodeAddress(addr)
	if err != nil {
		return nil, err
	}
	return scriptInfo, nil
}

// GetAddressInfo 通过锁定脚本, 获取地址信息.
func GetAddressInfo(ctx context.Context, pkScript []byte) (*types.AddressInfo, error) {
	scriptInfo, err := decoders.DecodePkScript(pkScript)
	if err != nil {
		return nil, err
	}
	return &types.AddressInfo{
		PKScript:  pkScript,
		Typ:       scriptInfo.Typ,
		ReqSigs:   scriptInfo.ReqSigs,
		Addresses: scriptInfo.Addresses,
	}, nil
}

// GetAddressBalance 通过地址, 获取地址的确认余额和未确认余额.
func GetAddressBalance(ctx context.Context, addr string) (confirmed int64, mempool int64, err error) {
	// if bitcoindrpc.IsInited() {
	// 	return bitcoindrpc.GetAddressBalance(ctx, addr)
	// }
	if mempoolapis.IsInited() {
		return mempoolapis.GetAddressBalance(ctx, addr)
	}
	return 0, 0, errors.New("btcapis: no client available")
}

// GetAddressUTXOs 通过地址, 获取地址拥有的UTXO.
func GetAddressUTXOs(ctx context.Context, addr string) ([]types.UTXO, error) {
	errRet := errors.New("btcapis: no client available or no utxos")

	// 全量扫UTXO耗时太长, 暂时使用mempool.space的API
	// if bitcoindrpc.IsInited() {
	// 	UTXODTOs, err := bitcoindrpc.GetAddressUTXOs(ctx, addr)
	// 	if UTXODTOs != nil && err == nil {
	// 		utxos := make([]types.UTXO, 0, len(UTXODTOs))
	// 		for _, dto := range UTXODTOs {
	// 			txidBytes, _ := hex.DecodeString(dto.TxID)
	// 			utxos = append(utxos, types.UTXO{
	// 				OutPoint: types.OutPoint{Hash: types.Hash32(txidBytes), Index: dto.Vout},
	// 				Value:    int64(dto.AmountBTC * 1e8),
	// 			})
	// 		}
	// 		return utxos, nil
	// 	}
	// 	errRet = err
	// }

	if mempoolapis.IsInited() {
		UTXODTOs, err := mempoolapis.GetAddressUTXOs(ctx, addr)
		if UTXODTOs != nil && err == nil {
			utxos := make([]types.UTXO, 0, len(UTXODTOs))
			for _, dto := range UTXODTOs {
				txidBytes, _ := hex.DecodeString(dto.Txid)
				u := types.UTXO{
					OutPoint: types.OutPoint{Hash: types.Hash32(txidBytes), Index: dto.Vout},
					Value:    dto.Value,
				}
				if dto.Status.Confirmed {
					u.Height = uint32(dto.Status.BlockHeight)
				}
				utxos = append(utxos, u)
			}
			return utxos, nil
		}
	}

	return nil, errRet
}
