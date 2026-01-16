package google

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeExchangeCodeToTokenResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp ExchangeCodeToTokenResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
