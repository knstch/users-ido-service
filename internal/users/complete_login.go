package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"github.com/redis/go-redis/v9"

	"users-service/internal/connector/google"
	"users-service/internal/domain/dto"
	"users-service/internal/domain/enum"
	"users-service/internal/users/repo"
)

func parseOAuthState(state string) (OAuthState, error) {
	if state == "" {
		return OAuthState{}, fmt.Errorf("state is empty")
	}

	raw, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return OAuthState{}, fmt.Errorf("base64.RawURLEncoding.DecodeStringe: %w", err)
	}

	var parsedState OAuthState
	if err = json.Unmarshal(raw, &parsedState); err != nil {
		return OAuthState{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	if parsedState.CSRF == "" {
		return OAuthState{}, fmt.Errorf("state csrf is empty")
	}
	if parsedState.Return == "" {
		return OAuthState{}, fmt.Errorf("state return is empty")
	}

	return parsedState, nil
}

func parseJWTClaimsNoVerify(idToken string) (GoogleIDTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return GoogleIDTokenClaims{}, fmt.Errorf("invalid jwt format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return GoogleIDTokenClaims{}, fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", err)
	}

	var claims GoogleIDTokenClaims
	if err = json.Unmarshal(payloadBytes, &claims); err != nil {
		return GoogleIDTokenClaims{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return claims, nil
}

func (s *ServiceImpl) CompleteLogin(ctx context.Context, state, code string) (dto.AccessTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "users: CompleteLogin")
	defer span.End()

	parsedState, err := parseOAuthState(state)
	if err != nil {
		return dto.AccessTokens{}, err
	}

	if err = s.redis.Get(ctx, parsedState.CSRF).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.AccessTokens{}, fmt.Errorf("CSRF token not found")
		}
		return dto.AccessTokens{}, fmt.Errorf("redis.Get: %w", err)
	}

	resp, err := s.google.ExchangeCodeToToken(ctx, google.ExchangeCodeToTokenRequest{
		ClientSecret:   s.cfg.GoogleAPI.GoogleOAuthClientSecret,
		Code:           code,
		GoogleClientID: s.cfg.GoogleAPI.GoogleClientID,
		RedirectURI:    parsedState.Return,
	})
	if err != nil {
		return dto.AccessTokens{}, fmt.Errorf("google.ExchangeCodeToToken: %w", err)
	}

	claims, err := parseJWTClaimsNoVerify(resp.IDToken)
	if err != nil {
		return dto.AccessTokens{}, fmt.Errorf("parseJWTClaimsNoVerify: %w", err)
	}

	fullName := strings.Split(claims.Name, " ")
	firstName := ""
	lastName := ""
	if len(fullName) != 2 {
		firstName = "Lovely"
		lastName = fullName[0]
	} else {
		firstName = fullName[0]
		lastName = fullName[1]
	}

	userExists := true
	userHasToBeUpdated := false
	user, err := s.repo.GetUser(ctx, repo.UserFilter{
		Email: claims.Email,
	})
	if err != nil {
		if errors.Is(err, svcerrs.ErrDataNotFound) {
			userExists = false
		} else {
			return dto.AccessTokens{}, fmt.Errorf("repo.GetUser: %w", err)
		}
	}
	if user.FirstName != firstName || user.LastName != lastName || user.ProfilePicture != claims.Picture {
		userHasToBeUpdated = true
	}

	var accessTokens dto.AccessTokens
	if err = s.repo.Transaction(func(st repo.Repository) error {
		var userID uint64
		if !userExists {
			userID, err = st.CreateUser(ctx, claims.Sub, claims.Email, firstName, lastName, claims.Picture)
			if err != nil {
				return fmt.Errorf("st.CreateUser: %w", err)
			}
		}
		if userHasToBeUpdated {
			if err = st.UpdateUserMetadata(ctx, user.ID, firstName, lastName, claims.Picture); err != nil {
				return fmt.Errorf("st.UpdateUserMetadata: %w", err)
			}
		}

		if userExists {
			userID = user.ID
		}
		accessTokens, err = s.mintJWT(userID, enum.User)
		if err != nil {
			return fmt.Errorf("mintJWT: %w", err)
		}

		if err = st.CreateAccessTokens(ctx, accessTokens.AccessToken, accessTokens.RefreshToken, userID); err != nil {
			return fmt.Errorf("st.CreateAccessTokens: %w", err)
		}

		return nil
	}); err != nil {
		return dto.AccessTokens{}, fmt.Errorf("repo.Transaction: %w", err)
	}

	return accessTokens, nil
}
