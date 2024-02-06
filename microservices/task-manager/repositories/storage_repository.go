package repositories

import (
	"archive/zip"
	"context"
	"errors"
	"fabc.it/task-manager/datasources"
	"fabc.it/task-manager/domain"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
	"os"
)

type StorageRepository struct {
	storage *datasources.Storage
}

func NewStorageRepository(storage *datasources.Storage) domain.StorageService {
	return &StorageRepository{storage: storage}
}

func (s StorageRepository) SaveImage(imageLabel string, image zip.File) error {
	source, err := image.Open()
	if err != nil {
		return err
	}
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			log.Print(err)
		}
	}(source)

	bytes := make([]byte, 512)

	_, err = source.Read(bytes)
	if err != nil {
		return err
	}

	detectedContent := http.DetectContentType(bytes)

	if detectedContent != "image/jpeg" && detectedContent != "image/png" {
		return errors.New("file is not an image")
	}

	destination, err := os.OpenFile(image.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, image.Mode())
	if err != nil {
		return err
	}

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	_, err = s.storage.FPutObject(context.Background(), "images", imageLabel, image.Name, minio.PutObjectOptions{
		ContentType: detectedContent,
	})

	return err
}
