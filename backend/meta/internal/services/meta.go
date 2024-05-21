package services

import (
	"context"
	"photobox-meta/proto"
)

type MetaService struct {
	MetaRepo MetaRepo
	proto.UnimplementedMetaServiceServer

	cfg MetaServiceConfig
}

type MetaServiceConfig struct{}

func NewMetaService(r MetaRepo, cfg MetaServiceConfig) *MetaService {
	return &MetaService{MetaRepo: r, cfg: cfg}
}

func (s *MetaService) DeleteMetaById(ctx context.Context, in *proto.DeleteMetaByIdRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) DeleteMetaByUser(ctx context.Context, in *proto.DeleteMetaByUserRequest) (*proto.DeleteMetaByUserResponse, error) {
	return &proto.DeleteMetaByUserResponse{}, nil
}

func (s *MetaService) GetAllMeta(ctx context.Context, in *proto.GetAllMetaRequest) (*proto.GetAllMetaResponse, error) {
	return &proto.GetAllMetaResponse{}, nil
}

func (s *MetaService) GetMetaById(ctx context.Context, in *proto.GetMetaByIdRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) GetMetaByUser(ctx context.Context, in *proto.GetMetaByUserRequest) (*proto.GetMetaByUserResponse, error) {
	return &proto.GetMetaByUserResponse{}, nil
}

func (s *MetaService) UpdateMeta(ctx context.Context, in *proto.UpdateMetaRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}

func (s *MetaService) UploadMeta(ctx context.Context, in *proto.UplodaMetaRequest) (*proto.MetaResponse, error) {
	return &proto.MetaResponse{}, nil
}
