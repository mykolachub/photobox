package services

import (
	"context"
	"photobox-meta/internal/entity"
)

type Storages struct {
	MetaRepo MetaRepo
	FileRepo FileRepo
}

type MetaRepo interface {
	CreateMeta(data entity.Meta) (entity.Meta, error)
	GetMetaById(id string) (entity.Meta, error)
	GetMetaByFileLocation(fileLocation string) (entity.Meta, error)
	GetAllMeta() ([]entity.Meta, error)
	GetAllMetaByUserId(user_id string) ([]entity.Meta, error)
	UpdateMeta(id string, data entity.Meta) (entity.Meta, error)
	DeleteMeta(id string) (entity.Meta, error)
	DeleteMetaByUser(user_id string) ([]entity.Meta, error)
}

type FileRepo interface {
	GetFile(filePath string) ([]byte, error)
	UploadFile(filePath string, file []byte) error
	DeleteFile(filePath string) error
	DeleteFiles(objectKeys []string) error
}

type MQ interface {
	Publish(ctx context.Context, name string, body []byte) error
}
