package domain

import "fabc.it/subtask-manager/models"

type MessageBrokerService interface {
	PublishCompletedSubtask(subtask *models.CompletedSubtaskMessage) error
	ConsumeNewSubtasks(consume func(message *models.SubtaskMessage) error)
}
