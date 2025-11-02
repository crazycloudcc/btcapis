package address

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/crazycloudcc/btcapis/pkg/logger"
	"github.com/crazycloudcc/btcapis/types"
	"golang.org/x/crypto/ripemd160"
)

func (c *Client) GenerateNew() (*types.WalletInfo, error) {
	return generateWallet(types.CurrentNetworkParams)
}

// generateWallet 生成单个钱包
func generateWallet(params *chaincfg.Params) (*types.WalletInfo, error) {
	// 1) 256-bit 熵 → 24 词助记词
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	// 2) 可选：加入 BIP39 passphrase（相当于第25词 加强抗窃取能力；空串也合法）
	passphrase := "" // 建议实际部署改为用户自定
	seed := bip39.NewSeed(mnemonic, passphrase)

	// 3. 从 seed 得到 Master Key
	master, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		return nil, err
	}

	// ========= 四类地址（index=0） =========
	p2pkh, err := addrBIP44_P2PKH(master, params, 0) // 1...
	if err != nil {
		return nil, err
	}
	p2sh, err := addrBIP49_P2SH_P2WPKH(master, params, 0) // 3...
	if err != nil {
		return nil, err
	}
	p2wpkh, err := addrBIP84_P2WPKH(master, params, 0) // bc1q...
	if err != nil {
		return nil, err
	}
	p2tr, err := addrBIP86_Taproot(master, params, 0) // bc1p...
	if err != nil {
		return nil, err
	}

	return &types.WalletInfo{
		P2PKH:    p2pkh,
		P2PSH:    p2sh,
		P2WPKH:   p2wpkh,
		P2TR:     p2tr,
		Mnemonic: mnemonic,
		XPRV:     master.String(),
	}, nil
}

