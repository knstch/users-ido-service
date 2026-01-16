package public

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/users-ido-api/public"
)

func MakeGoogleOAuthCallbackEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.GoogleOAuthCallback(ctx, request.(*public.GoogleOAuthCallbackRequest))
	}
}

func (c *Controller) GoogleOAuthCallback(ctx context.Context, req *public.GoogleOAuthCallbackRequest) (*public.GoogleOAuthCallbackResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: GoogleOAuthCallback")
	defer span.End()

	tokens, returnURL, err := c.svc.CompleteLogin(ctx, req.GetState(), req.GetCode())
	if err != nil {
		return nil, fmt.Errorf("svc.CompleteLogin: %w", err)
	}

	redirectURL, err := buildRedirectURL(c.cfg.PlatformURL, returnURL)
	if err != nil {
		return nil, fmt.Errorf("buildRedirectURL: %w", err)
	}

	return &public.GoogleOAuthCallbackResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		RedirectUrl:  redirectURL,
	}, nil
}

func buildRedirectURL(platformURL string, returnURL string) (string, error) {
	if returnURL == "" {
		return "", fmt.Errorf("empty returnURL")
	}

	if strings.HasPrefix(returnURL, "http://") || strings.HasPrefix(returnURL, "https://") {
		return returnURL, nil
	}

	base, err := url.Parse(platformURL)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("invalid platformURL")
	}

	if !strings.HasPrefix(returnURL, "/") {
		returnURL = "/" + returnURL
	}

	base.Path = strings.TrimRight(base.Path, "/") + returnURL
	base.RawQuery = ""
	base.Fragment = ""
	return base.String(), nil
}
