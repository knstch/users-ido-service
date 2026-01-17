package decoder

import (
	"context"
	"net/http"
	"reflect"

	public "github.com/knstch/users-ido-api/public"
)

// DecodeGoogleOAuthCallbackRequest manually decodes query parameters from query string
// and sets them to the fields of GoogleOAuthCallbackRequest using reflection.
//
// This is needed because form decoder doesn't work properly with protobuf structures
// due to their private fields (state, unknownFields, sizeCache).
func DecodeGoogleOAuthCallbackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")
	errorDescription := r.URL.Query().Get("error_description")
	scope := r.URL.Query().Get("scope")

	req := &public.GoogleOAuthCallbackRequest{}

	// Use reflection to set the fields
	v := reflect.ValueOf(req).Elem()

	if codeField := v.FieldByName("Code"); codeField.IsValid() && codeField.CanSet() {
		codeField.SetString(code)
	}
	if stateField := v.FieldByName("State"); stateField.IsValid() && stateField.CanSet() {
		stateField.SetString(state)
	}
	if errorField := v.FieldByName("Error"); errorField.IsValid() && errorField.CanSet() {
		errorField.SetString(errorParam)
	}
	if errorDescField := v.FieldByName("ErrorDescription"); errorDescField.IsValid() && errorDescField.CanSet() {
		errorDescField.SetString(errorDescription)
	}
	if scopeField := v.FieldByName("Scope"); scopeField.IsValid() && scopeField.CanSet() {
		scopeField.SetString(scope)
	}

	return req, nil
}
