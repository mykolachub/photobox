package response

import "time"

type User struct {
	ID          string    `json:"id"`
	GoogleID    string    `json:"google_id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Username    string    `json:"usename"`
	Picture     string    `json:"picture"`
	StorageUsed int64     `json:"storage_used"`
	MaxStorage  int64     `json:"max_storage"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
