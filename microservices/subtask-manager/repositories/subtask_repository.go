package repositories

import (
	"context"
	"fabc.it/subtask-manager/datasources"
	"fabc.it/subtask-manager/domain"
	"fabc.it/subtask-manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	db *datasources.Database
}

func NewSubtaskRepository(
	db *datasources.Database,
) domain.SubtaskService {
	return &TaskRepository{db: db}
}

func (t TaskRepository) AddUserIDToAssignees(subtasks []*models.SubtaskMessage, userID string) error {
	for _, subtask := range subtasks {
		filter := bson.M{"_id": subtask.Id}
		update := bson.M{"$addToSet": bson.M{"assignee": userID}}

		_, err := t.db.Collection(datasources.TasksCollection).UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t TaskRepository) GetSubtasks(numberOfSubtasks int, userId string) (task []*models.SubtaskMessage, err error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "assignee", Value: bson.D{{Key: "$ne", Value: userId}}}}}}, // Exclude the user ID from assignee list
		{{Key: "$addFields", Value: bson.D{{Key: "numAssigneValue: es", Value: bson.D{{Key: "$size", Value: "$assignee"}}}}}},
		{{Key: "$sort", Value: bson.D{{Key: "numAssignees", Value: 1}}}}, // Sort by the count of assignees in ascending order
		{{Key: "$limit", Value: numberOfSubtasks}},
	}

	// Aggregate using the pipeline
	cursor, err := t.db.Collection(datasources.TasksCollection).Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*models.SubtaskMessage

	// Iterate through the cursor and decode each document
	for cursor.Next(context.Background()) {
		var subtask models.SubtaskMessage
		err := cursor.Decode(&subtask)
		if err != nil {
			return nil, err
		}
		results = append(results, &subtask)
	}

	// Check for errors that may have occurred during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	t.AddUserIDToAssignees(results, userId)

	return results, nil
}

func (t *TaskRepository) UpdateSubtaskLabel(labelSubtask *models.LabelSubtask) (subtask *models.Subtask, err error) {
	field := "assignedLabels." + labelSubtask.AssignedLabel
	update := bson.M{"$inc": bson.M{field: 1}}

	// Build the filter to find the document with the matching ImageId
	filter := bson.M{
		"_id": labelSubtask.ImageId,
		field: bson.M{"$exists": true},
	}

	// return the document after the update
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	res := t.db.Collection(datasources.TasksCollection).FindOneAndUpdate(context.Background(), filter, update, opts)
	if res.Err() != nil {
		// TODO: if label not present what should I return? (for now "error": "mongo: no documents in result")
		return nil, res.Err()
	}

	err = res.Decode(&subtask)
	if err != nil {
		return nil, res.Err()
	}

	return subtask, nil
}

func (t *TaskRepository) CreateNewSubtask(subtask *models.Subtask) error {
	result, err := t.db.Collection(datasources.TasksCollection).InsertOne(context.Background(), subtask)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		subtask.Id = oid.Hex()
	}

	return nil
}