package services

import (
	"context"
	"fmt"
	"photobox/image/proto"
)

type ImageService struct {
	proto.UnimplementedImageServiceServer

	RekognitionRepo
	LabelsRepo
}

func InitImageService(rekRepo RekognitionRepo, labelsRepo LabelsRepo) *ImageService {
	return &ImageService{RekognitionRepo: rekRepo, LabelsRepo: labelsRepo}
}

func (s *ImageService) DetectImageLabels(ctx context.Context, in *proto.DetectImageLabelsRequest) (*proto.DetectImageLabelsResponce, error) {
	// Detect labels with AWS Rekognition
	labels, err := s.RekognitionRepo.DetectImageLabelsWithS3(in.FileLocation)
	if err != nil {
		return &proto.DetectImageLabelsResponce{}, fmt.Errorf("service %v", err.Error())
	}

	for _, v := range labels {
		// Create labels
		label, err := s.LabelsRepo.CreateLabel(v)
		if err != nil {
			return &proto.DetectImageLabelsResponce{}, err
		}

		// Check if metadata label exist
		metaLabels, err := s.LabelsRepo.GetMetadataLabelByMetalId(in.MetaId)
		if err != nil {
			return &proto.DetectImageLabelsResponce{}, err
		}
		metaLabelExists := false
		for _, v := range metaLabels {
			if v.LabelID == label.ID {
				metaLabelExists = true
				break
			}
		}

		// Create only new metadata labels
		if !metaLabelExists {
			_, err := s.LabelsRepo.CreateMetadataLabel(in.MetaId, label.ID)
			if err != nil {
				return &proto.DetectImageLabelsResponce{}, err
			}
		}
	}

	return &proto.DetectImageLabelsResponce{}, nil
}
