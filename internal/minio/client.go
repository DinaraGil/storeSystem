package minio

import (
	"context"
	"storeSystem/internal/config"
	"storeSystem/internal/helpers"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client интерфейс для взаимодействия с Minio
type Client interface {
	InitMinio() error
	CreateOne(file helpers.FileDataType) (*UploadedFile, error)
	GetOne(objectID string) (string, error)
}

// minioClient реализация интерфейса MinioClient
type minioClient struct {
	mc *minio.Client
}

func NewMinioClient() Client {
	return &minioClient{}
}

func (m *minioClient) InitMinio() error {
	ctx := context.Background()

	client, err := minio.New(config.AppConfig.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.MinioRootUser, config.AppConfig.MinioRootPassword, ""),
		Secure: config.AppConfig.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	m.mc = client

	exists, err := m.mc.BucketExists(ctx, config.AppConfig.BucketName)
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, config.AppConfig.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
