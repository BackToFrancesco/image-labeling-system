package repositories

import (
	"context"
	"fabc.it/task-manager/datasources"
	"fabc.it/task-manager/domain"
	"fabc.it/task-manager/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (t TaskRepository) GetTask(taskId string) (task *models.Task, err error) {
	objectId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, err
	}

	result := t.db.Collection(datasources.TasksCollection).FindOne(context.Background(), bson.M{"_id": objectId})

	if result.Err() != nil {
		return nil, result.Err()
	}

	err = result.Decode(&task)
	if err != nil {
		return nil, err
	}

	return
}

func (t TaskRepository) UpdateTask(task *models.Task) error {
	objectId, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{{"subtasks", task.Subtasks}}}}

	_, err = t.db.Collection(datasources.TasksCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
