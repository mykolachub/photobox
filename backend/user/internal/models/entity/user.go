package entity

import (
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
