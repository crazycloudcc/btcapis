package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/crazycloudcc/btcapis"
)

const (
	network = "mainnet"
	rpcUser = "cc"
	rpcPass = "ccc"
	rpcUrl  = "http://192.168.1.16:8332"
	timeout = 30 * time.Second
)

func main() {
	btcapis.Init(
		network,
		rpcUrl,
		rpcUser,
		rpcPass,
		int(timeout.Seconds()),
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

	// utxos, err := btcapis.GetAddressUTXOs(ctx, addr)
	// if err != nil {
	// 	log.Fatalf("GetAddressUTXOs: %v", err)
	// }
	// outUtxos, _ := json.MarshalIndent(utxos, "", "  ")
	// fmt.Println(string(outUtxos))
	// fmt.Println("--------------------------------")

	testrpc()
}

// 测试rpc是否正常
func testrpc() {
	node := rpcUrl
	user := rpcUser
	pass := rpcPass

	body := []byte(`{"jsonrpc":"1.0","id":"go","method":"getblockcount","params":[]}`)
	req, _ := http.NewRequest("POST", node, bytes.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(out))
}
