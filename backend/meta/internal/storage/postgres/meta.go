package postgres

import (
	"database/sql"
	"photobox-meta/internal/models/entity"

	_ "github.com/lib/pq"
)

type MetaRepo struct {
	db *sql.DB
}

func InitMetaRepo(db *sql.DB) *MetaRepo {
	return &MetaRepo{db: db}
}

func (r *MetaRepo) CreateMeta(data entity.Meta) (entity.Meta, error) {
	return entity.Meta{}, nil
}

func (r *MetaRepo) GetMeta(id string) (entity.Meta, error) {
	return entity.Meta{}, nil
}

func (r *MetaRepo) GetMetaByUser(user_id string) ([]entity.Meta, error) {
	return []entity.Meta{}, nil
}

func (r *MetaRepo) GetAllMeta() ([]entity.Meta, error) {
	return []entity.Meta{}, nil
}

func (r *MetaRepo) UpdateMeta(id string, data entity.Meta) (entity.Meta, error) {
	return entity.Meta{}, nil
}

func (r *MetaRepo) DeleteMeta(id string) (entity.Meta, error) {
	return entity.Meta{}, nil
}

func (r *MetaRepo) DeleteMetaByUser(user_id string) ([]entity.Meta, error) {
	return []entity.Meta{}, nil
}
