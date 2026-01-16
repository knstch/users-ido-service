package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"gorm.io/gorm"

	"users-service/internal/users/modles"
)

func (r *DBRepo) UpdateUserMetadata(ctx context.Context, id uint64, firstName, lastName, profilePic string) error {
	ctx, span := tracing.StartSpan(ctx, "repo: UpdateUserMetadata")
	defer span.End()

	var user modles.User
	if err := r.db.WithContext(ctx).Model(&modles.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("db.First: %w", svcerrs.ErrDataNotFound)
		}
		return fmt.Errorf("db.First: %w", err)
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfilePic = profilePic

	if err := r.db.Save(&user).Error; err != nil {
		return fmt.Errorf("db.Save: %w", err)
	}

	return nil
}
