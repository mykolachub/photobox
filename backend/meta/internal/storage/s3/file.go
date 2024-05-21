package s3

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type FileRepo struct {
	client *s3.Client
	cfg    FileRepoConfig
}

type FileRepoConfig struct {
	BucketName string
}

func InitFileRepo(client *s3.Client, cfg FileRepoConfig) FileRepo {
	return FileRepo{client: client, cfg: cfg}
}

func (r FileRepo) GetFile(filePath string) ([]byte, error) {
	obj, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.cfg.BucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(obj.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r FileRepo) UploadFile(filePath string, file []byte) error {
	contentLen := int64(len(file))
	contentType := http.DetectContentType(file)

	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(r.cfg.BucketName),
		Key:           aws.String(filePath),
		Body:          bytes.NewReader(file),
		ContentLength: &contentLen,
		ContentType:   &contentType,
	})
	return err
}

func (r FileRepo) DeleteFile(filePath string) error {
	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.cfg.BucketName),
		Key:    aws.String(filePath),
	})
	return err
}

func (r FileRepo) DeleteFiles(objectKeys []string) error {
	objectIdentifiers := make([]types.ObjectIdentifier, 0, len(objectKeys))
	for _, key := range objectKeys {
		objectIdentifiers = append(objectIdentifiers, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	_, err := r.client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(r.cfg.BucketName),
		Delete: &types.Delete{
			Objects: objectIdentifiers,
		},
	})
	return err
}
