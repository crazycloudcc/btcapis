package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
)

func main() {
	client := btcapis.BuildClient(
		"",                      // bitcoind url
		"",                      // bitcoind user
		"",                      // bitcoind pass
		"https://mempool.space", // https://mempool.space/signet
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	txid := "8e2955284d9f66f56df9c89fedd50103fe68f84ec9f138e4e20c67db15de68ee"

	raw, err := client.GetRawTransaction(ctx, txid)
	if err != nil {
		log.Fatalf("GetRawTransaction: %v", err)
	}
	fmt.Printf("元数据 raw: %x\n", raw)
	fmt.Println("--------------------------------")

	tx, err := client.GetTransaction(ctx, txid)
	if err != nil {
		log.Fatalf("GetTransaction: %v", err)
	}
	outTx, _ := json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(outTx))
	fmt.Println("--------------------------------")

	// 解析输入
	if tx.Vin != nil && len(tx.Vin) > 0 {
		for i, in := range tx.Vin {
			info, err := btcapis.Tx.AnalyzeTxIn(tx, &in)
			if err != nil {
				log.Fatalf("AnalyzeTxIn: %v", err)
			}

			out, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println("-------------------------------- index: ", i)
			fmt.Println(string(out))
			fmt.Println("--------------------------------")
		}
	}

	// addr := "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq"
	// balance, err := client.GetAddressBalance(ctx, addr)
	// if err != nil {
	// 	log.Fatalf("GetAddressBalance: %v", err)
	// }
	// fmt.Printf("balance: %s\n", balance)
	// fmt.Println("--------------------------------")

	// utxos, err := client.GetAddressUTXOs(ctx, addr)
	// if err != nil {
	// 	log.Fatalf("GetAddressUTXOs: %v", err)
	// }
	// outUtxos, _ := json.MarshalIndent(utxos, "", "  ")
	// fmt.Println(string(outUtxos))
	// fmt.Println("--------------------------------")

	// ops, asm, _ := btcapis.Tx.DisasmScriptPubKey(tx, 0)
	// fmt.Printf("ops: %+v\n", ops)
	// fmt.Printf("asm: %s\n", asm)
}
