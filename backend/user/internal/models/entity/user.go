package entity

import (
	"photobox-user/internal/models/response"
	"time"
)

type User struct {
	ID          string
	Email       string
	Username    string
	Password    string
	StorageUsed int64
	MaxStorage  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) ToResponse() response.User {
	return response.User{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		Password:    u.Password,
		StorageUsed: u.StorageUsed,
		MaxStorage:  u.MaxStorage,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
