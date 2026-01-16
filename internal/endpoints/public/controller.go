package public

import (
	"net/http"

	"users-service/config"
	"users-service/internal/endpoints/encoder"
	"users-service/internal/users"

	"github.com/knstch/knstch-libs/log"

	public "github.com/knstch/users-ido-api/public"

	"github.com/knstch/knstch-libs/endpoints"
	"github.com/knstch/knstch-libs/transport"
)

type Controller struct {
	svc users.Service
	lg  *log.Logger
	cfg *config.Config

	public.UnimplementedUsersServer
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
			Method:  http.MethodGet,
			Path:    "/authViaGoogle",
			Handler: MakeAuthViaGoogleEndpoint(c),
			Decoder: transport.DecodeQueryRequest[public.AuthViaGoogleRequest],
			Encoder: encoder.EncodeAuthViaGoogleResponse,
		},
		{
			Method:  http.MethodGet,
			Path:    "/googleOAuthCallback",
			Handler: MakeGoogleOAuthCallbackEndpoint(c),
			Decoder: transport.DecodeQueryRequest[public.GoogleOAuthCallbackRequest],
			Encoder: encoder.EncodeGoogleOAuthCallbackResponse,
		},
	}
}