// ======= BIP44: m/44'/0'/0'/0/index -> 1... =======
func addrBIP44_P2PKH(master *hdkeychain.ExtendedKey, params *chaincfg.Params, index uint32) (string, error) {
	key, err := derivePath(master, []uint32{
		hdkeychain.HardenedKeyStart + 44,
		hdkeychain.HardenedKeyStart + 0,
		hdkeychain.HardenedKeyStart + 0,
		0, index,
	})
	if err != nil {
		return "", err
	}
	pub, err := key.Neuter()
	if err != nil {
		return "", err
	}
	pk, err := pub.ECPubKey()
	if err != nil {
		return "", err
	}
	h160 := hash160(pk.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(h160, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// ======= BIP49: m/49'/0'/0'/0/index -> 3... (P2SH-P2WPKH) =======
func addrBIP49_P2SH_P2WPKH(master *hdkeychain.ExtendedKey, params *chaincfg.Params, index uint32) (string, error) {
	key, err := derivePath(master, []uint32{
		hdkeychain.HardenedKeyStart + 49,
		hdkeychain.HardenedKeyStart + 0,
		hdkeychain.HardenedKeyStart + 0,
		0, index,
	})
	if err != nil {
		return "", err
	}
	pub, err := key.Neuter()
	if err != nil {
		return "", err
	}
	pk, err := pub.ECPubKey()
	if err != nil {
		return "", err
	}
	// redeemScript = 0 <20-byte-pubkey-hash>
	wpkh := hash160(pk.SerializeCompressed())
	redeem := buildWitnessProgramV0(wpkh) // OP_0 0x14 <20-byte>
	sh := hash160(redeem)                 // HASH160(redeemScript)
	addr, err := btcutil.NewAddressScriptHashFromHash(sh, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// ======= BIP84: m/84'/0'/0'/0/index -> bc1q... (P2WPKH) =======
func addrBIP84_P2WPKH(master *hdkeychain.ExtendedKey, params *chaincfg.Params, index uint32) (string, error) {
	key, err := derivePath(master, []uint32{
		hdkeychain.HardenedKeyStart + 84,
		hdkeychain.HardenedKeyStart + 0,
		hdkeychain.HardenedKeyStart + 0,
		0, index,
	})
	if err != nil {
		return "", err
	}
	pub, err := key.Neuter()
	if err != nil {
		return "", err
	}
	pk, err := pub.ECPubKey()
	if err != nil {
		return "", err
	}
	h160 := hash160(pk.SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(h160, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// ======= BIP86: m/86'/0'/0'/0/index -> bc1p... (Taproot, bech32m) =======
// 需要 Tweaked x-only pubkey
func addrBIP86_Taproot(master *hdkeychain.ExtendedKey, params *chaincfg.Params, index uint32) (string, error) {
	// 对 Taproot，优先使用私钥派生，拿到内公钥再做 tweak
	key, err := derivePath(master, []uint32{
		hdkeychain.HardenedKeyStart + 86,
		hdkeychain.HardenedKeyStart + 0,
		hdkeychain.HardenedKeyStart + 0,
		0, index,
	})
	if err != nil {
		return "", err
	}
	// child private key
	sk, err := key.ECPrivKey()
	if err != nil {
		return "", err
	}
	internal := sk.PubKey() // *btcec.PublicKey

	// 优先尝试库函数（若你的 btcd/txscript 版本提供） ---
	tweaked := txscript.ComputeTaprootKeyNoScript(internal) // 某些版本提供
	xOnlyTweaked := tweaked.X().Bytes()

	addr, err := btcutil.NewAddressTaproot(xOnlyTweaked, params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// SLIP-0010 master key for Ed25519
func slip10MasterKeyEd25519(seed []byte) ([]byte, []byte) {
	mac := hmac.New(sha512.New, []byte("ed25519 seed"))
	mac.Write(seed)
	I := mac.Sum(nil)
	return I[:32], I[32:]
}

// Child key derivation for Ed25519 (hardened only)
func slip10CKDPrivEd25519(k, c []byte, i uint32) ([]byte, []byte) {
	// data: 0x00 || k || ser32(i)
	data := make([]byte, 0, 1+32+4)
	data = append(data, 0x00)
	data = append(data, k...)
	var iBytes [4]byte
	binary.BigEndian.PutUint32(iBytes[:], i)
	data = append(data, iBytes[:]...)

	mac := hmac.New(sha512.New, c)
	mac.Write(data)
	I := mac.Sum(nil)
	return I[:32], I[32:]
}

// ======= 工具函数 =======

// derivePath 依次派生（允许混合 hardened/normal）
func derivePath(key *hdkeychain.ExtendedKey, path []uint32) (*hdkeychain.ExtendedKey, error) {
	k := key
	var err error
	for _, i := range path {
		k, err = k.Derive(i)
		if err != nil {
			return nil, err
		}
	}
	return k, nil
}

// HASH160 = RIPEMD160(SHA256(data))
func hash160(data []byte) []byte {
	sha := sha256.Sum256(data)
	r := ripemd160.New()
	_, _ = r.Write(sha[:])
	return r.Sum(nil)
}

// 构造 P2WPKH 见证程序: OP_0 <20-byte>
func buildWitnessProgramV0(hash20 []byte) []byte {
	if len(hash20) != 20 {
		panic("wpkh must be 20 bytes")
	}
	// OP_0 (0x00) + push_20 (0x14) + hash20
	return append([]byte{txscript.OP_0, 0x14}, hash20...)
}

// TapTweak = taggedHash("TapTweak", xOnly || merkleRoot?)
// 返回 32 字节
func tapTweakTaggedHash(xOnly []byte, merkleRoot []byte) [32]byte {
	tag := []byte("TapTweak")
	tagHash := sha256.Sum256(tag)
	h := sha256.New()
	// tagged hash: H(H(tag)||H(tag)||msg)
	h.Write(tagHash[:])
	h.Write(tagHash[:])
	h.Write(xOnly)
	if merkleRoot != nil {
		h.Write(merkleRoot)
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

func must(err error) {
	if err != nil {
		logger.Error("Error: %v", err)
	}
}
