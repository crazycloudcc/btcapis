# BTCAPIs

**English** | [简体中文](README.zh-CN.md)

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go library for Bitcoin on-chain operations, used by services such as [chainbox](https://github.com/chainboxapp). It covers address queries, transaction build/broadcast, PSBT, chain RPC, mempool queries, fee estimation, and script parsing — with multi-source fallback across Bitcoin Core RPC, mempool.space, and ElectrumX.

Current version: **v0.4.0** (see [CHANGELOG.md](CHANGELOG.md))

## Features

- **Address**: balance/UTXO queries, validation, wallet generation; ElectrumX batch queries and custom Provider injection
- **Transaction**: raw/verbose tx queries, mempool pre-check, PSBT v0 create/validate/broadcast
- **Chain RPC**: network info, chain state, block header/block, hash by height
- **Fees**: unified `types.ChainFees` (sat/vB); default fallback chain + API chain (unisat/okx)
- **Mempool**: stats, tx status, recommended fees, mempool txid/entry
- **Script/Decode**: address↔script conversion, ASM parsing, raw tx decoding
- **Networks**: mainnet / testnet / signet

## Install

```bash
go get github.com/crazycloudcc/btcapis@v0.4.0
```

## Quick Start

### Create a client

```go
package main

import (
    "context"
    "log"

    "github.com/crazycloudcc/btcapis"
)

func main() {
    client := btcapis.New(&btcapis.Config{
        Network:         "testnet",
        Timeout:         30,
        RPCUrl:          "http://127.0.0.1:18332",
        RPCUser:         "rpcuser",      // inject via env or config — never hardcode
        RPCPass:         "rpcpassword",
        MempoolSpaceUrl: "https://mempool.space/testnet",
        ElectrumXUrl:    "", // optional, e.g. https://blockstream.info/testnet/api
    })
    if client == nil {
        log.Fatal("client init failed")
    }

    ctx := context.Background()
    info, err := client.GetBlockChainInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("chain=%s blocks=%d", info.Chain, info.Blocks)
}
```

Shortcut constructor (auto-fills mempool.space URL by network):

```go
client := btcapis.NewWithElectrumX(
    "testnet",
    "http://127.0.0.1:18332", "rpcuser", "rpcpassword",
    "wss://electrum.example.com:50002", // ElectrumX, may be empty
    30,
)
```

> **Security**: RPC passwords, Unisat/OKX API keys, and other credentials must be injected via environment variables or local config files. **Never commit them to git.** See `.gitignore`.

### Address queries

```go
// default: mempool.space → bitcoin core
confirmed, mempool, err := client.GetAddressBalance(ctx, "tb1q...")
utxos, err := client.GetAddressUTXOs(ctx, "tb1q...")

// with custom ElectrumX Provider (e.g. TCP — implement types.ElectrumXAddressProvider)
balance, err := client.GetAddressBalanceSatsWithOptions(ctx, "tb1q...", types.AddressProviderOptions{
    ElectrumX: myElectrumProvider,
})

addrType, err := client.DecodeAddressToType("tb1q...")
```

### Transactions & PSBT

```go
import "github.com/crazycloudcc/btcapis/types"

tx, err := client.GetTx(ctx, "txid...")
verbose, err := client.GetTxVerbose(ctx, "txid...", 1)

psbtBase64, err := client.CreatePSBT(ctx, &types.TxInputParams{
    FromAddress:   []string{"tb1p..."},
    ToAddress:     []string{"tb1q..."},
    AmountBTC:     []float64{0.001},
    FeeRate:       1.0,
    ChangeAddress: "tb1p...",
    Replaceable:   true,
})

txid, err := client.FinalizePSBTAndBroadcast(ctx, signedPSBT)
```

### Fee estimation

```go
// default fallback: mempool.space → electrumX → bitcoin core
fees, err := client.EstimateChainFees(ctx, 6)
// fees.High / Medium / Low / FeeRate (sat/vB, 2 decimal places)

// full API fallback chain (create extension clients with API keys)
import (
    "github.com/crazycloudcc/btcapis/extensions/okx"
    "github.com/crazycloudcc/btcapis/extensions/unisat"
)

unisatClient, _ := unisat.New(unisat.Config{APIKey: os.Getenv("UNISAT_API_KEY")})
okxClient, _ := okx.New(okx.Config{
    APIKey:     os.Getenv("OKX_API_KEY"),
    SecretKey:  os.Getenv("OKX_SECRET_KEY"),
    Passphrase: os.Getenv("OKX_PASSPHRASE"),
})
fees, err = client.EstimateChainFeesForAPI(ctx, 6, btcapis.APIFeesProviderOptions{
    Unisat: unisatClient,
    OKX:    okxClient,
})

// legacy API
high, low, err := client.EstimateFeeRate(ctx, 6)
```

## Public API

### Address — `btcapis_address.go` / `btcapis_mempool.go`

| Method | Description |
|--------|-------------|
| `CreateNewWallet` | Generate new wallet (mnemonic/address/private key) |
| `GetAddressBalance` | Confirmed/unconfirmed balance (BTC) |
| `GetAddressUTXOs` | Address UTXO list |
| `GetAddressBalanceSats` | Balance in sats via default provider |
| `GetAddressUTXOsForAddress` | UTXOs as `types.AddressUTXO` |
| `GetAddressBalanceSatsWithOptions` | Balance: electrumX → mempool.space |
| `GetAddressUTXOsWithOptions` | UTXOs: electrumX → mempool.space |
| `GetAddressHistoryTxs` | Tx history (requires `ElectrumXAddressProvider`) |
| `GetAddressUnconfirmedTxs` | Unconfirmed txs (requires ElectrumX) |
| `BatchGetAddressBalanceSats` | Batch balances (requires ElectrumX) |
| `GetAddressBalanceWithElectrumX` | ElectrumX balance |
| `GetAddressBalanceWithElectrumXByXPRV` | Balance via xprv-derived addresses |
| `GetAddressBalanceWithElectrumXByPrivateKey` | Balance via WIF |
| `FilterAddressesWithBalanceWithElectrumX` | Filter addresses with balance |
| `BatchGetBalancesWithElectrumX` | ElectrumX batch balances |
| `ValidateAddress` | Bitcoin Core `validateaddress` |

### Transaction — `btcapis_tx.go`

| Method | Description |
|--------|-------------|
| `GetTx` / `GetTxRaw` | Query and decode transaction |
| `GetTxVerbose` | Verbose tx (core → mempool fallback) |
| `CreatePSBT` | Create PSBT v0 |
| `FinalizePSBTAndBroadcast` | Finalize PSBT and broadcast |
| `BroadcastRawTx` | Broadcast raw transaction |
| `ValidateUnsignedPsbtBase64` | Validate unsigned PSBT |
| `ValidateSignedPsbtBase64` | Validate signed PSBT |
| `TransferAllToNewAddress` | Transfer all (WIF signing) |
| `TestMempoolAccept` | Mempool accept pre-check |
| `GetMempoolTxIds` / `GetMempoolTxEntry` | Mempool tx queries |

### Chain — `btcapis_chain.go`

| Method | Description |
|--------|-------------|
| `EstimateChainFees` | Default fee fallback chain |
| `EstimateChainFeesWithOptions` | Custom ElectrumX fee provider |
| `EstimateChainFeesWithProviders` | Fully custom provider order |
| `EstimateChainFeesForAPI` | Full API fallback chain |
| `EstimateFeeRate` | Legacy API (high/low sat/vB) |
| `GetNetworkInfo` / `GetBlockChainInfo` | Node/chain state |
| `GetBlockStats` / `GetChainTips` | Block stats/fork tips |
| `GetBlockHeader` / `GetBlock` | Block header/block (verbosity=1) |
| `GetBlockHashByHeight` | Block hash by height |

### Mempool — `btcapis_mempool.go`

| Method | Description |
|--------|-------------|
| `GetMempoolStats` | Mempool statistics |
| `GetMempoolTxStatus` | Transaction confirmation status |
| `GetMempoolFeesRecommend` | mempool.space recommended fees |

### Script/Decode — `btcapis_scripts.go`

| Method | Description |
|--------|-------------|
| `DecodeAddressToScriptInfo` / `DecodeAddressToPkScript` / `DecodeAddressToType` | Address parsing |
| `DecodePkScriptToAddressInfo` / `DecodePKScriptToType` | Script → address |
| `DecodePkScriptToAsmString` | Script → opcodes/ASM |
| `DecodeRawTx` / `DecodeRawTxString` | Raw tx decoding |

### Format checkers — `btcapis_checkers.go`

`CheckFormatAddress` / `CheckFormatTxid` / `CheckFormatHex` / `CheckFormatBase64`

### Extension packages

| Package | Purpose |
|---------|---------|
| `extensions/unisat` | Unisat Open API fee provider |
| `extensions/okx` | OKX Web3 API fee provider |

### Custom Provider interfaces — `types/`

```go
// Address data source (implement in caller, e.g. TCP ElectrumX)
type ElectrumXAddressProvider interface {
    Name() string
    GetBalanceSats(ctx context.Context, addr string) (*AddressBalanceSats, error)
    GetUTXOs(ctx context.Context, addr string) ([]AddressUTXO, error)
    GetHistoryTxs(ctx context.Context, addr string) ([]AddressHistoryTx, error)
    GetUnconfirmedTxs(ctx context.Context, addr string) ([]AddressUnconfirmedTx, error)
    BatchGetBalanceSats(ctx context.Context, addrs []string) ([]AddressBalanceSats, error)
}

// Fee data source
type ChainFeesEstimator interface { ... }
```

## Data source configuration

Fill `Config` fields as needed; adapters without config are not initialized:

| Field | Description |
|-------|-------------|
| `Network` | `mainnet` / `testnet` / `signet` |
| `RPCUrl` / `RPCUser` / `RPCPass` | Bitcoin Core JSON-RPC |
| `MempoolSpaceUrl` | mempool.space API base URL |
| `ElectrumXUrl` | ElectrumX HTTP endpoint |
| `Timeout` | HTTP/RPC timeout (seconds) |

Default mempool.space URLs (`NewWithElectrumX` auto-selects):

- mainnet → `https://mempool.space`
- testnet → `https://mempool.space/testnet`
- signet → `https://mempool.space/signet`

## Address types

| Type | Prefix | Supported |
|------|--------|-----------|
| P2PKH | `1...` / `m...` / `n...` | ✅ |
| P2SH | `3...` / `2...` | ✅ |
| P2WPKH / P2WSH | `bc1q...` / `tb1q...` | ✅ |
| P2TR (Taproot) | `bc1p...` / `tb1p...` | ✅ |

## Project structure

```
btcapis.go / btcapis_*.go     # public facade API
types/                        # public types and Provider interfaces
extensions/unisat|okx/        # optional fee extensions
internal/
  adapters/
    bitcoindrpc/              # Bitcoin Core RPC
    mempoolapis/              # mempool.space REST
    electrumx/                # ElectrumX HTTP
  address/                    # address queries and orchestration
  chain/                      # chain RPC, fees, mempool
  tx/                         # tx build, PSBT, broadcast
  decoders/                   # address/script/tx decoding
  ordinals/ / runes/          # internal only, not exported yet
pkg/logger/                   # logging
docs/                         # ElectrumX docs and RPC API lists
```

## Development

```bash
git clone https://github.com/crazycloudcc/btcapis.git
cd btcapis
export GOCACHE=$PWD/.gocache
go mod tidy
go test ./...
golangci-lint run   # requires local install
```

Use `TestClient` (`btcapis_tests.go`) for integration tests against real bitcoindrpc/mempool/electrumx backends.

## Documentation

| File | Content |
|------|---------|
| [CHANGELOG.md](CHANGELOG.md) | Version history |
| [docs/electrumx_quickstart.md](docs/electrumx_quickstart.md) | ElectrumX quickstart |
| [docs/electrumx_implementation.md](docs/electrumx_implementation.md) | ElectrumX implementation |
| [docs/2.bitcoin-core-rpc接口清单.txt](docs/2.bitcoin-core-rpc接口清单.txt) | Bitcoin Core RPC list |
| [docs/3.mempool-space接口清单.txt](docs/3.mempool-space接口清单.txt) | mempool.space API list |

## Dependencies

- [btcsuite/btcd](https://github.com/btcsuite/btcd) — tx/script/address primitives
- [shopspring/decimal](https://github.com/shopspring/decimal) — decimal arithmetic
- [sirupsen/logrus](https://github.com/sirupsen/logrus) — logging

Public API uses custom `types/` structs, decoupled from btcd types.

## License

This project is licensed under the [MIT License](LICENSE).
