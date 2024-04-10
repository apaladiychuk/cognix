package storage

import (
	"context"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinioConfig struct {
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID,required"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY,required"`
	Endpoint        string `env:"MINIO_ENDPOINT,required"`
	UseSSL          bool   `env:"MINIO_USE_SSL"`
	BucketName      string `env:"MINIO_BUCKET_NAME,required"`
	Region          string `env:"MINIO_REGION,required"`
}
type MinIOClient interface {
	Upload(ctx context.Context, filename string, reader io.Reader) (string, error)
	GetObject(ctx context.Context, filename string, writer io.Writer) error
}
type minIOClient struct {
	BucketName string
	Region     string
	client     *minio.Client
}

func (c *minIOClient) Upload(ctx context.Context, filename string, reader io.Reader) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *minIOClient) GetObject(ctx context.Context, filename string, writer io.Writer) error {
	//TODO implement me
	panic("implement me")
}

func NewMinIOClient(cfg *MinioConfig) (MinIOClient, error) {

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKeyID,
			cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &minIOClient{
		BucketName: cfg.BucketName,
		Region:     cfg.Region,
		client:     minioClient,
	}, nil
}
