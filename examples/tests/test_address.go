package tests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/crazycloudcc/btcapis"
	"github.com/crazycloudcc/btcapis/internal/utils"
	"github.com/crazycloudcc/btcapis/types"
)

func TestAddress(client *btcapis.Client) {
	fmt.Println("Testing address...")

	addr := "" // 测试地址

	switch types.CurrentNetwork {
	case types.Mainnet:
		addr = "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq"
		// addr := "1xxxxxxxxxxxxxxxxxxxxxxxxxy1kmdGr" // P2PKH 示例地址
		// addr := "3DHgSaYsxCj62UKU2yFG3vKbwjua3ViHUS" // P2SH 示例地址
		// addr := "bc1qgnmdx4pyaxrkhtgeqgh0g93cvar7achq8kjtnm" // P2WPKH 示例地址 - 0个UTXO
		// addr := "bc1ps2wwxjhw5t33r5tp46yh9x5pukkalsd2vtye07p353fgt7hln5tq763upq" // P2TR 示例地址
	case types.Signet:
		addr = "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw"
	case types.Testnet:
		addr = "tb1pu32s67eye07d05llxr8klr4lj3em3fd6glse5nujmym835x7aw3shp2ffw"
	case types.Regtest:
		addr = ""
	}

	// 查询地址余额
	confirmed, mempool, err := client.GetAddressBalance(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Balance(BTC): %.8f(%.8f)\n", utils.SatsToBTC(confirmed), utils.SatsToBTC(mempool))
	fmt.Println("--------------------------------")

	// 查询地址UTXO
	utxos, err := client.GetAddressUTXOs(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	utxosJson, _ := json.MarshalIndent(utxos, "", "  ")
	fmt.Printf("UTXOs: %s\n", string(utxosJson))
	fmt.Println("--------------------------------")

	// 查询地址脚本信息
	scriptInfo, err := client.GetAddressScriptInfo(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	scriptInfoJson, _ := json.MarshalIndent(scriptInfo, "", "  ")
	fmt.Printf("scriptInfo: %s\n", string(scriptInfoJson))
	fmt.Println("--------------------------------")

	// 通过地址转换为锁定脚本
	pkScript, err := client.AddressToPkScript(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pkScript: %s\n", hex.EncodeToString(pkScript))
	fmt.Println("--------------------------------")

	// 通过地址转类型
	typ, err := client.AddressToType(context.Background(), addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("typ: %s\n", typ)
	fmt.Println("--------------------------------")

	// 通过脚本转类型
	typ2 := client.PKScriptToType(context.Background(), pkScript)
	fmt.Printf("typ2: %s\n", typ2)
	fmt.Println("--------------------------------")

	fmt.Println("Test address done.")
	fmt.Println("================================")
}
