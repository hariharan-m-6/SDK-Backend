package service

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client *s3.Client
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &S3Service{
		client: s3.NewFromConfig(cfg),
	}, nil
}

func (s *S3Service) GetFile(ctx context.Context, key string) ([]byte, error) {
	bucket := os.Getenv("AWS_BUCKET_NAME")

	base64Key := os.Getenv("S3_ENCRYPTION_KEY")

	decodedKey, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	hash := md5.Sum(decodedKey)
	keyMD5 := base64.StdEncoding.EncodeToString(hash[:])

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),

		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String(base64.StdEncoding.EncodeToString(decodedKey)),
		SSECustomerKeyMD5:    aws.String(keyMD5),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}