package domain

import "fabc.it/task-manager/models"

type TaskService interface {
	CreateNewTask(task *models.Task) error
}
