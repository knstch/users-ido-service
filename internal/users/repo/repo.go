package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/knstch/knstch-libs/log"
	"gorm.io/gorm"

	"users-service/internal/domain/dto"
)

type DBRepo struct {
	lg *log.Logger
	db *gorm.DB
}

type Repository interface {
	Transaction(fn func(st Repository) error) error
	CreateUser(ctx context.Context, googleSub, email, firstName, lastName, profilePic string) (uint64, error)
	CreateAccessTokens(ctx context.Context, accessToken, refreshToken string, userID uint64) error
	GetAccessTokens(ctx context.Context, filter AccessTokenFilter) (dto.AccessTokens, error)
	GetUser(ctx context.Context, filters UserFilter) (dto.User, error)
	UpdateUserMetadata(ctx context.Context, id uint64, firstName, lastName, profilePic string) error
}

func (r *DBRepo) NewDBRepo(db *gorm.DB) *DBRepo {
	if db == nil {
		db = r.db.Session(&gorm.Session{NewDB: true})
	}
	return &DBRepo{
		db: db,
		lg: r.lg,
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *DBRepo) Transaction(fn func(st Repository) error) error {
	db := r.db.Session(&gorm.Session{NewDB: true})
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := fn(r.NewDBRepo(tx)); err != nil {
			return fmt.Errorf("fn: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("db.Transaction: %w", err)
	}
	return nil
}
