package public

import (
	"context"
	"fmt"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/tracing"

	public "github.com/knstch/users-ido-api/public"
)

func MakeAuthViaGoogleEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.AuthViaGoogle(ctx, request.(*public.AuthViaGoogleRequest))
	}
}

// AuthViaGoogle is an HTTP endpoint handler that starts Google OAuth.
//
// The request field `state` is treated as the return URL/path from which the
// user initiated login. The response contains a login URL which is then used by
// the HTTP encoder to perform a redirect.
func (c *Controller) AuthViaGoogle(ctx context.Context, req *public.AuthViaGoogleRequest) (*public.AuthViaGoogleResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: AuthViaGoogle")
	defer span.End()

	xfp, _ := ctx.Value(httptransport.ContextKeyRequestXForwardedProto).(string)
	if xfp == "" {
		xfp = "http"
	}

	loginURL, err := c.svc.AuthViaGoogle(ctx, req.GetLocation(), xfp)
	if err != nil {
		return nil, fmt.Errorf("svc.AuthViaGoogle: %w", err)
	}

	return &public.AuthViaGoogleResponse{
		LoginUrl: loginURL,
	}, nil
}
