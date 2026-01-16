package encoder

import (
	"context"
	"fmt"
	"net/http"

	public "github.com/knstch/users-ido-api/public"
)

func EncodeAuthViaGoogleResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(*public.AuthViaGoogleResponse)
	if !ok {
		return fmt.Errorf("unexpected response type %T", response)
	}
	if resp.GetLoginUrl() == "" {
		return fmt.Errorf("empty login url")
	}

	w.Header().Set("Location", resp.GetLoginUrl())
	w.WriteHeader(http.StatusFound)
	return nil
}
