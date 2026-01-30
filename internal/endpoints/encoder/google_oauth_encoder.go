package encoder

import (
	"context"
	"fmt"
	"net/http"

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

	domain, secure := cookieDomainAndSecureFromRequestContext(ctx)

	deleteOldCookies(w, domain, secure)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
		Domain:   domain,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
		MaxAge:   604800,
	})

	w.Header().Set("Location", resp.RedirectUrl)
	w.WriteHeader(http.StatusFound)
	return nil
}
