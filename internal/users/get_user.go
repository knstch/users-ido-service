package users

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"users-service/internal/domain/dto"
	"users-service/internal/users/filters"
)

func (s *ServiceImpl) GetUser(ctx context.Context, userToFind dto.GetUser) (dto.User, error) {
	ctx, span := tracing.StartSpan(ctx, "users: GetUser")
	defer span.End()

	user, err := s.repo.GetUser(ctx, filters.UserFilter{
		ID:        userToFind.ID,
		GoogleSub: userToFind.GoogleSub,
		Email:     userToFind.Email,
		FirstName: userToFind.FirstName,
		LastName:  userToFind.LastName,
	})
	if err != nil {
		return dto.User{}, fmt.Errorf("repo.GetUser: %w", err)
	}

	return user, nil
}
