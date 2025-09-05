package tests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/crazycloudcc/btcapis"
	"github.com/crazycloudcc/btcapis/types"
)

// 测试交易相关接口
func TestTxs(client *btcapis.Client) {
	fmt.Println("Testing tx...")

	txid := ""
	switch types.CurrentNetwork {
	case types.Mainnet:
		txid = "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b" // btc链创世交易
	// txid := "8e2955284d9f66f56df9c89fedd50103fe68f84ec9f138e4e20c67db15de68ee" // main 839639
	// txid := "2aafdfd00e5a27ef211d7e6ed76a14dd7edd6b88a19df2e535f0963c05372a8d" // main 910115
	case types.Signet:
		txid = "12ad394188be183dccc53357a4374a7eab067810c3535063012986dc437e8a3c"
		// txid := "0d31e59675c85f17d942f4510bb4760d9ed4b661df22af3b7cd5ef3c2116626b" // 测试交易
	case types.Testnet:
		txid = "" // 空交易
	case types.Regtest:
		txid = ""
	}

	// 交易详情
	txDetail, err := client.GetTx(context.Background(), txid)
	if err != nil {
		panic(err)
	}
	txDetailJson, _ := json.MarshalIndent(txDetail, "", "  ")
	fmt.Printf("txDetail: %s\n", string(txDetailJson))
	fmt.Println("--------------------------------")

	// 交易hex数据
	txHex, err := client.GetTxRaw(context.Background(), txid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("txHex: %s\n", hex.EncodeToString(txHex))
	fmt.Println("--------------------------------")

	// 解析交易hex数据
	decodedTx, err := client.DecodeRawTx(context.Background(), txHex)
	if err != nil {
		panic(err)
	}
	decodedTxJson, _ := json.MarshalIndent(decodedTx, "", "  ")
	fmt.Printf("decodedTx: %s\n", string(decodedTxJson))
	fmt.Println("--------------------------------")

	// 构建、填充、签名交易 测试
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

	psbtBase64, err := client.CreatePSBT(context.Background(), txInputParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("psbt base64: %s\n", psbtBase64)
	fmt.Println("--------------------------------")

	fmt.Println("Test tx done.")
	fmt.Println("================================")
}
