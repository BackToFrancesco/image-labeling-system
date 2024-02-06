package repositories

import (
	"context"
	"fabc.it/task-manager/datasources"
	"fabc.it/task-manager/domain"
	"fabc.it/task-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository struct {
	db *datasources.Database
}

func NewTaskRepository(
	db *datasources.Database,
) domain.TaskService {
	return &TaskRepository{db: db}
}

func (t TaskRepository) CreateNewTask(task *models.Task) error {
	result, err := t.db.Collection(datasources.TasksCollection).InsertOne(context.Background(), task)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		task.Id = oid.Hex()
	}

	return nil
}
