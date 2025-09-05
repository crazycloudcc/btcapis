package tx

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// // 将交易 MsgTx 转为 PSBTv2 格式;
// func (c *Client) MsgTxToPSBTV2(ctx context.Context, tx *wire.MsgTx, inputParams *types.TxInputParams, utxos []*types.TxUTXO) (*types.TxUnsignedPSBT, error) {
// 	w := bytes.NewBuffer()
// }

// 将交易 MsgTx 转为 PSBTv0 格式;
// 因为btcd还不支持v2, 暂时使用v0;
func (c *Client) MsgTxToPSBTV0(ctx context.Context, tx *wire.MsgTx, inputParams *types.TxInputParams, utxos []*types.TxUTXO) (*types.TxUnsignedPSBT, error) {
	packet, err := psbt.NewFromUnsignedTx(tx)
	if err != nil {
		return nil, fmt.Errorf("创建 PSBT 失败: %v", err)
	}

	upd, err := psbt.NewUpdater(packet)
	if err != nil {
		return nil, fmt.Errorf("创建 PSBT 更新器失败: %v", err)
	}

	addrScriptInfo, err := decoders.DecodeAddress(inputParams.FromAddress[0])
	if err != nil {
		return nil, fmt.Errorf("解析地址失败: %v", err)
	}

	for i := 0; i < len(utxos); i++ {
		switch addrScriptInfo.Typ {
		case types.AddrP2PKH: // P2PKH 需要 NonWitnessTx
			txRaw, err := c.bitcoindrpcClient.TxGetRaw(ctx, utxos[i].OutPoint.Hash.String(), false)
			if err != nil {
				return nil, fmt.Errorf("failed to get raw tx for %s: %w", utxos[i].OutPoint.Hash.String(), err)
			}
			var prevTx wire.MsgTx
			r := bytes.NewReader(txRaw)
			if err := prevTx.Deserialize(r); err != nil {
				return nil, fmt.Errorf("deserialize prev tx failed: %w", err)
			}
			upd.AddInNonWitnessUtxo(&prevTx, i)
		case types.AddrP2SH: // P2SH 需要 RedeemScript
			if len(addrScriptInfo.RedeemScriptHashHex) > 0 && isSegwitProgram(addrScriptInfo.RedeemScriptHashHex) {
				txout := &wire.TxOut{Value: utxos[i].Value, PkScript: addrScriptInfo.ScriptPubKeyHex}
				upd.AddInWitnessUtxo(txout, i)

				upd.AddInRedeemScript(addrScriptInfo.RedeemScriptHashHex, i)

				if len(addrScriptInfo.WitnessProgramHex) > 0 {
					upd.AddInWitnessScript(addrScriptInfo.WitnessProgramHex, i)
				}
			} else { // 普通P2SH, 非嵌套Segwit => 和P2PKH逻辑相同
				txRaw, err := c.bitcoindrpcClient.TxGetRaw(ctx, utxos[i].OutPoint.Hash.String(), false)
				if err != nil {
					return nil, fmt.Errorf("failed to get raw tx for %s: %w", utxos[i].OutPoint.Hash.String(), err)
				}
				var prevTx wire.MsgTx
				r := bytes.NewReader(txRaw)
				if err := prevTx.Deserialize(r); err != nil {
					return nil, fmt.Errorf("deserialize prev tx failed: %w", err)
				}
				upd.AddInNonWitnessUtxo(&prevTx, i)
			}
		case types.AddrP2WPKH: // 需要 WitnessScript
			txout := &wire.TxOut{Value: utxos[i].Value, PkScript: addrScriptInfo.ScriptPubKeyHex}
			upd.AddInWitnessUtxo(txout, i)
		case types.AddrP2WSH: // 需要 WitnessScript, 额外需要witness script
			txout := &wire.TxOut{Value: utxos[i].Value, PkScript: addrScriptInfo.ScriptPubKeyHex}
			upd.AddInWitnessUtxo(txout, i)
			// 额外需要witness script
			if len(addrScriptInfo.WitnessProgramHex) > 0 {
				upd.AddInWitnessScript(addrScriptInfo.WitnessProgramHex, i)
			}
		case types.AddrP2TR: // P2TR 需要 WitnessScript
			txout := &wire.TxOut{Value: utxos[i].Value, PkScript: addrScriptInfo.ScriptPubKeyHex}
			upd.AddInWitnessUtxo(txout, i)
		default:
			fmt.Printf("Unsupported address type for PSBT input: %v\n", addrScriptInfo.Typ)
		}
	}

	errCheck := packet.SanityCheck()
	if errCheck != nil {
		return nil, fmt.Errorf("PSBT Sanity Check failed: %v", errCheck)
	}

	psbtBase64, err := packet.B64Encode()
	if err != nil {
		var buf bytes.Buffer
		if err2 := packet.Serialize(&buf); err2 != nil {
			return nil, fmt.Errorf("PSBT 编码失败: %v / %v", err, err2)
		}
		psbtBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	}

	var raw bytes.Buffer
	if err := tx.Serialize(&raw); err != nil {
		return nil, fmt.Errorf("序列化交易失败: %v", err)
	}

	return &types.TxUnsignedPSBT{
		PSBTBase64: psbtBase64,
		UnsignedTx: hex.EncodeToString(raw.Bytes()),
	}, nil
}

// 接收OKX签名后的交易数据并解析
func (c *Client) FinalizePSBT(ctx context.Context, signedPSBT string) ([]byte, error) {
	// 兼容 OKX psbtHex 与 base64 两种输入
	normalized := strings.TrimSpace(signedPSBT)
	var psbtBase64 string

	// 判定十六进制
	isHex := func(s string) bool {
		if len(s)%2 != 0 || len(s) == 0 {
			return false
		}
		for i := 0; i < len(s); i++ {
			ch := s[i]
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				return false
			}
		}
		return true
	}
	fmt.Printf("psbt input len: %d, isHex: %v\n", len(normalized), isHex(normalized))

	// 十六进制转base64
	if isHex(normalized) {
		bin, err := hex.DecodeString(normalized)
		if err != nil {
			return nil, err
		}
		psbtBase64 = base64.StdEncoding.EncodeToString(bin)
	} else {
		psbtBase64 = normalized
	}
	fmt.Printf("psbt base64 len: %d\n", len(psbtBase64))
	fmt.Printf("psbt base64: %s\n", psbtBase64)

	// finalizepsbt -> 原始交易hex
	finalizeData, err := c.bitcoindrpcClient.TxFinalizePsbt(ctx, psbtBase64)
	if err != nil || !finalizeData.Complete {
		fmt.Printf("finalizepsbt error: %v, complete: %v\n", err, finalizeData.Complete)
		return nil, err
	}
	fmt.Println("finalize rawHex: ", finalizeData.Hex)

	rawTx, err := hex.DecodeString(finalizeData.Hex)
	if err != nil {
		return nil, err
	}

	fmt.Println("broadcast tx len: ", len(rawTx))
	return rawTx, nil
}

// helper: 判定“SegWit 程序”本体（用于 P2SH redeemScript：P2WPKH/P2WSH/Taproot）
func isSegwitProgram(script []byte) bool {
	n := len(script)
	if n < 4 || n > 42 {
		return false
	}
	ver := script[0]
	push := int(script[1])
	if ver != 0x00 && (ver < 0x51 || ver > 0x60) {
		return false
	}
	if 2+push != n {
		return false
	}
	// 仅接受 20/32（v0-pkh / v0-wsh / v1-taproot）的典型长度
	return push == 20 || push == 32
}
