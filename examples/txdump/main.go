package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
)

func main() {
	client := btcapis.BuildClient(
		"mainnet",
		"",                      // bitcoind url
		"",                      // bitcoind user
		"",                      // bitcoind pass
		"https://mempool.space", // https://mempool.space/signet
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// txid := "8e2955284d9f66f56df9c89fedd50103fe68f84ec9f138e4e20c67db15de68ee"
	// txid := "2aafdfd00e5a27ef211d7e6ed76a14dd7edd6b88a19df2e535f0963c05372a8d"
	// txid := "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b" // btc链创世交易

	// raw, err := client.GetRawTransaction(ctx, txid)
	// if err != nil {
	// 	log.Fatalf("GetRawTransaction: %v", err)
	// }
	// fmt.Printf("元数据 raw: %x\n", raw)
	// fmt.Println("--------------------------------")

	// tx, err := client.GetTransaction(ctx, txid)
	// if err != nil {
	// 	log.Fatalf("GetTransaction: %v", err)
	// }
	// outTx, _ := json.MarshalIndent(tx, "", "  ")
	// fmt.Println(string(outTx))
	// fmt.Println("--------------------------------")

	// // 解析输入
	// if tx.Vin != nil && len(tx.Vin) > 0 {
	// 	for i, in := range tx.Vin {
	// 		info, err := btcapis.Tx.AnalyzeTxIn(tx, &in)
	// 		if err != nil {
	// 			log.Fatalf("AnalyzeTxIn: %v", err)
	// 		}

	// 		out, _ := json.MarshalIndent(info, "", "  ")
	// 		fmt.Println("-------------------------------- index: ", i)
	// 		fmt.Println(string(out))
	// 		fmt.Println("--------------------------------")
	// 	}
	// }

	// ops, asm, _ := btcapis.Tx.DisasmScriptPubKey(tx, 0)
	// fmt.Printf("ops: %+v\n", ops)
	// fmt.Printf("asm: %s\n", asm)

	// addr := "1xxxxxxxxxxxxxxxxxxxxxxxxxy1kmdGr" // P2PKH 示例地址
	// addr := "3DHgSaYsxCj62UKU2yFG3vKbwjua3ViHUS" // P2SH 示例地址
	// addr := "bc1qgnmdx4pyaxrkhtgeqgh0g93cvar7achq8kjtnm" // P2WPKH 示例地址
	addr := "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq" // P2TR 示例地址

	scriptInfo, err := btcapis.DecodeAddress(addr)
	if err != nil {
		log.Fatalf("DecodeAddress: %v", err)
	}
	fmt.Printf("scriptInfo: %+v\n", scriptInfo)
	fmt.Println("--------------------------------")

	pkScript := scriptInfo.ScriptPubKeyHex
	addrInfo, err := btcapis.DecodePkScript(pkScript)
	if err != nil {
		log.Fatalf("DecodePkScript: %v", err)
	}
	fmt.Printf("addrInfo: %+v\n", addrInfo)
	fmt.Println("--------------------------------")

	confirmed, mempool, err := client.GetAddressBalance(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressBalance: %v", err)
	}
	fmt.Printf("Balance(BTC): %.8f(%.8f)\n", confirmed, mempool)
	fmt.Println("--------------------------------")

	// utxos, err := client.GetAddressUTXOs(ctx, addr)
	// if err != nil {
	// 	log.Fatalf("GetAddressUTXOs: %v", err)
	// }
	// outUtxos, _ := json.MarshalIndent(utxos, "", "  ")
	// fmt.Println(string(outUtxos))
	// fmt.Println("--------------------------------")

}
