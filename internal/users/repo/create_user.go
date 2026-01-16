package repo

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"
)

func (r *DBRepo) CreateUser(ctx context.Context, googleSub, email, firstName, lastName, profilePic string) (uint64, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetPassword")
	defer span.End()

	user := &User{
		GoogleSub:  googleSub,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		ProfilePic: profilePic,
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Create(user).Error; err != nil {
		return 0, fmt.Errorf("db.Create: %w", err)
	}

	return user.ID, nil
}
