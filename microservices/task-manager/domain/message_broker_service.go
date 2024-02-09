package domain

import "fabc.it/task-manager/models"

type MessageBrokerService interface {
	PublishNewSubtask(subtask *models.SubtaskMessage) error

	ConsumeCompletedSubtasks(consume func(message *models.CompletedSubtaskMessage) error)
}
