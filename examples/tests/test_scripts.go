package tests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/crazycloudcc/btcapis"
)

func TestScripts(client *btcapis.Client) {
	fmt.Println("Testing scripts...")

	addr := getAddr()
	txid := getTxid()

	// 通过地址获取地址信息
	scriptInfo, err := client.DecodeAddressToScriptInfo(addr)
	if err != nil {
		panic(err)
	}
	scriptInfoJson, _ := json.MarshalIndent(scriptInfo, "", "  ")
	fmt.Printf("scriptInfo: %s\n", string(scriptInfoJson))
	fmt.Println("--------------------------------")

	// 通过地址获取锁定脚本
	pkScript, err := client.DecodeAddressToPkScript(addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pkScript: %s\n", hex.EncodeToString(pkScript))
	fmt.Println("--------------------------------")

	// 通过地址获取类型
	addrType, err := client.DecodeAddressToType(addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("addrType: %s\n", addrType)
	fmt.Println("--------------------------------")

	// 通过脚本获取地址信息
	addrInfo, err := client.DecodePkScriptToAddressInfo(pkScript)
	if err != nil {
		panic(err)
	}
	addrInfoJson, _ := json.MarshalIndent(addrInfo, "", "  ")
	fmt.Printf("addrInfo: %s\n", string(addrInfoJson))
	fmt.Println("--------------------------------")

	// 通过脚本获取类型
	scriptType, err := client.DecodePKScriptToType(pkScript)
	if err != nil {
		panic(err)
	}
	fmt.Printf("scriptType: %s\n", scriptType)
	fmt.Println("--------------------------------")

	// 解析脚本为操作码
	ops, asm, err := client.DecodePkScriptToAsmString(pkScript)
	if err != nil {
		panic(err)
	}
	opsJson, _ := json.MarshalIndent(ops, "", "  ")
	fmt.Printf("ops: %s\n", string(opsJson))
	fmt.Printf("asm: %s\n", asm)
	fmt.Println("--------------------------------")

	// 解析交易hex数据
	txHex, err := client.GetTxRaw(context.Background(), txid)
	if err != nil {
		panic(err)
	}
	decodedTx, err := client.DecodeRawTx(txHex)
	if err != nil {
		panic(err)
	}
	decodedTxJson, _ := json.MarshalIndent(decodedTx, "", "  ")
	fmt.Printf("decodedTx: %s\n", string(decodedTxJson))
	fmt.Println("--------------------------------")

	fmt.Println("Test scripts done.")
	fmt.Println("================================")
}
