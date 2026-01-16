package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/knstch/knstch-libs/tracing"
)

func (r *DBRepo) RevokeAccessToken(ctx context.Context, refreshToken string) error {
	ctx, span := tracing.StartSpan(ctx, "repo: RevokeAccessToken")
	defer span.End()

	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&AccessToken{}).Where("refresh_token = ?", refreshToken).Updates(&AccessToken{
		RevokedAt: &now,
	}).Error; err != nil {
		return fmt.Errorf("db.Updates: %w", err)
	}

	return nil
}
