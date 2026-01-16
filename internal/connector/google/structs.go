package google

type ExchangeCodeToTokenRequest struct {
	ClientSecret   string
	Code           string
	GoogleClientID string
	RedirectURI    string
}

type ExchangeCodeToTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}
