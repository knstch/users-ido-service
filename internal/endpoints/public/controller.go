package public

import (
	"net/http"

	"github.com/knstch/knstch-libs/middleware"
	"github.com/knstch/knstch-libs/transport"

	"users-service/internal/endpoints/decoder"

	"users-service/config"
	"users-service/internal/endpoints/encoder"
	"users-service/internal/users"

	"github.com/knstch/knstch-libs/log"

	httptransport "github.com/go-kit/kit/transport/http"
	public "github.com/knstch/users-ido-api/public"

	"github.com/knstch/knstch-libs/endpoints"
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
	defaultMiddlewares := []middleware.Middleware{middleware.WithCookieAuth(c.cfg.JwtSecret)}

	return []endpoints.Endpoint{
		{
			Method:  http.MethodGet,
			Path:    "/authViaGoogle",
			Handler: MakeAuthViaGoogleEndpoint(c),
			Decoder: decoder.DecodeAuthViaGoogleRequest,
			Encoder: encoder.EncodeAuthViaGoogleResponse,
		},
		{
			Method:  http.MethodGet,
			Path:    "/googleOAuthCallback",
			Handler: MakeGoogleOAuthCallbackEndpoint(c),
			Decoder: decoder.DecodeGoogleOAuthCallbackRequest,
			Encoder: encoder.EncodeGoogleOAuthCallbackResponse,
		},
		{
			Method:  http.MethodPost,
			Path:    "/refreshAccessToken",
			Handler: MakeRefreshAccessTokenEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.RefreshAccessTokenRequest],
			Encoder: encoder.EncodeRefreshAccessTokenResponse,
		},
		{
			Method:  http.MethodGet,
			Path:    "/getUser",
			Handler: MakeGetUserEndpoint(c),
			Decoder: transport.DecodeQueryRequest[public.GetUserRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     defaultMiddlewares,
		},
	}
}
