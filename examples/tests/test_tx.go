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

	txid := getTxid()

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
