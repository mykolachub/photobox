package s3

import (
	"context"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3BucketConfig struct {
	AccessKey string
	Secret    string
	Region    string
	Name      string
	Endpoint  string
}

func InitS3Connection(s3Config S3BucketConfig) (*s3.Client, error) {
	appCreds := credentials.NewStaticCredentialsProvider(
		s3Config.AccessKey,
		s3Config.Secret,
		"",
	)

	// a := rekognition.NewFromConfig(*aws.NewConfig())
	// a.DetectLabels(context.TODO(), &rekognition.DetectLabelsInput{Image: &types.Image{Bytes: })
	cfg, err := awsCfg.LoadDefaultConfig(
		context.TODO(),
		awsCfg.WithRegion(s3Config.Region),
		awsCfg.WithCredentialsProvider(appCreds),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
