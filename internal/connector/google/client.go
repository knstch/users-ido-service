package google

import (
	"context"
	"net/http"

	"github.com/knstch/knstch-libs/log"

	"users-service/config"
)

type Client interface {
	// ExchangeCodeToToken exchanges an OAuth authorization code for tokens.
	ExchangeCodeToToken(ctx context.Context, req ExchangeCodeToTokenRequest) (*ExchangeCodeToTokenResponse, error)
}

// GetClient constructs a Google OAuth connector client.
func GetClient(cfg *config.Config, lg *log.Logger) (*ClientImpl, error) {
	httpClient := http.DefaultClient

	client, err := NewClient(lg, cfg, httpClient)
	if err != nil {
		return nil, err
	}

	return client, nil
}
