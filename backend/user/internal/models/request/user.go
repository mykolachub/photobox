package request

import (
	"photobox-user/internal/models/entity"
	"time"
)

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"usename"`
	Password    string    `json:"password"`
	StorageUsed int64     `json:"storage_used"`
	MaxStorage  int64     `json:"max_storage"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u User) ToEntity() entity.User {
	return entity.User{
		ID:          u.ID,
		Email:       u.Email,
		Password:    u.Password,
		Username:    u.Username,
		StorageUsed: u.StorageUsed,
		MaxStorage:  u.MaxStorage,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
