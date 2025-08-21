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
	btcapis.Init(
		"mainnet",
		"http://192.168.1.16:8332", // bitcoind url
		"cc",                       // bitcoind user
		"ccc",                      // bitcoind pass
		10,                         // 超时时间(秒)
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

	// addr := "1xxxxxxxxxxxxxxxxxxxxxxxxxy1kmdGr" // P2PKH 示例地址
	addr := "3DHgSaYsxCj62UKU2yFG3vKbwjua3ViHUS" // P2SH 示例地址
	// addr := "bc1qgnmdx4pyaxrkhtgeqgh0g93cvar7achq8kjtnm" // P2WPKH 示例地址 - 0个UTXO
	// addr := "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq" // P2TR 示例地址

	scriptInfo, err := btcapis.GetAddressScriptInfo(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressScriptInfo: %v", err)
	}
	fmt.Printf("scriptInfo: %+v\n", scriptInfo)
	fmt.Println("--------------------------------")

	pkScript := scriptInfo.ScriptPubKeyHex
	addrInfo, err := btcapis.GetAddressInfo(ctx, pkScript)
	if err != nil {
		log.Fatalf("GetAddressInfo: %v", err)
	}
	fmt.Printf("addrInfo: %+v\n", addrInfo)
	fmt.Println("--------------------------------")

	confirmed, mempool, err := btcapis.GetAddressBalance(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressBalance: %v", err)
	}
	fmt.Printf("Balance(BTC): %.8f(%.8f)\n", btcapis.SatsToBTC(confirmed), btcapis.SatsToBTC(mempool))
	fmt.Println("--------------------------------")

	utxos, err := btcapis.GetAddressUTXOs(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressUTXOs: %v", err)
	}
	outUtxos, _ := json.MarshalIndent(utxos, "", "  ")
	fmt.Println(string(outUtxos))
	fmt.Println("--------------------------------")

	// utxos, err := btcapis.GetAddressUTXOs(ctx, addr)
	// if err != nil {
	// 	log.Fatalf("GetAddressUTXOs: %v", err)
	// }
	// outUtxos, _ := json.MarshalIndent(utxos, "", "  ")
	// fmt.Println(string(outUtxos))
	// fmt.Println("--------------------------------")

}
