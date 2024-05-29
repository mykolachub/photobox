package services

import "photobox/image/internal/entities"

type Storages struct {
	LabesRepo LabelsRepo
}

type LabelsRepo interface {
	CreateLabel(data entities.Label) (entities.Label, error)
	CreateMetadataLabel(metadata_id, label_id string) (entities.MetadataLabel, error)
	GetMetadataLabelByMetalId(metadata_id string) ([]entities.MetadataLabel, error)
}

type Processors struct {
	RekognitionRepo RekognitionRepo
}

type RekognitionRepo interface {
	DetectImageLabelsWithS3(fileLocation string) ([]entities.Label, error)
}
