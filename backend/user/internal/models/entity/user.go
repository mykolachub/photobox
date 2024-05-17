package entity

import (
	"photobox-user/internal/models/response"
	"time"
)

type User struct {
	ID          string
	GoogleID    string
	Email       string
	Password    string
	Username    string
	Picture     string
	StorageUsed int64
	MaxStorage  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) ToResponse() response.User {
	return response.User{
		ID:          u.ID,
		GoogleID:    u.GoogleID,
		Email:       u.Email,
		Password:    u.Password,
		Username:    u.Username,
		Picture:     u.Picture,
		StorageUsed: u.StorageUsed,
		MaxStorage:  u.MaxStorage,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
