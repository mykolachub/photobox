package processors

import (
	"context"
	"photobox/image/internal/entities"

	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

const MIN_CONFIDENCE = 90

type RekognitionRepo struct {
	client *rekognition.Client
	cfg    RekognitionRepoConfig
}

type RekognitionRepoConfig struct {
	BuckerName string
}

func InitRekognition(client *rekognition.Client, cfg RekognitionRepoConfig) RekognitionRepo {
	return RekognitionRepo{client, cfg}
}

func (r RekognitionRepo) DetectImageLabelsWithS3(fileLocation string) ([]entities.Label, error) {
	labels := []entities.Label{}
	args := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: &r.cfg.BuckerName, Name: &fileLocation}}}

	detectedLabels, err := r.client.DetectLabels(context.TODO(), args)
	if err != nil {
		return []entities.Label{}, err
	}

	for _, v := range detectedLabels.Labels {
		if v.Confidence != nil && *v.Confidence > MIN_CONFIDENCE {
			label := entities.Label{Name: *v.Name}
			labels = append(labels, label)
		}
	}

	return labels, nil
}
