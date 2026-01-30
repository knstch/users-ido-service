package users

// OAuthState is encoded into OAuth "state" parameter and later validated on callback.
// It contains a CSRF token and the original return URL/path.
type OAuthState struct {
	CSRF   string `json:"csrf"`
	Return string `json:"return"`
	Scheme string `json:"scheme"`
}

type GoogleIDTokenClaims struct {
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
}
