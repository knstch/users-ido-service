package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/knstch-libs/auth"
	"github.com/knstch/knstch-libs/tracing"
	public "github.com/knstch/users-ido-api/public"

	"users-service/internal/domain/dto"
)

func MakeGetUserEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.GetUser(ctx, request.(*public.GetUserRequest))
	}
}

func (c *Controller) GetUser(ctx context.Context, req *public.GetUserRequest) (*public.GetUserResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: RefreshAccessToken")
	defer span.End()

	userFromCookie, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	if req.GetID() != 0 && req.GetID() != uint64(userFromCookie.UserID) {
		return nil, ErrAccessDenied
	}

	user, err := c.svc.GetUser(ctx, dto.GetUser{
		ID:    req.GetID(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, fmt.Errorf("svc.GetUser: %w", err)
	}

	if user.ID != uint64(userFromCookie.UserID) {
		return nil, ErrAccessDenied
	}

	return &public.GetUserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfilePicture: user.ProfilePicture,
	}, nil
}
