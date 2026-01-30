package users

import (
	"context"

	"github.com/knstch/knstch-libs/log"
	"github.com/redis/go-redis/v9"

	"users-service/config"
	"users-service/internal/connector/google"
	"users-service/internal/domain/dto"
	"users-service/internal/users/repo"
)

type ServiceImpl struct {
	lg *log.Logger

	repo  repo.Repository
	redis *redis.Client

	google google.Client

	cfg config.Config
}

type Service interface {
	// AuthViaGoogle returns a Google login URL.
	// stateURL is the original page URL/path to return to after successful login.
	AuthViaGoogle(ctx context.Context, stateURL string, scheme string) (string, error)
	// CompleteLogin completes the OAuth flow by exchanging `code` for tokens and
	// issuing service JWTs. It returns tokens, the validated return URL/path, and the original request scheme.
	CompleteLogin(ctx context.Context, state, code string) (dto.AccessTokens, string, string, error)
	// GetUser returns a user by filter fields.
	GetUser(ctx context.Context, userToFind dto.GetUser) (dto.User, error)
	// RefreshAccessToken revokes an old tokens pair and makes a new one.
	RefreshAccessToken(ctx context.Context, refreshToken string) (dto.AccessTokens, error)
}

// NewService constructs the Users service.
func NewService(
	lg *log.Logger,
	repo repo.Repository,
	cfg config.Config,
	googleClient google.Client,
	redis *redis.Client,
) *ServiceImpl {
	return &ServiceImpl{
		lg:     lg,
		repo:   repo,
		cfg:    cfg,
		google: googleClient,
		redis:  redis,
	}
}
