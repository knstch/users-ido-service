package modles

import "time"

type User struct {
	ID         uint64    `gorm:"primaryKey;column:id"`
	GoogleSub  string    `gorm:"column:google_sub"`
	Email      string    `gorm:"column:email"`
	FirstName  string    `gorm:"column:first_name"`
	LastName   string    `gorm:"column:last_name"`
	ProfilePic string    `gorm:"column:profile_picture"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	DeletedAt  time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

type AccessToken struct {
	ID           uint64     `gorm:"primaryKey;column:id"`
	UserID       uint64     `gorm:"column:user_id"`
	AccessToken  string     `gorm:"column:access_token"`
	RefreshToken string     `gorm:"column:refresh_token"`
	RevokedAt    *time.Time `gorm:"column:revoked_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
