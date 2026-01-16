package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/tracing"

	public "github.com/knstch/users-ido-api/public"
)

func MakeAuthViaGoogleEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.AuthViaGoogle(ctx, request.(*public.AuthViaGoogleRequest))
	}
}

func (c *Controller) AuthViaGoogle(ctx context.Context, req *public.AuthViaGoogleRequest) (*public.AuthViaGoogleResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: AuthViaGoogle")
	defer span.End()

	loginURL, err := c.svc.AuthViaGoogle(ctx, req.GetState())
	if err != nil {
		return nil, fmt.Errorf("svc.AuthViaGoogle: %w", err)
	}

	return &public.AuthViaGoogleResponse{
		LoginUrl: loginURL,
	}, nil
}
