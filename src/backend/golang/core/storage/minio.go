package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"io"
)

type (
	MinioConfig struct {
		AccessKey       string `env:"MINIO_ACCESS_KEY"`
		SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
		Endpoint        string `env:"MINIO_ENDPOINT"`
		UseSSL          bool   `env:"MINIO_USE_SSL"`
		BucketName      string `env:"MINIO_BUCKET_NAME"`
		Region          string `env:"MINIO_REGION"`
		Mocked          bool   `env:"MINIO_MOCKED" envDefault:"true"`
	}
	MinIOClient interface {
		Upload(ctx context.Context, filename, contentType string, reader io.Reader) (string, string, error)
		GetObject(ctx context.Context, filename string, writer io.Writer) error
	}
	minIOClient struct {
		BucketName string
		Region     string
		client     *minio.Client
	}
	minIOMockClient struct{}
)

func (c *minIOClient) Upload(ctx context.Context, filename, contentType string, reader io.Reader) (string, string, error) {
	objectName := fmt.Sprintf("%s-%s", filename, uuid.New().String())
	client := *c.client

	res, err := client.PutObject(ctx, c.BucketName, objectName, reader, -1,
		minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		return "", "", utils.Internal.Wrapf(err, "cannot upload file: %s", err.Error())
	}
	return res.Key, res.ChecksumCRC32C, nil
}

func (c *minIOClient) GetObject(ctx context.Context, filename string, writer io.Writer) error {
	object, err := c.client.GetObject(ctx, c.BucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()
	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}
	return nil
}

func NewMinIOClient(cfg *MinioConfig) (MinIOClient, error) {

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKey,
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

func NewMinIOMockClient() (MinIOClient, error) {
	zap.S().Info("Run with mocked minio client")
	return &minIOMockClient{}, nil
}

func (m minIOMockClient) Upload(ctx context.Context, filename, contentType string, reader io.Reader) (string, string, error) {
	return fmt.Sprintf("bucket/%s", filename), "sign", nil
}

func (m minIOMockClient) GetObject(ctx context.Context, filename string, writer io.Writer) error {
	return nil
}
