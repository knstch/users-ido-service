package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"

	"users-service/internal/users/utils"
	"users-service/internal/users/validator"
)

const (
	randomStringLength = 16

	authExpirationTime = time.Minute * 15
)

func (s *ServiceImpl) AuthViaGoogle(ctx context.Context, stateURL string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "users: AuthViaGoogle")
	defer span.End()

	if safe := validator.IsSafeRedirectURL(stateURL, s.cfg.PlatformURL); !safe {
		return "", fmt.Errorf("unknown return url: %w", svcerrs.ErrInvalidData)
	}

	securityCode, err := utils.RandomString(randomStringLength)
	if err != nil {
		return "", fmt.Errorf("utils.RandomString: %w", err)
	}

	stateJSON, err := json.Marshal(&OAuthState{
		CSRF:   securityCode,
		Return: stateURL,
	})
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	state := base64.RawURLEncoding.EncodeToString(stateJSON)

	loginURL, err := buildGoogleAuthURL(
		s.cfg.GoogleAPI.GoogleClientID,
		s.cfg.GoogleAPI.GoogleRedirectURI,
		state,
		s.cfg.GoogleAPI.GoogleAuthHost,
	)
	if err != nil {
		return "", err
	}

	if err = s.redis.Set(ctx, "oauth:state:"+securityCode, 1, authExpirationTime).Err(); err != nil {
		return "", fmt.Errorf("redis.Set: %w", err)
	}

	return loginURL, nil
}

func buildGoogleAuthURL(clientID, redirectURI, state, authURL string) (string, error) {
	if clientID == "" {
		return "", fmt.Errorf("client_id is required: %w", svcerrs.ErrInvalidData)
	}
	if redirectURI == "" {
		return "", fmt.Errorf("redirect_uri is required: %w", svcerrs.ErrInvalidData)
	}
	if state == "" {
		return "", fmt.Errorf("state is required: %w", svcerrs.ErrInvalidData)
	}

	u, err := url.Parse(authURL)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile")
	q.Set("state", state)
	q.Set("access_type", "offline")
	q.Set("prompt", "consent")

	u.RawQuery = q.Encode()

	return u.String(), nil
}
