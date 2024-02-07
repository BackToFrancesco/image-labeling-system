package models

type Task struct {
	Id       string     `bson:"_id,omitempty"`
	Name     string     `json:"name" binding:"required" bson:"name"`
	Labels   []*string  `json:"labels" binding:"required" bson:"labels"`
	Subtasks []*Subtask `json:"-" bson:"subtasks"`
}
