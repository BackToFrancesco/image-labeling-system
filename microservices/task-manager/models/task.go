package models

type Task struct {
	Id       string     `bson:"-"`
	Name     string     `json:"name" binding:"required"`
	Labels   []*string  `json:"labels" binding:"required"`
	Subtasks []*Subtask `json:"-"`
}
