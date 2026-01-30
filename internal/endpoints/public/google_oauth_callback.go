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

// GoogleOAuthCallback handles Google's redirect back to the service.
//
// It exchanges the `code` for tokens, mints service JWTs and returns a response
// containing access/refresh tokens and the final redirect URL. The HTTP encoder
// then sets cookies and performs a redirect to RedirectUrl.
func (c *Controller) GoogleOAuthCallback(ctx context.Context, req *public.GoogleOAuthCallbackRequest) (*public.GoogleOAuthCallbackResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: GoogleOAuthCallback")
	defer span.End()

	tokens, returnURL, scheme, err := c.svc.CompleteLogin(ctx, req.GetState(), req.GetCode())
	if err != nil {
		return nil, fmt.Errorf("svc.CompleteLogin: %w", err)
	}

	redirectURL, err := buildRedirectURL(c.cfg.PlatformURL, returnURL, scheme)
	if err != nil {
		return nil, fmt.Errorf("buildRedirectURL: %w", err)
	}

	return &public.GoogleOAuthCallbackResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		RedirectUrl:  redirectURL,
	}, nil
}

// buildRedirectURL converts a returnURL (absolute or path) into an absolute URL
// under platformURL. It is used to redirect the browser to the platform domain
// after authentication is complete.
func buildRedirectURL(platformURL string, returnURL string, scheme string) (string, error) {
	if returnURL == "" {
		return "", fmt.Errorf("empty returnURL")
	}

	if strings.HasPrefix(returnURL, "http://") || strings.HasPrefix(returnURL, "https://") {
		return returnURL, nil
	}

	if scheme != "http" && scheme != "https" {
		scheme = "http"
	}

	base, err := url.Parse(platformURL)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}
	if base.Host == "" {
		return "", fmt.Errorf("invalid platformURL")
	}

	if !strings.HasPrefix(returnURL, "/") {
		returnURL = "/" + returnURL
	}

	redirectURL := fmt.Sprintf("%s://%s%s", scheme, base.Host, returnURL)
	return redirectURL, nil
}
