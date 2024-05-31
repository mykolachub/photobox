package services

import (
	"context"
	"errors"
	"path/filepath"
	"photobox-meta/internal/entity"
	"photobox-meta/internal/utils"
	"photobox-meta/logger"
	"photobox-meta/proto"

	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetaService struct {
	MetaRepo MetaRepo
	FileRepo FileRepo
	MQ       MQ
	proto.UnimplementedMetaServiceServer

	userClient proto.UserServiceClient

	cfg    MetaServiceConfig
	logger logger.Logger
}

type MetaServiceConfig struct{}

func NewMetaService(metaRepo MetaRepo, fileRepo FileRepo, userClient proto.UserServiceClient, cfg MetaServiceConfig, logger logger.Logger, mq MQ) *MetaService {
	return &MetaService{MetaRepo: metaRepo, FileRepo: fileRepo, cfg: cfg, logger: logger, userClient: userClient, MQ: mq}
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
		FileWidth:        int(in.FileWidth),
		FileHeight:       int(in.FileHeight),
	})
	if err != nil {
		return &proto.MetaResponse{}, nil
	}

	// Produsing image_detect message to image microservice for image labels detection
	protoReqBody := &proto.DetectImageLabelsRequest{
		FileLocation: meta.FileLocation,
		MetaId:       meta.ID,
	}

	bytesReqBody, err := pb.Marshal(protoReqBody)
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	err = s.MQ.Publish(ctx, "image_detect", bytesReqBody)
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	res := MakeMetaResponse(meta)
	return &res, nil
}

func (s *MetaService) GetMetaById(ctx context.Context, in *proto.GetMetaByIdRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) GetMetaByUser(ctx context.Context, in *proto.GetMetaByUserRequest) (*proto.GetMetaByUserResponse, error) {
	metas, err := s.MetaRepo.GetAllMetaByUserId(in.UserId)
	if err != nil {
		return &proto.GetMetaByUserResponse{}, err
	}

	res := proto.GetMetaByUserResponse{}
	for _, v := range metas {
		meta := MakeMetaResponse(v)
		res.Metas = append(res.Metas, &meta)
	}

	return &res, nil
}

func (s *MetaService) GetAllMeta(ctx context.Context, in *proto.GetAllMetaRequest) (*proto.GetAllMetaResponse, error) {
	return &proto.GetAllMetaResponse{}, nil
}

func (s *MetaService) UpdateMeta(ctx context.Context, in *proto.UpdateMetaRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) DeleteMetaById(ctx context.Context, in *proto.DeleteMetaByIdRequest) (*proto.MetaResponse, error) {
	meta, err := s.MetaRepo.GetMetaById(in.Id)
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	err = s.FileRepo.DeleteFile(meta.FileLocation)
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	meta, err = s.MetaRepo.DeleteMeta(in.Id)
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	_, err = s.userClient.UpdateStorageUsed(ctx, &proto.UpdateStorageUsedRequest{Id: meta.UserID, FileSize: int64(meta.FileSize) * -1})
	if err != nil {
		return &proto.MetaResponse{}, err
	}

	res := MakeMetaResponse(meta)
	return &res, nil
}

func (s *MetaService) DeleteMetaByUser(ctx context.Context, in *proto.DeleteMetaByUserRequest) (*proto.DeleteMetaByUserResponse, error) {
	return &proto.DeleteMetaByUserResponse{}, nil
}

func (s *MetaService) GetFileByKey(ctx context.Context, in *proto.GetFileByKeyRequest) (*proto.FileResponse, error) {
	// Check if user has this file
	meta, err := s.MetaRepo.GetMetaByFileLocation(in.Key)
	if err != nil {
		return &proto.FileResponse{}, err
	}

	if meta.UserID != in.UserId {
		return &proto.FileResponse{}, errors.New("file does not belong to user")
	}

	// Get file from S3
	file, err := s.FileRepo.GetFile(in.Key)
	if err != nil {
		return &proto.FileResponse{}, err
	}

	return &proto.FileResponse{File: file}, nil
}

func MakeMetaResponse(meta entity.Meta) proto.MetaResponse {
	res := proto.MetaResponse{
		Id:               meta.ID,
		UserId:           meta.UserID,
		FileLocation:     meta.FileLocation,
		FileName:         meta.FileName,
		FileSize:         int64(meta.FileSize),
		FileExt:          meta.FileExt,
		FileLastModified: timestamppb.New(meta.FileLastModified),
		CreatedAt:        timestamppb.New(meta.CreatedAt),
		FileWidth:        int32(meta.FileWidth),
		FileHeight:       int32(meta.FileHeight),
	}

	labels := []*proto.Label{}
	for _, v := range meta.Labels {
		label := proto.Label{Id: v.ID, Value: v.Value, MetadataLabelId: v.MetadataLabelID}
		labels = append(labels, &label)
	}
	res.Labels = labels

	return res
}
