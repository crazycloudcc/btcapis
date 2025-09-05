package btcapis

// import (
// 	"encoding/base64"
// 	"encoding/hex"
// 	"errors"
// 	"strings"
// 	"unicode"
// )

// // 校验地址是否合法
// func (c *Client) ValidateAddress(s string) error {
// 	if s == "" {
// 		return errors.New("address is empty")
// 	}

// 	n := len(s)
// 	if n < 14 || n > 74 {
// 		return errors.New("address length is invalid")
// 	}

// 	if strings.TrimSpace(s) != s {
// 		return errors.New("address has leading or trailing whitespace")
// 	}

// 	for _, r := range s {
// 		if unicode.IsSpace(r) || unicode.IsControl(r) {
// 			return errors.New("address contains whitespace or control characters")
// 		}
// 	}

// 	// Bech32 大小写不能混合
// 	lower, upper := strings.ToLower(s), strings.ToUpper(s)
// 	if s != lower && s != upper {
// 		return errors.New("bech32 address cannot mix upper and lower case")
// 	}

// 	return nil
// }

// // 校验交易id是否合法
// func (c *Client) ValidateTxid(s string) error {
// 	if len(s) != 64 {
// 		return errors.New("invalid txid length")
// 	}

// 	_, err := hex.DecodeString(s)
// 	if err != nil {
// 		return errors.New("txid is not a valid hex string")
// 	}

// 	return nil
// }

// // 校验交易元数据是否合法
// func (c *Client) ValidateRawTx(s string) error {
// 	if err := c.ValidateHex(s); err != nil {
// 		return err
// 	}

// 	_, err := hex.DecodeString(s)
// 	if err != nil {
// 		return errors.New("raw tx is not valid hex")
// 	}

// 	return nil
// }

// // 校验PSBT数据是否合法
// func (c *Client) ValidatePSBTBase64(s string) error {
// 	if s == "" {
// 		return errors.New("psbt is empty")
// 	}

// 	_, err := base64.StdEncoding.DecodeString(s)
// 	if err != nil {
// 		return errors.New("psbt is not valid base64")
// 	}

// 	return nil
// }

// // helper 校验hex字符串是否合法
// func (c *Client) ValidateHex(s string) error {
// 	if s == "" {
// 		return errors.New("hex string is empty")
// 	}

// 	if len(s)%2 != 0 {
// 		return errors.New("hex string length must be even")
// 	}

// 	if _, err := hex.DecodeString(s); err != nil {
// 		return errors.New("hex string is not valid hex")
// 	}

// 	return nil
// }

// // helper 校验base64字符串是否合法
// func (c *Client) ValidateBase64String(s string) error {
// 	if s == "" {
// 		return errors.New("base64 string is empty")
// 	}

// 	if _, err := base64.StdEncoding.DecodeString(s); err != nil {
// 		return errors.New("string is not valid base64")
// 	}

// 	return nil
// }
