package postgres

import (
	"database/sql"
	"photobox-user/internal/models/entity"

	_ "github.com/lib/pq"
)

type UserRepo struct {
	db *sql.DB
}

func InitUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(data entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) GetUser(id string) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) GetAllUsers() ([]entity.User, error) {
	return []entity.User{}, nil
}

func (r *UserRepo) UpdateUser(id string, data entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) DeleteUser(id string) (entity.User, error) {
	return entity.User{}, nil
}
