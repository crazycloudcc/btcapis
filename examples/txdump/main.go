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
	fmt.Printf("raw: %x\n", raw)

	tx, err := client.GetTransaction(ctx, cfg.TxID)
	if err != nil {
		log.Fatalf("GetTransaction: %v", err)
	}

	// 解析输入0
	info, err := btcapis.Tx.AnalyzeInput(tx, 0)
	if err != nil {
		log.Fatalf("AnalyzeInput: %v", err)
	}

	out, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(out))

	// ops, asm, _ := btcapis.Tx.DisasmScriptPubKey(tx, 0)
	// fmt.Printf("ops: %+v\n", ops)
	// fmt.Printf("asm: %s\n", asm)
}
