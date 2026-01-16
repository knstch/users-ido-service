package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"gorm.io/gorm"

	"users-service/internal/domain/dto"
)

func (r *DBRepo) GetAccessTokens(ctx context.Context, filter AccessTokenFilter) (dto.AccessTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetAccessTokens")
	defer span.End()

	var tokenPair AccessToken
	if err := r.db.WithContext(ctx).Scopes(filter.ToScope()).First(&tokenPair).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AccessTokens{}, fmt.Errorf("db.First: %w", svcerrs.ErrDataNotFound)
		}
		return dto.AccessTokens{}, fmt.Errorf("db.First: %w", err)
	}

	return dto.AccessTokens{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
