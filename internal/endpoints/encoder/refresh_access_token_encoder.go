package encoder

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	public "github.com/knstch/users-ido-api/public"
)

// EncodeRefreshAccessTokenResponse sets auth cookies and returns refreshed tokens.
//
// It sets HttpOnly cookies for access/refresh tokens and writes a JSON body with
// the same tokens. Status code is 200.
func EncodeRefreshAccessTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(*public.RefreshAccessTokenResponse)
	if !ok {
		return fmt.Errorf("unexpected response type %T", response)
	}

	domain, secure := cookieDomainAndSecureFromRequestContext(ctx)

	deleteOldCookies(w, domain, secure)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.GetAccessToken(),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
		MaxAge:   3600,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.GetRefreshToken(),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
		MaxAge:   604800,
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func deleteOldCookies(w http.ResponseWriter, currentDomain string, secure bool) {
	cookieNames := []string{"access_token", "refresh_token"}
	domains := []string{"", currentDomain}
	if currentDomain != "" && !strings.HasPrefix(currentDomain, ".") {
		domains = append(domains, "."+currentDomain)
	}
	
	for _, name := range cookieNames {
		for _, domain := range domains {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   secure,
				SameSite: http.SameSiteLaxMode,
				Domain:   domain,
				MaxAge:   -1,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   secure,
				SameSite: http.SameSiteNoneMode,
				Domain:   domain,
				MaxAge:   -1,
			})
		}
	}
}

func cookieDomainAndSecureFromRequestContext(ctx context.Context) (domain string, secure bool) {
	xfp, _ := ctx.Value(httptransport.ContextKeyRequestXForwardedProto).(string)
	if xfp == "" {
		// fallback: if not behind proxy, assume http
		xfp = "http"
	}
	secure = xfp == "https"

	host, _ := ctx.Value(httptransport.ContextKeyRequestHost).(string)
	host = strings.TrimSpace(host)
	if host == "" {
		return "", secure
	}

	hostname := host
	if h, _, err := net.SplitHostPort(host); err == nil {
		hostname = h
	}

	if hostname == "" || hostname == "localhost" || net.ParseIP(hostname) != nil {
		return "", secure
	}

	parts := strings.Split(hostname, ".")
	if len(parts) >= 2 {
		return "." + strings.Join(parts[len(parts)-2:], "."), secure
	}
	return hostname, secure
}
