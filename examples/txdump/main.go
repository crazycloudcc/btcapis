package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/crazycloudcc/btcapis"
	"github.com/crazycloudcc/btcapis/examples/common"
)

func main() {
	envPath := flag.String("env", "examples/env.json", "path to env.json")
	flag.Parse()

	cfg, err := common.LoadConfig(*envPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if cfg.TxID == "" {
		log.Fatal("txid not set in env.json")
	}
	client := cfg.BuildClient()

	ctx, cancel := cfg.Ctx()
	defer cancel()

	raw, err := client.GetRawTransaction(ctx, cfg.TxID)
	if err != nil {
		log.Fatalf("GetRawTransaction: %v", err)
	}
	fmt.Printf("元数据 raw: %x\n", raw)
	fmt.Println("--------------------------------")

	tx, err := client.GetTransaction(ctx, cfg.TxID)
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

	// balance, err := client.GetAddressBalance(ctx, "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq")
	// if err != nil {
	// 	log.Fatalf("GetAddressBalance: %v", err)
	// }
	// fmt.Printf("balance: %s\n", balance)
	// fmt.Println("--------------------------------")

	// utxos, err := client.GetAddressUTXOs(ctx, "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq")
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
