package tests

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crazycloudcc/btcapis"
	"github.com/crazycloudcc/btcapis/internal/utils"
)

func TestAddress(client *btcapis.Client) {
	fmt.Println("Testing address...")

	addr := getAddr()

	// 查询地址余额
	confirmed, mempool, err := client.GetAddressBalance(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Balance(BTC): %.8f(%.8f)\n", utils.SatsToBTC(confirmed), utils.SatsToBTC(mempool))
	fmt.Println("--------------------------------")

	// 查询地址UTXO
	utxos, err := client.GetAddressUTXOs(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	utxosJson, _ := json.MarshalIndent(utxos, "", "  ")
	fmt.Printf("UTXOs: %s\n", string(utxosJson))
	fmt.Println("--------------------------------")

	fmt.Println("Test address done.")
	fmt.Println("================================")
}
