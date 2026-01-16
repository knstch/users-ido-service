package users

type OAuthState struct {
	CSRF   string `json:"csrf"`
	Return string `json:"return"`
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
