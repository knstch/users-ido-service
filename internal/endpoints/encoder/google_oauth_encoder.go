package encoder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	public "github.com/knstch/users-ido-api/public"
)

type GoogleOAuthCallbackHTTPResponse struct {
	AccessToken  string
	RefreshToken string
	RedirectURL  string
}

// EncodeGoogleOAuthCallbackResponse sets auth cookies and redirects to the platform.
//
// It sets HttpOnly cookies for access/refresh tokens and returns HTTP 302 with
// Location pointing to RedirectUrl. No JSON body is written.
func EncodeGoogleOAuthCallbackResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(*public.GoogleOAuthCallbackResponse)
	if !ok {
		return fmt.Errorf("unexpected response type %T", response)
	}

	platformDomain, secure := cookieDomainAndSecureFromPlatform(resp.RedirectUrl)

	// Access token cookie (short-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   platformDomain,
	})

	// Refresh token cookie (longer-lived; max-age left to browser session unless you want persistence)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   platformDomain,
	})

	// Redirect to the original page.
	w.Header().Set("Location", resp.RedirectUrl)
	w.WriteHeader(http.StatusFound)
	return nil
}

// cookieDomainAndSecureFromPlatform derives cookie domain and secure flag from redirect URL.
// Domain is best-effort: for localhost we skip Domain, otherwise we set the host (without port).
func cookieDomainAndSecureFromPlatform(redirect string) (domain string, secure bool) {
	u, err := url.Parse(redirect)
	if err != nil {
		return "", false
	}
	secure = u.Scheme == "https"
	host := u.Hostname()
	if host == "" || host == "localhost" {
		return "", secure
	}
	// allow subdomain sharing by setting ".example.com" if host is "app.example.com"
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return "." + strings.Join(parts[len(parts)-2:], "."), secure
	}
	return host, secure
}
