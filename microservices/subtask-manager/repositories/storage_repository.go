package repositories

import (
	"archive/zip"
	"context"
	"errors"
	"fabc.it/subtask-manager/datasources"
	"fabc.it/subtask-manager/domain"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type StorageRepository struct {
	storage *datasources.Storage
}

func NewStorageRepository(storage *datasources.Storage) domain.StorageService {
	return &StorageRepository{storage: storage}
}

func (s StorageRepository) SaveImage(imageLabel string, image *zip.File) error {
	filePath := filepath.Join("", image.Name)

	sourceType, err := image.Open()
	if err != nil {
		return err
	}
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			log.Print(err)
		}
	}(sourceType)

	bytes := make([]byte, 512)

	_, err = sourceType.Read(bytes)
	if err != nil {
		return err
	}

	detectedContent := http.DetectContentType(bytes)

	if detectedContent != "image/jpeg" && detectedContent != "image/png" {
		return errors.New("file is not an image")
	}

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

	destination, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, image.Mode())
	if err != nil {
		return err
	}

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	_, err = s.storage.FPutObject(context.Background(), "images", imageLabel, filePath, minio.PutObjectOptions{
		ContentType: detectedContent,
	})

	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			log.Print(err)
		}
	}()

	if err != nil {
		return err
	}

	return nil
}
