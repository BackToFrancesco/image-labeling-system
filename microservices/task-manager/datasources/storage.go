package datasources

import (
	"fabc.it/task-manager/config"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type Storage struct {
	*minio.Client
}

func NewStorage(env *config.Env) *Storage {
	uri := fmt.Sprintf("%s:%s", env.MinioHost, env.MinioPort)

	client, err := minio.New(uri, &minio.Options{
		Creds: credentials.NewStaticV2(env.MinioUsername, env.MinioPassword, ""),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{client}
}
