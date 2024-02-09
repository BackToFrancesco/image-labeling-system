package datasources

import (
	"context"
	"errors"
	"fabc.it/task-manager/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Database struct {
	*mongo.Database
}

const (
	TasksCollection = "tasks"
)

func NewDatabase(env *config.Env) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s", env.DBHost, env.DBPort)

	credentials := options.Credential{
		Username: env.DBUsername,
		Password: env.DBPassword,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credentials))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(env.DBName)

	err = db.CreateCollection(ctx, TasksCollection)
	if err != nil && !errors.As(err, &mongo.CommandError{}) {
		log.Fatal(err)
	}

	return &Database{db}
}
