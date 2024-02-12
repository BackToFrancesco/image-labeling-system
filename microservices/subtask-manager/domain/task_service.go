package domain

import "fabc.it/subtask-manager/models"

type TaskService interface {
	GetSubtasks(numberOfSubtasks int, userId string) (task *[]models.Subtask, err error)
	UpdateSubtaskLabels(labelSubtask *models.LabelSubtask) (err error)
}
