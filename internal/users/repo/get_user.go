package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"gorm.io/gorm"

	"users-service/internal/domain/dto"
	"users-service/internal/users/filters"
	"users-service/internal/users/models"
)

func (r *DBRepo) GetUser(ctx context.Context, filters filters.UserFilter) (dto.User, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetUser")
	defer span.End()

	// Apply filter scope and return the first matching row.
	var user models.User
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
