package public

import (
	"net/http"

	"users-service/config"
	"users-service/internal/users"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/knstch/knstch-libs/log"

	public "github.com/knstch/template-api/public"

	"github.com/knstch/knstch-libs/endpoints"
	"github.com/knstch/knstch-libs/transport"
)

type Controller struct {
	svc users.Service
	lg  *log.Logger
	cfg *config.Config

	public.UnimplementedTemplateServer
}

func NewController(svc users.Service, lg *log.Logger, cfg *config.Config) *Controller {
	return &Controller{
		svc: svc,
		cfg: cfg,
		lg:  lg,
	}
}

func (c *Controller) Endpoints() []endpoints.Endpoint {
	return []endpoints.Endpoint{
		{
			Method:  http.MethodPost,
			Path:    "/handler",
			Handler: MakeHandlerEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.Request],
			Encoder: httptransport.EncodeJSONResponse,
		},
	}
}
