package repositories

import (
	"context"

	"fabc.it/subtask-manager/datasources"
	"fabc.it/subtask-manager/domain"
	"fabc.it/subtask-manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
/*
func (t TaskRepository) UpdateSubtask(subtask *models.Subtask) error {
	taskId := strings.Split(subtask.Id, "-")[0]

	objectId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{"_id", objectId},
		{"subtasks", bson.D{{"$elemMatch", bson.D{{"id", subtask.Id}}}}},
	}
	update := bson.D{{"$set", bson.D{{"subtasks.$.label", subtask.Label}}}}

	_, err = t.db.Collection(datasources.TasksCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
*/
func (t TaskRepository) AddUserIDToAssignees(subtasks []models.Subtask, userID string) error {
    for _, subtask := range subtasks {
        filter := bson.M{"_id": subtask.Id}
        update := bson.M{"$addToSet": bson.M{"assignee": userID}} // this will add the userID only if it's not already present
        _, err := t.db.Collection(datasources.TasksCollection).UpdateOne(context.Background(), filter, update)
        if err != nil {
            return err
        }
    }
    return nil
}


func (t TaskRepository) GetSubtasks(numberOfSubtasks int, userId string) (task *[]models.Subtask, err error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"assignee", bson.D{{"$ne", userId}}}}}}, // Exclude the user ID from assignee list
		{{"$addFields", bson.D{{"numAssignees", bson.D{{"$size", "$assignee"}}}}}},
		{{"$sort", bson.D{{"numAssignees", 1}}}}, // Sort by the count of assignees in ascending order
		{{"$limit", numberOfSubtasks}},
	}

	// Aggregate using the pipeline
	cursor, err := t.db.Collection(datasources.TasksCollection).Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []models.Subtask

	// Iterate through the cursor and decode each document
	for cursor.Next(context.Background()) {
		var task models.Subtask
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		results = append(results, task)
	}

	// Check for errors that may have occurred during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// TODO: handle concurrency
	t.AddUserIDToAssignees(results, userId)

	return &results, nil
}
