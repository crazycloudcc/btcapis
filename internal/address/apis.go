package address

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/crazycloudcc/btcapis/types"
)

// GetAddressBalanceWithElectrumX 通过ElectrumX获取地址的余额
func (c *Client) GetAddressBalanceWithElectrumX(ctx context.Context, addr string) (float64, float64, error) {
	if c.electrumxClient != nil {
		confirmedSats, unconfirmedSats, err := c.electrumxClient.AddressGetBalance(ctx, addr)
		if err != nil {
			return 0, 0, err
		}
		return float64(confirmedSats), float64(unconfirmedSats), nil
	}
	return 0, 0, errors.New("btcapis: no client available")
}

// 通过扩展私钥查询余额
func (c *Client) GetAddressBalanceWithElectrumXByXPRV(ctx context.Context, xprv string, numAddresses uint32) ([]types.AddressBalanceInfo, error) {
	if c.electrumxClient != nil {
		return c.electrumxClient.GetBalancesByXPRV(ctx, xprv, numAddresses)
	}
	return nil, errors.New("btcapis: not implemented yet")
}

// 通过私钥查询余额
func (c *Client) GetAddressBalanceWithElectrumXByPrivateKey(ctx context.Context, privateKeyWIF string) (*types.AddressBalanceInfo, error) {
	if c.electrumxClient != nil {
		return c.electrumxClient.GetBalancesByPrivateKey(ctx, privateKeyWIF)
	}
	return nil, errors.New("btcapis: not implemented yet")
}

// 批量查询并过滤有余额的地址
func (c *Client) FilterAddressesWithBalanceWithElectrumX(ctx context.Context, addresses []string, concurrent int) ([]types.AddressBalanceInfo, error) {
	if c.electrumxClient != nil {
		return c.electrumxClient.FilterAddressesWithBalance(ctx, addresses, 0, concurrent)
	}
	return nil, errors.New("btcapis: no electrumx client available")
}

// 批量查询所有地址余额
func (c *Client) BatchGetBalancesWithElectrumX(ctx context.Context, addresses []string, concurrent int) ([]types.AddressBalanceInfo, error) {
	if c.electrumxClient != nil {
		return c.electrumxClient.BatchGetBalances(ctx, addresses, concurrent)
	}
	return nil, errors.New("btcapis: no electrumx client available")
}

// // GetAddressUTXOsWithElectrumX 通过ElectrumX获取地址的UTXO
// func (c *Client) GetAddressUTXOsWithElectrumX(ctx context.Context, addr string) ([]types.TxUTXO, error) {
// 	if c.electrumxClient != nil {
// 		UTXODTOs, err := c.electrumxClient.AddressGetUTXOs(ctx, addr)
// 		if err != nil {
// 			return nil, err
// 		}

// 		utxos := make([]types.TxUTXO, 0, len(UTXODTOs))
// 		for _, dto := range UTXODTOs {
// 			txidBytes, _ := hex.DecodeString(dto.TxHash)
// 			u := types.TxUTXO{
// 				OutPoint: types.TxOutPoint{Hash: types.Hash32(txidBytes), Index: dto.TxPos},
// 				Value:    dto.Value,
// 			}
// 			if dto.Height > 0 {
// 				u.Height = uint32(dto.Height)
// 			}
// 			utxos = append(utxos, u)
// 		}
// 		return utxos, nil
// 	}
// 	return nil, errors.New("btcapis: no electrumx client available")
// }

// GetAddressBalance 通过地址, 获取地址的确认余额和未确认余额.
func (c *Client) GetAddressBalance(ctx context.Context, addr string) (float64, float64, error) {
	if c.mempoolapisClient != nil {
		return c.mempoolapisClient.AddressGetBalance(ctx, addr)
	}
	return 0, 0, errors.New("btcapis: no client available")
}

// GetAddressUTXOs 通过地址, 获取地址拥有的UTXO.
func (c *Client) GetAddressUTXOs(ctx context.Context, addr string) ([]types.TxUTXO, error) {
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

	if c.mempoolapisClient != nil {
		UTXODTOs, err := c.mempoolapisClient.AddressGetUTXOs(ctx, addr)
		if UTXODTOs != nil && err == nil {
			utxos := make([]types.TxUTXO, 0, len(UTXODTOs))
			for _, dto := range UTXODTOs {
				txidBytes, _ := hex.DecodeString(dto.Txid)
				u := types.TxUTXO{
					OutPoint: types.TxOutPoint{Hash: types.Hash32(txidBytes), Index: dto.Vout},
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
