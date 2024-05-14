package response

import "time"

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
