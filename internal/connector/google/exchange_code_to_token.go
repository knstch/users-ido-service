package google

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"
)

func (c *ClientImpl) ExchangeCodeToToken(ctx context.Context, req ExchangeCodeToTokenRequest) (*ExchangeCodeToTokenResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "google: ExchangeCodeToToken")
	defer span.End()

	resp, err := c.endpoints.ExchangeCodeToToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("endpoints.ExchangeCodeToToken: %w", err)
	}

	return resp.(*ExchangeCodeToTokenResponse), nil
}
