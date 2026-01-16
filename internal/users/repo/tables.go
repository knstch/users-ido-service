package repo

import "time"

type User struct {
	ID         uint64
	GoogleSub  string
	Email      string
	FirstName  string
	LastName   string
	ProfilePic string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

func (User) TableName() string {
	return "users"
}

type AccessToken struct {
	ID           uint64
	UserID       uint64
	AccessToken  string
	RefreshToken string
	RevokedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
