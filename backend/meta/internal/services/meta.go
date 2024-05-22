package services

import (
	"context"
	"path/filepath"
	"photobox-meta/internal/models/entity"
	"photobox-meta/internal/utils"
	"photobox-meta/logger"
	"photobox-meta/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetaService struct {
	MetaRepo MetaRepo
	FileRepo FileRepo
	proto.UnimplementedMetaServiceServer

	userClient proto.UserServiceClient

	cfg    MetaServiceConfig
	logger logger.Logger
}

type MetaServiceConfig struct{}

func NewMetaService(metaRepo MetaRepo, fileRepo FileRepo, userClient proto.UserServiceClient, cfg MetaServiceConfig, logger logger.Logger) *MetaService {
	return &MetaService{MetaRepo: metaRepo, FileRepo: fileRepo, cfg: cfg, logger: logger, userClient: userClient}
}

func (s *MetaService) UploadMeta(ctx context.Context, in *proto.UplodaMetaRequest) (*proto.MetaResponse, error) {
	// Process metadata
	fileName := in.Filename
	fileSize := int64(len(in.File))
	fileExt := filepath.Ext(fileName)
	fileLocation := utils.GenerateS3FileLocation(fileName, fileExt)

	// Check storage availability
	_, err := s.userClient.UpdateStorageUsed(ctx, &proto.UpdateStorageUsedRequest{Id: in.UserId, FileSize: fileSize})
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	// Save file to AWS S3 Bucket
	err = s.FileRepo.UploadFile(fileLocation, in.File)
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	// Save metadata to Postgres
	meta, err := s.MetaRepo.CreateMeta(entity.Meta{
		UserID:           in.UserId,
		FileLocation:     fileLocation,
		FileName:         fileName,
		FileSize:         int(fileSize),
		FileExt:          fileExt,
		FileLastModified: in.FileLastModified.AsTime(),
	})
	if err != nil {
		return &proto.MetaResponse{}, nil
	}

	res := MakeMetaResponse(meta)
	return &res, nil
}

func (s *MetaService) GetMetaById(ctx context.Context, in *proto.GetMetaByIdRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) GetMetaByUser(ctx context.Context, in *proto.GetMetaByUserRequest) (*proto.GetMetaByUserResponse, error) {
	return &proto.GetMetaByUserResponse{}, nil
}

func (s *MetaService) GetAllMeta(ctx context.Context, in *proto.GetAllMetaRequest) (*proto.GetAllMetaResponse, error) {
	return &proto.GetAllMetaResponse{}, nil
}

func (s *MetaService) UpdateMeta(ctx context.Context, in *proto.UpdateMetaRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) DeleteMetaById(ctx context.Context, in *proto.DeleteMetaByIdRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) DeleteMetaByUser(ctx context.Context, in *proto.DeleteMetaByUserRequest) (*proto.DeleteMetaByUserResponse, error) {
	return &proto.DeleteMetaByUserResponse{}, nil
}

func MakeMetaResponse(meta entity.Meta) proto.MetaResponse {
	return proto.MetaResponse{
		Id:               meta.ID,
		UserId:           meta.UserID,
		FileLocation:     meta.FileLocation,
		FileName:         meta.FileName,
		FileSize:         int64(meta.FileSize),
		FileExt:          meta.FileExt,
		FileLastModified: timestamppb.New(meta.FileLastModified),
		CreatedAt:        timestamppb.New(meta.CreatedAt),
	}
}
