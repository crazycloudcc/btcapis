package psbt

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/crazycloudcc/btcapis/internal/decoders"
	"github.com/crazycloudcc/btcapis/types"
)

func TestClassifyScript(t *testing.T) {
	tests := []struct {
		name     string
		pkScript []byte
		expected types.AddressType
	}{
		{
			name:     "P2PKH",
			pkScript: []byte{0x76, 0xa9, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0xac},
			expected: types.AddrP2PKH,
		},
		{
			name:     "P2SH",
			pkScript: []byte{0xa9, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x87},
			expected: types.AddrP2SH,
		},
		{
			name:     "P2WPKH",
			pkScript: []byte{0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: types.AddrP2WPKH,
		},
		{
			name:     "P2WSH",
			pkScript: []byte{0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: types.AddrP2WSH,
		},
		{
			name:     "P2TR",
			pkScript: []byte{0x51, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: types.AddrP2TR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decoders.PKScriptToType(tt.pkScript)
			if result != tt.expected {
				t.Errorf("PKScriptToType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseWitnessScriptPubkeys(t *testing.T) {
	// 创建2-of-3多签脚本: OP_2 <pubkey1> <pubkey2> <pubkey3> OP_3 OP_CHECKMULTISIG
	pubkey1 := make([]byte, 33)
	pubkey2 := make([]byte, 33)
	pubkey3 := make([]byte, 33)

	script := []byte{0x52}        // OP_2
	script = append(script, 0x21) // PUSH33
	script = append(script, pubkey1...)
	script = append(script, 0x21) // PUSH33
	script = append(script, pubkey2...)
	script = append(script, 0x21) // PUSH33
	script = append(script, pubkey3...)
	script = append(script, 0x53, 0xae) // OP_3, OP_CHECKMULTISIG

	p := &Packet{}
	pubkeys, err := p.parseWitnessScriptPubkeys(script)
	if err != nil {
		t.Fatalf("parseWitnessScriptPubkeys() error = %v", err)
	}

	if len(pubkeys) != 3 {
		t.Errorf("Expected 3 pubkeys, got %d", len(pubkeys))
	}
}

func TestHash160(t *testing.T) {
	data := []byte("test data")
	result := hash160(data)

	if len(result) != 20 {
		t.Errorf("hash160() returned %d bytes, expected 20", len(result))
	}

	// 验证结果不为空
	allZero := true
	for _, b := range result {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("hash160() returned all zeros")
	}
}

func TestNewV0FromUnsignedTx(t *testing.T) {
	// 创建测试交易
	tx := &wire.MsgTx{
		Version: 1,
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 0,
				},
				Sequence: 0xffffffff,
			},
		},
		TxOut: []*wire.TxOut{
			{
				Value:    1000000,
				PkScript: []byte{0x76, 0xa9, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0xac},
			},
		},
	}

	psbt := NewV0FromUnsignedTx(tx)

	if psbt.Version != VersionV0 {
		t.Errorf("Expected version %d, got %d", VersionV0, psbt.Version)
	}

	if len(psbt.Inputs) != 1 {
		t.Errorf("Expected 1 input, got %d", len(psbt.Inputs))
	}

	if len(psbt.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(psbt.Outputs))
	}

	if psbt.Outputs[0].Value != 1000000 {
		t.Errorf("Expected output value 1000000, got %d", psbt.Outputs[0].Value)
	}

	fmt.Printf("psbt: %+v\n", psbt)
}
