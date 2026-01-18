package repo

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"users-service/internal/users/models"
)

func (r *DBRepo) CreateAccessTokens(ctx context.Context, accessToken, refreshToken string, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "repo: CreateAccessToken")
	defer span.End()

	if err := r.db.WithContext(ctx).Model(&models.AccessToken{}).Create(&models.AccessToken{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}).Error; err != nil {
		return fmt.Errorf("db.Create: %w", err)
	}

	return nil
}
