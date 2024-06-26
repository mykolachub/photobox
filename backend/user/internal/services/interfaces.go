package services

import "photobox-user/internal/models/entity"

type Storages struct {
	UserRepo UserRepo
}

type UserRepo interface {
	CreateUser(data entity.User) (entity.User, error)
	GetUser(id string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
	UpdateUser(id string, data entity.User) (entity.User, error)
	UpdateStorage(id string, file_size int64) (entity.User, error)
	DeleteUser(id string) (entity.User, error)
}
