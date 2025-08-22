package mempoolapis

import (
	"context"
	"fmt"
	"path"
)

// 估算交易费率
func (c *Client) EstimateFeeRate(ctx context.Context, targetBlocks int) (*FeeRateDTO, error) {
	u := *c.base
	u.Path = path.Join(u.Path, "/api/v1/fees/recommended")
	var dto FeeRateDTO
	if err := c.getJSON(ctx, u.String(), &dto); err != nil {
		return nil, err
	}

	fmt.Printf("mempooldto: %+v\n", dto)

	return &dto, nil
}
