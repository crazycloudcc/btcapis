package tx

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

// 将交易 MsgTx 转为 PSBTv0 格式;
// 因为btcd还不支持v2, 暂时使用v0;
func (c *Client) MsgTxToPsbtV0(ctx context.Context, tx *wire.MsgTx, inputParams *types.TxInputParams, utxos []*types.TxUTXO) (*types.TxUnsignedPSBT, error) {
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

// // 3. 将TxUTXO转为PsbtUTXO结构
// 	psbtUTXOs := make([]psbt.PsbtUTXO, 0, len(selectedUTXOs))
// for _, utxo := range selectedUTXOs {

// 	pkScript := utxo.PkScript
// 	if len(pkScript) == 0 {
// 		// 如果没有 pkScript，尝试通过地址解析
// 		pkScript = addrScriptInfo.ScriptPubKeyHex
// 	}

// 	nonWitnessTxHex := ""
// 	if decoders.PKScriptToType(utxo.PkScript) == types.AddrP2PKH {
// 		txRaw, err := c.bitcoindrpcClient.TxGetRaw(ctx, utxo.OutPoint.Hash.String(), false)
// 		if err != nil {
// 			fmt.Print("Error fetching raw tx \n")
// 			return nil, fmt.Errorf("failed to get raw tx for %s: %w", utxo.OutPoint.Hash.String(), err)
// 		}
// 		nonWitnessTxHex = hex.EncodeToString(txRaw)
// 	}

// 	psbtUTXOs = append(psbtUTXOs, psbt.PsbtUTXO{
// 		TxID:             utxo.OutPoint.Hash.String(),
// 		Vout:             utxo.OutPoint.Index,
// 		ValueSat:         utxo.Value,
// 		ScriptPubKeyHex:  hex.EncodeToString(pkScript),
// 		NonWitnessTxHex:  nonWitnessTxHex,
// 		RedeemScriptHex:  hex.EncodeToString(addrScriptInfo.RedeemScriptHashHex),
// 		WitnessScriptHex: hex.EncodeToString(addrScriptInfo.WitnessProgramHex),
// 	})
// }
