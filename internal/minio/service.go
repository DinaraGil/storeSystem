package minio

import (
	"bytes"
	"context"
	"fmt"
	"storeSystem/internal/config"
	"storeSystem/internal/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (m *minioClient) CreateOne(file helpers.FileDataType) (string, error) {
	objectID := uuid.New().String() + "_" + file.FileName

	reader := bytes.NewReader(file.Data)

	_, err := m.mc.PutObject(
		context.Background(),
		config.AppConfig.BucketName,
		objectID,
		reader,
		int64(len(file.Data)),
		minio.PutObjectOptions{
			ContentType: "text/csv",
		},
	)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании объекта %s: %v", file.FileName, err)
	}

	url, err := m.mc.PresignedGetObject(
		context.Background(),
		config.AppConfig.BucketName,
		objectID,
		time.Hour*24,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании URL для объекта %s: %v", file.FileName, err)
	}

	return url.String(), nil
}

func (m *minioClient) GetOne(objectID string) (string, error) {
	// Получение предварительно подписанного URL для доступа к объекту Minio.
	url, err := m.mc.PresignedGetObject(context.Background(), config.AppConfig.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении URL для объекта %s: %v", objectID, err)
	}

	return url.String(), nil
}
