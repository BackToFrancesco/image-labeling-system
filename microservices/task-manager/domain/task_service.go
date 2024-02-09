package domain

import "fabc.it/task-manager/models"

type TaskService interface {
	CreateNewTask(task *models.Task) error
	GetTask(taskId string) (*models.Task, error)
	UpdateTask(task *models.Task) error
	UpdateSubtask(subtask *models.Subtask) error
}
