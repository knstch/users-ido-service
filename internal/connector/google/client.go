package google

import (
	"context"
	"net/http"

	"github.com/knstch/knstch-libs/log"

	"users-service/config"
)

type Client interface {
	ExchangeCodeToToken(ctx context.Context, req ExchangeCodeToTokenRequest) (*ExchangeCodeToTokenResponse, error)
}

func GetClient(cfg *config.Config, lg *log.Logger) (*ClientImpl, error) {
	httpClient := http.DefaultClient

	client, err := NewClient(lg, cfg, httpClient)
	if err != nil {
		return nil, err
	}

	return client, nil
}
