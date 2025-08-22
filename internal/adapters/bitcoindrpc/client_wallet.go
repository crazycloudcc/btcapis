// 钱包相关操作
// 封装只能对本地钱包操作的rpc接口
package bitcoindrpc

// 4. 钱包（生成/加载/余额/地址/UTXO）
// createwallet <name>、
// loadwallet <name>、
// listwallets、
// getwalletinfo
// getnewaddress [label] [address_type]（legacy/p2sh-segwit/bech32/bech32m；Taproot 用 bech32m）
// validateaddress <address>（补充校验）
// getaddressesbylabel <label>（可罗列同标签地址）
// listunspent [minconf] [maxconf] [addresses] [include_unsafe] [query_options]
// lockunspent false/true [...]（锁定/解锁 UTXO，防止并发争抢）
// getbalance、
// getbalances
// gettransaction <txid>（钱包视角，含确认数、收支变动）

// 6. 交易构建（原生/PSBT）
// 快速转账：sendtoaddress <address> <amount> [comment] [comment_to] [subtractfeefromamount]

// 1. 链上数据查询 / 分析
// gettxout <txid> <vout> <include_mempool>get
// 可选索引：getblockfilter <hash>（需 -blockfilterindex=1）

// 5. 描述符 / 扫描（非钱包地址的 UTXO 查询）
// importdescriptors、
// deriveaddresses <descriptor> [range]
// getdescriptorinfo <descriptor>
// 6. 交易构建（原生/PSBT）
//
// 原始交易流：createrawtransaction → fundrawtransaction → signrawtransactionwithwallet|withkey → sendrawtransaction
// PSBT流：walletcreatefundedpsbt →（外部签名或 walletprocesspsbt）→ finalizepsbt → sendrawtransaction
// utxoupdatepsbt、
// decodepsbt（调试/补全）
