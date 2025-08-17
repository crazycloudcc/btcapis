// 通过交易元数据, 解析交易的详细信息并组装新的数据结构返回.
package btcdapis

import "github.com/crazycloudcc/btcapis/types"

func ParseTx(tx *types.Tx) (*types.Tx, error) {
	return tx, nil
}
