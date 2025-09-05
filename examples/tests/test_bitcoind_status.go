// 内部测试节点是否正常, 接口不对外
package tests

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crazycloudcc/btcapis"
)

func TestBitcoindStatus(testClient *btcapis.TestClient) {
	fmt.Println("Testing bitcoind status...")

	// 测试 bitcoind rpc 是否正常
	netInfo, err := testClient.GetNetworkInfo(context.Background())
	if err != nil {
		panic(err)
	}
	netInfoJson, _ := json.MarshalIndent(netInfo, "", "  ")
	fmt.Printf("NetworkInfo: %s\n", string(netInfoJson))
	fmt.Println("--------------------------------")

	// 获取链信息
	chainInfo, err := testClient.GetBlockChainInfo(context.Background())
	if err != nil {
		panic(err)
	}
	chainInfoJson, _ := json.MarshalIndent(chainInfo, "", "  ")
	fmt.Printf("ChainInfo: %s\n", string(chainInfoJson))
	fmt.Println("--------------------------------")

	// 获取区块统计信息
	blockStats, err := testClient.GetBlockStats(context.Background(), int64(chainInfo.Blocks))
	if err != nil {
		panic(err)
	}
	blockStatsJson, _ := json.MarshalIndent(blockStats, "", "  ")
	fmt.Printf("BlockStats: %s\n", string(blockStatsJson))
	fmt.Println("--------------------------------")

	// 获取链顶信息
	chainTips, err := testClient.GetChainTips(context.Background())
	if err != nil {
		panic(err)
	}
	chainTipsJson, _ := json.MarshalIndent(chainTips, "", "  ")
	fmt.Printf("ChainTips: %s\n", string(chainTipsJson))
	fmt.Println("--------------------------------")

	fmt.Println("Test bitcoind status done.")
	fmt.Println("================================")
}
