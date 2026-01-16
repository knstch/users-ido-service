package public

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/template-api/public"
)

func MakeHandlerEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.Handle(ctx, request.(*public.Request))
	}
}

func (c *Controller) Handle(ctx context.Context, req *public.Request) (*public.Response, error) {
	ctx, span := tracing.StartSpan(ctx, "public: Handle")
	defer span.End()

	return c.svc.Handle(ctx, req)
}
