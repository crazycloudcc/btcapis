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
	c := btcapis.New(
		btcapis.WithBitcoindRPC(os.Getenv("BITCOIND_URL"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS")),
		btcapis.WithMempoolSpace("https://mempool.space"),
	)
	tx, err := c.GetTransaction(context.Background(), "your-txid-here")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.TxID, len(tx.Vin), len(tx.Vout), tx.Vsize)
}
