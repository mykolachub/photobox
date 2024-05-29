package processors

import (
	"context"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
)

type RekognitionConfig struct {
	AccessKey string
	Secret    string
	Region    string
}

func NewRekognitionConnection(rekognitionConfig RekognitionConfig) (*rekognition.Client, error) {
	cfg, err := awsCfg.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	cfg.Region = rekognitionConfig.Region

	return rekognition.NewFromConfig(cfg), nil
}
