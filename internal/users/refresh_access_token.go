package users

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"users-service/internal/domain/dto"
	"users-service/internal/domain/enum"
	"users-service/internal/users/filters"
	"users-service/internal/users/repo"
)

func (s *ServiceImpl) RefreshAccessToken(ctx context.Context, refreshToken string) (dto.AccessTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "users: RefreshAccessToken")
	defer span.End()

	oldTokenPair, err := s.repo.GetAccessTokens(ctx, filters.AccessTokenFilter{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return dto.AccessTokens{}, fmt.Errorf("repo.GetAccessTokens: %w", err)
	}

	newTokenPair, err := s.mintJWT(oldTokenPair.UserID, enum.User)
	if err != nil {
		return dto.AccessTokens{}, fmt.Errorf("mintJWT: %w", err)
	}

	if err = s.repo.Transaction(func(st repo.Repository) error {
		if err = st.RevokeAccessToken(ctx, refreshToken); err != nil {
			return fmt.Errorf("st.RevokeAccessToken: %w", err)
		}

		if err = st.CreateAccessTokens(ctx, newTokenPair.AccessToken, newTokenPair.RefreshToken, oldTokenPair.UserID); err != nil {
			return fmt.Errorf("st.CreateAccessTokens: %w", err)
		}

		return nil
	}); err != nil {
		return dto.AccessTokens{}, fmt.Errorf("repo.Transaction: %w", err)
	}

	return newTokenPair, nil
}
