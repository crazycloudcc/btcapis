package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
)

// const (
//
//	network = "mainnet"
//	rpcUser = "cc"
//	rpcPass = "ccc"
//	rpcUrl  = "http://192.168.1.16:8332"
//	timeout = 10 * time.Second
//
// )
const (
	network = "signet"
	rpcUser = "cc"
	rpcPass = "ccc"
	rpcUrl  = "http://localhost:18332"
	timeout = 10 * time.Second
)

func main() {
	client, testClient := btcapis.New(
		network,
		rpcUrl,
		rpcUser,
		rpcPass,
		int(timeout.Seconds()),
	)

	if client == nil {
		log.Fatalf("New: %v", errors.New("client is nil"))
	}

	if testClient == nil {
		log.Fatalf("New: %v", errors.New("testClient is nil"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// txid := "0d31e59675c85f17d942f4510bb4760d9ed4b661df22af3b7cd5ef3c2116626b" // 测试交易
	// txid := "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b" // btc链创世交易
	// txid := "8e2955284d9f66f56df9c89fedd50103fe68f84ec9f138e4e20c67db15de68ee"
	// txid := "2aafdfd00e5a27ef211d7e6ed76a14dd7edd6b88a19df2e535f0963c05372a8d"

	addr := "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw" // 测试地址
	// addr := "1xxxxxxxxxxxxxxxxxxxxxxxxxy1kmdGr" // P2PKH 示例地址
	// addr := "3DHgSaYsxCj62UKU2yFG3vKbwjua3ViHUS" // P2SH 示例地址
	// addr := "bc1qgnmdx4pyaxrkhtgeqgh0g93cvar7achq8kjtnm" // P2WPKH 示例地址 - 0个UTXO
	// addr := "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq" // P2TR 示例地址

	scriptInfo, err := client.GetAddressScriptInfo(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressScriptInfo: %v", err)
	}
	fmt.Printf("scriptInfo: %+v\n", scriptInfo)
	fmt.Println("--------------------------------")

	pkScript := scriptInfo.ScriptPubKeyHex
	addrInfo, err := client.GetAddressInfo(ctx, pkScript)
	if err != nil {
		log.Fatalf("GetAddressInfo: %v", err)
	}
	fmt.Printf("addrInfo: %+v\n", addrInfo)
	fmt.Println("--------------------------------")

	confirmed, mempool, err := client.GetAddressBalance(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressBalance: %v", err)
	}
	fmt.Printf("Balance(BTC): %.8f(%.8f)\n", btcapis.SatsToBTC(confirmed), btcapis.SatsToBTC(mempool))
	fmt.Println("--------------------------------")

	utxos, err := client.GetAddressUTXOs(ctx, addr)
	if err != nil {
		log.Fatalf("GetAddressUTXOs: %v", err)
	}
	outUtxos, _ := json.MarshalIndent(utxos, "", "  ")
	fmt.Println(string(outUtxos))
	fmt.Println("--------------------------------")

	// testrpc(client, testClient)
}

// 测试rpc是否正常
func testrpc(client *btcapis.Client, testClient *btcapis.TestClient) {
	res, err := testClient.GetNetworkInfo(context.Background())
	if err != nil {
		log.Fatalf("GetNetworkInfo: %v", err)
	}
	fmt.Printf("res: %+v\n", res)
	fmt.Println("--------------------------------")

	feerate1, feerate2, err := client.EstimateFeeRate(context.Background(), 1)
	if err != nil {
		log.Fatalf("EstimateFeeRate: %v", err)
	}
	fmt.Printf("feerate1: %.2f (sats/vB)\n", feerate1*1e8/1000.0)
	fmt.Printf("feerate2: %.2f (sats/vB)\n", feerate2)
	fmt.Println("--------------------------------")

	// rawtx, err := client.BuildTx(context.Background())
	// if err != nil {
	// 	log.Fatalf("BuildTx: %v", err)
	// }
	// fmt.Printf("rawtx: %s\n", hex.EncodeToString(rawtx))
	// fmt.Println("--------------------------------")

	// result, err := client.FundTx(context.Background(), hex.EncodeToString(rawtx))
	// if err != nil {
	// 	log.Fatalf("FundTx: %v", err)
	// }
	// fmt.Printf("result: %+v\n", result)
	// fmt.Println("--------------------------------")
}
