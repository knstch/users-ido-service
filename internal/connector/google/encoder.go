package google

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func EncodeExchangeCodeToTokenRequest(_ context.Context, r *http.Request, request interface{}) error {
	req, ok := request.(ExchangeCodeToTokenRequest)
	if !ok {
		return fmt.Errorf("request does not implement ExchangeCodeToTokenRequest")
	}

	form := url.Values{}
	form.Set("code", req.Code)
	form.Set("client_id", req.GoogleClientID)
	form.Set("client_secret", req.ClientSecret)
	form.Set("redirect_uri", req.RedirectURI)
	form.Set("grant_type", "authorization_code")

	encodedForm := form.Encode()
	r.Body = io.NopCloser(strings.NewReader(form.Encode()))
	r.ContentLength = int64(len(encodedForm))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return nil
}
