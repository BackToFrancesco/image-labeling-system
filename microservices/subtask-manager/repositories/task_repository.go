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

func (t *TaskRepository) UpdateSubtaskLabels(labelSubtask *models.LabelSubtask) (err error) {
	for _, pair := range labelSubtask.AssignedLabels {
		// Build the field increment expression for each label
		field := "assignedLabels." + pair.Label
		update := bson.M{"$inc": bson.M{field: 1}}

		// Build the filter to find the document with the matching ImageId
		filter := bson.M{"_id": pair.ImageId}

		// TODO: if label not exist what to do?
		// TODO: control if the user that is going to assign a label is in assignee?
		// TODO: handle concurrency?
		_, err := t.db.Collection(datasources.TasksCollection).UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}
