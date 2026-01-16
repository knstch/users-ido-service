package users

import (
	"github.com/knstch/knstch-libs/log"
	"github.com/redis/go-redis/v9"

	"users-service/config"
	"users-service/internal/connector/google"
	"users-service/internal/users/repo"
)

type ServiceImpl struct {
	lg *log.Logger

	repo  repo.Repository
	redis *redis.Client

	google google.Client

	cfg config.Config
}

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
