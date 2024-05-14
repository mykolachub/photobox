package postgres

import (
	"database/sql"
	"photobox-user/internal/models/entity"
	"time"

	_ "github.com/lib/pq"
)

type UserRepo struct {
	db *sql.DB
}

func InitUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(data entity.User) (entity.User, error) {
	user := entity.User{}
	now := time.Now().UTC()

	query := `
	INSERT INTO
		users(username, email, password, created_at)
	VALUES
		($1, $2, $3, $4) 
	RETURNING
		user_id, email, password,
		username, storage_used, max_storage,
		created_at, updated_at`

	rows := r.db.QueryRow(query, data.Username, data.Email, data.Password, now)
	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Password,
		&user.StorageUsed,
		&user.MaxStorage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetUser(id string) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) GetUserByEmail(email string) (entity.User, error) {
	user := entity.User{}

	query := `
	SELECT
		user_id, email, password,	
		username, storage_used, max_storage,	
		created_at, updated_at	
	FROM
		users	
	WHERE
		email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.StorageUsed,
		&user.MaxStorage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
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
