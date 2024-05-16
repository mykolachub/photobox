package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"photobox-user/internal/models/entity"
	"strings"
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
		users (username, email, password, created_at)
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
	user := entity.User{}

	query := `
	SELECT
		user_id, email, password,	
		username, storage_used, max_storage,	
		created_at, updated_at	
	FROM
		users	
	WHERE
		user_id = $1`

	err := r.db.QueryRow(query, id).Scan(
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
	users := []entity.User{}

	query := `
	SELECT
		user_id, email, password,	
		username, storage_used, max_storage,	
		created_at, updated_at
	FROM
		users`

	rows, err := r.db.Query(query)
	if err != nil {
		return []entity.User{}, err
	}
	for rows.Next() {
		user := entity.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Username,
			&user.StorageUsed,
			&user.MaxStorage,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return []entity.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) UpdateUser(id string, data entity.User) (entity.User, error) {
	user := entity.User{}

	updates := []string{}
	args := []interface{}{id}

	if data.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, data.Email)
	}
	if data.Password != "" {
		updates = append(updates, fmt.Sprintf("password = $%d", len(args)+1))
		args = append(args, data.Password)
	}
	if data.Username != "" {
		updates = append(updates, fmt.Sprintf("username = $%d", len(args)+1))
		args = append(args, data.Username)
	}

	if len(updates) == 0 {
		return entity.User{}, errors.New("empty update body")
	}

	query := `
	UPDATE
		users
	SET ` + strings.Join(updates, ", ") + ` 
	WHERE
		user_id = $1
	RETURNING
		user_id, email, password,	
		username, storage_used, max_storage,	
		created_at, updated_at
	`

	err := r.db.QueryRow(query, args...).Scan(
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

func (r *UserRepo) DeleteUser(id string) (entity.User, error) {
	user := entity.User{}

	query := `
	DELETE FROM
		users
	WHERE
		user_id = $1
	RETURNING
		*
	`

	err := r.db.QueryRow(query, id).Scan(
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
