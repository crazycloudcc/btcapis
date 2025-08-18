// Package main 基础使用示例
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/crazycloudcc/btcapis"
)

func main() {
	c := btcapis.BuildClient(
		os.Getenv("BITCOIND_URL"),
		os.Getenv("BITCOIND_USER"),
		os.Getenv("BITCOIND_PASS"),
		os.Getenv("MEMPOOL_BASE_URL"),
	)
	tx, err := c.GetTransaction(context.Background(), "your-txid-here")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.TxID, len(tx.Vin), len(tx.Vout), tx.Vsize)
}
