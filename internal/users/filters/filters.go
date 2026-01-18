package filters

import (
	"gorm.io/gorm"

	"users-service/internal/users/models"
)

type AccessTokenFilter struct {
	AccessToken  string
	RefreshToken string
	UserID       uint64
}

func (f *AccessTokenFilter) ToScope() func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&models.AccessToken{})

		if f.AccessToken != "" {
			tx = tx.Where("access_token = ?", f.AccessToken)
		}
		if f.RefreshToken != "" {
			tx = tx.Where("refresh_token = ?", f.RefreshToken)
		}
		if f.UserID != 0 {
			tx = tx.Where("user_id = ?", f.UserID)
		}

		return tx
	}
}

type UserFilter struct {
	ID        uint64
	GoogleSub string
	Email     string
	FirstName string
	LastName  string
}

func (f *UserFilter) ToScope() func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&models.User{})

		if f.ID != 0 {
			tx = tx.Where("id = ?", f.ID)
		}

		if f.GoogleSub != "" {
			tx = tx.Where("google_sub = ?", f.GoogleSub)
		}

		if f.Email != "" {
			tx = tx.Where("email = ?", f.Email)
		}

		if f.FirstName != "" {
			tx = tx.Where("first_name = ?", f.FirstName)
		}

		if f.LastName != "" {
			tx = tx.Where("last_name = ?", f.LastName)
		}

		return tx
	}
}
