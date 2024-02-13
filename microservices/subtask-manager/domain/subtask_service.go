package domain

import "fabc.it/subtask-manager/models"

type SubtaskService interface {
	CreateNewSubtask(task *models.Subtask) error
	GetSubtasks(numberOfSubtasks int, userId string) (task []*models.SubtaskMessage, err error)
	UpdateSubtaskLabel(labelSubtask *models.LabelSubtask) (subtask *models.Subtask, err error)
}
