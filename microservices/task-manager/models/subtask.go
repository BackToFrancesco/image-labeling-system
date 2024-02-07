package models

type Subtask struct {
	Id    string  `bson:"id"`
	Label *string `bson:"label"`
}
