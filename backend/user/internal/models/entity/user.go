package entity

import "photobox-user/internal/models/response"

type User struct {
	ID       string
	Email    string
	Username string
	Password string
}

func (u *User) ToResponse() response.User {
	return response.User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Password: u.Password,
	}
}
