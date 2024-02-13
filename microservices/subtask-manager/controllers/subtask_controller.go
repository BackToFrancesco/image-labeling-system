package controllers

import (
	"fabc.it/subtask-manager/domain"
	"fabc.it/subtask-manager/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	// Number of labels required to complete a subtask
	labelsRequired = 10
)

const subtaskId = "subtaskId"

type SubtaskController struct {
	subtaskService domain.SubtaskService
	//storageService       domain.StorageService
	messageBrokerService domain.MessageBrokerService
}

func NewSubtaskController(
	subtaskService domain.SubtaskService,
	//storageService domain.StorageService,
	messageBrokerService domain.MessageBrokerService,
) *SubtaskController {
	subtaskController := &SubtaskController{
		subtaskService: subtaskService,
		//storageService:       storageService,
		messageBrokerService: messageBrokerService,
	}

	go func() {
		subtaskController.messageBrokerService.ConsumeNewSubtasks(subtaskController.ConsumeNewSubTask)
	}()

	return subtaskController
}

func (t *SubtaskController) ConsumeNewSubTask(message *models.SubtaskMessage) error {

	// Translation of SubtaskMessage to Subtask
	assignedLabels := make(map[string]int)

	for _, labelPtr := range message.Labels {
		if labelPtr != nil {
			assignedLabels[*labelPtr] = 0
		}
	}

	subtask := models.Subtask{
		Id:             message.Id,
		Labels:         message.Labels,
		Assignee:       []string{},
		AssignedLabels: assignedLabels,
	}

	err := t.subtaskService.CreateNewSubtask(&subtask)
	if err != nil {
		return err
	}

	return nil
}

func (t *SubtaskController) PublishCompletedSubtask(message *models.Subtask) error {

	completedSubtask := models.CompletedSubtaskMessage{
		Id:             message.Id,
		AssignedLabels: &message.AssignedLabels,
	}

	err := t.messageBrokerService.PublishCompletedSubtask(&completedSubtask)
	if err != nil {
		return err
	}

	return nil
}

func (t *SubtaskController) GetSubtasks(c *gin.Context) {

	input := &models.RequestSubtasks{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := t.subtaskService.GetSubtasks(input.NumberOfSubtasks, input.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subtasks": res})
}

func (t *SubtaskController) UpdateSubtaskLabel(c *gin.Context) {
	input := &models.LabelSubtask{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ImageId = c.Param(subtaskId)

	res, err := t.subtaskService.UpdateSubtaskLabel(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: refactor in a function?
	// total number of assigned labels
	totalLabels := 0
	for _, count := range res.AssignedLabels {
		totalLabels += count
	}

	// publish message in RabbitMq
	if labelsRequired <= totalLabels {
		err = t.PublishCompletedSubtask(res)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"newTotal:": totalLabels})
}
