package tests

import (
	"context"
	"fmt"
	"log"

	"github.com/crazycloudcc/btcapis"
)

func TestMempoolSpaceStatus(client *btcapis.Client) {
	fmt.Println("Testing mempool space status...")

	// 测试 mempool apis 是否正常
	feerate1, feerate2, err := client.EstimateFeeRate(context.Background(), 1)
	if err != nil {
		log.Fatalf("EstimateFeeRate: %v", err)
	}
	fmt.Printf("feerate1: %.2f (sats/vB)\n", feerate1*1e8/1000.0)
	fmt.Printf("feerate2: %.2f (sats/vB)\n", feerate2)
	fmt.Println("--------------------------------")

	fmt.Println("Test mempool space status done.")
	fmt.Println("================================")
}
