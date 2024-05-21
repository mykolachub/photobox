package services

import "photobox-meta/internal/models/entity"

type Storages struct {
	MetaRepo MetaRepo
}

type MetaRepo interface {
	CreateMeta(data entity.Meta) (entity.Meta, error)
	GetMeta(id string) (entity.Meta, error)
	GetMetaByUser(user_id string) ([]entity.Meta, error)
	GetAllMeta() ([]entity.Meta, error)
	UpdateMeta(id string, data entity.Meta) (entity.Meta, error)
	DeleteMeta(id string) (entity.Meta, error)
	DeleteMetaByUser(user_id string) ([]entity.Meta, error)
}
