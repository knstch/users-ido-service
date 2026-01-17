package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/users-ido-api/public"
)

func MakeRefreshAccessTokenEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.RefreshAccessToken(ctx, request.(*public.RefreshAccessTokenRequest))
	}
}

func (c *Controller) RefreshAccessToken(ctx context.Context,
	req *public.RefreshAccessTokenRequest) (*public.RefreshAccessTokenResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: RefreshAccessToken")
	defer span.End()

	tokens, err := c.svc.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("svc.RefreshAccessToken: %w", err)
	}

	return &public.RefreshAccessTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
