package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
	"github.com/crazycloudcc/btcapis/examples/tests"
)

// const (
//
//	network = "mainnet"
//	rpcUser = "cc"
//	rpcPass = "ccc"
//	rpcUrl  = "http://192.168.1.16:8332"
//	timeout = 10 * time.Second
//
// )
const (
	network = "signet"
	rpcUser = "cc"
	rpcPass = "ccc"
	rpcUrl  = "http://localhost:18332"
	timeout = 10 * time.Second
)

func main() {
	client := btcapis.New(
		network,
		rpcUrl,
		rpcUser,
		rpcPass,
		int(timeout.Seconds()),
	)

	if client == nil {
		log.Fatalf("New: %v", errors.New("client is nil"))
	}

	testClient := btcapis.NewTestClient(client)
	if testClient == nil {
		log.Fatalf("NewTestClient: %v", errors.New("testClient is nil"))
	}

	fmt.Println("test all start...")
	testNodesOnly(client, testClient)
	// tests.TestBitcoindStatus(testClient)
	// tests.TestMempoolSpaceStatus(client)

	// tests.TestAddress(client)
	// tests.TestTxs(client)
	tests.TestScripts(client)
	fmt.Println("test all done.")
	fmt.Println("================================")
}

// 测试节点是否正常
func testNodesOnly(client *btcapis.Client, testClient *btcapis.TestClient) {
	fmt.Println("Testing nodes...")

	// 测试 bitcoind rpc 是否正常
	netInfo, err := testClient.GetNetworkInfo(context.Background())
	if err != nil {
		log.Fatalf("GetNetworkInfo: %v", err)
	}

	netJson, _ := json.MarshalIndent(netInfo, "", "  ")
	fmt.Printf("NetworkInfo: %s\n", string(netJson))
	fmt.Println("--------------------------------")

	// 测试 mempool apis 是否正常
	feerate1, feerate2, err := client.EstimateFeeRate(context.Background(), 1)
	if err != nil {
		log.Fatalf("EstimateFeeRate: %v", err)
	}
	fmt.Printf("feerate1: %.2f (sats/vB)\n", feerate1*1e8/1000.0)
	fmt.Printf("feerate2: %.2f (sats/vB)\n", feerate2)
	fmt.Println("--------------------------------")

	fmt.Println("Test nodes done.")
	fmt.Println("================================")
}
