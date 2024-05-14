package request

import "photobox-user/internal/models/entity"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"usename"`
	Password string `json:"password"`
}

func (u User) ToEntity() entity.User {
	return entity.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Username: u.Username,
	}
}
