package repo

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"users-service/internal/metrics"
	"users-service/internal/users/models"
)

func (r *DBRepo) CreateUser(ctx context.Context, googleSub, email, firstName, lastName, profilePic string) (uint64, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: CreateUser")
	defer span.End()

	// Insert user record.
	user := &models.User{
		GoogleSub:  googleSub,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		ProfilePic: profilePic,
	}

	if err := r.db.WithContext(ctx).Model(&models.User{}).Create(user).Error; err != nil {
		return 0, fmt.Errorf("db.Create: %w", err)
	}

	// Metric: count newly created users.
	metrics.IncUsersCreated()
	return user.ID, nil
}
