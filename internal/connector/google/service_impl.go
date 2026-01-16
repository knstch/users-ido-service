package google

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/knstch/knstch-libs/log"

	"users-service/config"
)

type ClientImpl struct {
	lg         *log.Logger
	cfg        *config.Config
	httpClient *http.Client
	endpoints  struct {
		ExchangeCodeToToken endpoint.Endpoint
	}
}

func NewClient(lg *log.Logger, cfg *config.Config, client *http.Client) (*ClientImpl, error) {
	cl := &ClientImpl{
		lg:         lg,
		cfg:        cfg,
		httpClient: client,
	}

	if err := cl.initEndpoints(); err != nil {
		return nil, fmt.Errorf("initEndpoints: %w", err)
	}

	return cl, nil
}

func (c *ClientImpl) initEndpoints() error {
	if err := c.initExchangeCodeToTokenEndpoint(); err != nil {
		return fmt.Errorf("initExchangeCodeToTokenEndpoint: %w", err)
	}

	return nil
}

func (c *ClientImpl) initExchangeCodeToTokenEndpoint() error {
	googleApiURL, err := url.Parse(c.cfg.GoogleAPI.GoogleAPIHost + "/token")
	if err != nil {
		return fmt.Errorf("url.Parse: %w", err)
	}

	opts := []kitHttp.ClientOption{
		kitHttp.SetClient(c.httpClient),
	}

	c.endpoints.ExchangeCodeToToken = kitHttp.NewClient(
		http.MethodPost,
		googleApiURL,
		EncodeExchangeCodeToTokenRequest,
		DecodeExchangeCodeToTokenResponse,
		opts...,
	).Endpoint()

	return nil
}
