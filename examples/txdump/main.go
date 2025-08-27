package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/crazycloudcc/btcapis"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/internal/types"
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
	client := btcapis.New(
		network,
		rpcUrl,
		rpcUser,
		rpcPass,
		int(timeout.Seconds()),
	)

	if client == nil {
		log.Fatalf("New: %v", errors.New("client is nil"))
	}

	testClient := btcapis.NewTestClient(client)
	if testClient == nil {
		log.Fatalf("NewTestClient: %v", errors.New("testClient is nil"))
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
	// fmt.Printf("scriptInfo: %+v\n", scriptInfo)
	fmt.Println("--------------------------------")

	pkScript := scriptInfo.ScriptPubKeyHex
	addrInfo, err := client.GetAddressInfo(ctx, pkScript)
	if err != nil {
		log.Fatalf("GetAddressInfo: %v", err)
	}
	fmt.Printf("addrInfo: %+v\n", addrInfo)
	fmt.Println("--------------------------------")

	typ := decoders.PKScriptToType(pkScript)
	fmt.Printf("typ: %s\n", typ)
	fmt.Println("--------------------------------")

	confirmed, mempool, err := client.GetAddressBalanceBTC(ctx, addr)
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

	testrpc(client, testClient)
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

	// // 构建、填充、签名交易 测试
	txInputParams := &types.TxInputParams{
		FromAddress: []string{ // 必须是钱包内地址
			"tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw",
		},
		ToAddress: []string{
			"tb1pn3rx2vtzfrazpqzfrftjzye5k3szvjskqlwvfu9pv9dtv6a8wv3sper2ce",
		},
		AmountBTC: []float64{
			0.0001,
		},
		FeeRate:       1.0, // sat/vB
		Locktime:      0,   // 0=默认最新区块高度
		Replaceable:   true,
		Data:          "hello world", // 可选
		PublicKey:     "9a235c04856d94389042a8e12a500fe9a80dbdb090ec9b235762a706a475b20a",
		ChangeAddress: "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw",
	}

	psbt, err := client.CreatePSBT(context.Background(), txInputParams)
	if err != nil {
		log.Fatalf("CreatePSBT: %v", err)
	}
	fmt.Printf("psbt base64: %s\n", psbt.PSBTBase64)
	fmt.Printf("unsigned tx hex: %s\n", psbt.UnsignedTxHex)
	fmt.Printf("estimated vsize: %d vB\n", psbt.EstimatedVSize)
	fmt.Printf("fee sat: %d sats\n", psbt.FeeSat)
	fmt.Printf("change output index: %d\n", psbt.ChangeOutputIdx)
	fmt.Printf("psbt struct: %+v\n", psbt.Packet)
	fmt.Println("--------------------------------")

	// 	// 输出给调用方/前端（包含 OKX 可用的 PSBT base64）
	// type BuildResult struct {
	// 	PSBTBase64      string  `json:"psbt_base64"`        // 给 OKX
	// 	UnsignedTxHex   string  `json:"unsigned_tx_hex"`    // 调试/核对
	// 	Packet          *Packet `json:"-"`                  // 你自定义 psbt 结构，便于回写签名/Finalize
	// 	EstimatedVSize  int     `json:"estimated_vsize_vb"` // 估算
	// 	FeeSat          int64   `json:"fee_sat"`
	// 	ChangeOutputIdx int     `json:"change_output_index"` // -1 表示没有找零
	// }

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
