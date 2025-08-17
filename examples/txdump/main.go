package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

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
	out, _ := json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(out))
}
