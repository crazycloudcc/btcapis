package tx

import "context"

// 将给定私钥+对应fromAddress的所有余额转移到toAddress
// 这是一个紧急避险功能，用于将泄露私钥的地址的余额快速转移到安全地址
func (c *Client) TransferAllToNewAddress(
	ctx context.Context,
	toAddress string,
	privateKeyWIF string,
	fromAddress string) (string, error) {
}
