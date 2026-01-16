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

func (r *DBRepo) GetUser(ctx context.Context, filters UserFilter) (dto.User, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetUser")
	defer span.End()

	var user User
	if err := r.db.WithContext(ctx).Scopes(filters.ToScope()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.User{}, fmt.Errorf("db.First: %w", svcerrs.ErrDataNotFound)
		}
		return dto.User{}, fmt.Errorf("db.First: %w", err)
	}

	return dto.User{
		ID:             user.ID,
		GoogleSub:      user.GoogleSub,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfilePicture: user.ProfilePic,
	}, nil
}
