package decoder

import (
	"context"
	"net/http"
	"reflect"

	public "github.com/knstch/users-ido-api/public"
)

// DecodeAuthViaGoogleRequest manually decodes the location parameter from query string
// and sets it to the Location field of AuthViaGoogleRequest using reflection.
//
// This is needed because form decoder doesn't work properly with protobuf structures
// due to their private fields (state, unknownFields, sizeCache).
func DecodeAuthViaGoogleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	location := r.URL.Query().Get("location")

	req := &public.AuthViaGoogleRequest{}

	// Use reflection to set the Location field
	v := reflect.ValueOf(req).Elem()
	locationField := v.FieldByName("Location")
	if locationField.IsValid() && locationField.CanSet() {
		locationField.SetString(location)
	}

	return req, nil
}
