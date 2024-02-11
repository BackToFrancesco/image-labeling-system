package controllers

import (
	//"archive/zip"
	"fabc.it/subtask-manager/domain"
	//"fabc.it/subtask-manager/models"
	"fmt"

	"github.com/gin-gonic/gin"

	//"log"
	//"mime/multipart"
	"net/http"
	//"os"
	//"slices"
	//"sync"
)

const (
	taskId = "taskId"
)

type TaskController struct {
	taskService          domain.TaskService
	storageService       domain.StorageService
	//messageBrokerService domain.MessageBrokerService
}

func NewTaskController(
	taskService domain.TaskService,
	storageService domain.StorageService,
	//messageBrokerService domain.MessageBrokerService,
) *TaskController {
	taskController := &TaskController{
		taskService:          taskService,
		storageService:       storageService,
		//messageBrokerService: messageBrokerService,
	}
	/*
	go func() {
		taskController.messageBrokerService.ConsumeCompletedSubtasks(taskController.consumeCompletedTask)
	}()*/

	return taskController
}
/*
func (t *TaskController) consumeCompletedTask(message *models.CompletedSubtaskMessage) error {
	m := map[string]int{}

	var maxCnt int
	var mostFreq string

	for _, label := range message.AssignedLabels {
		m[*label]++
		if m[*label] > maxCnt {
			maxCnt = m[*label]
			mostFreq = *label
		}
	}

	err := t.taskService.UpdateSubtask(&models.Subtask{
		Id:    message.Id,
		Label: mostFreq,
	})
	if err != nil {
		return err
	}

	return nil
}
*/
/*
func (t *TaskController) CreateNewTask(c *gin.Context) {
	input := &models.Task{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = t.taskService.CreateNewTask(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": input.Id})
}
*/
func (t *TaskController) SendImages(c *gin.Context) {
	fmt.Println("Images asked from user")
	//err := db
	// take subtask with less labelings
	// send 100 images to classify
	// 
	res, err := t.taskService.GetSubtasks(10, "userID3")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println("Tasks2:", res)
	c.JSON(http.StatusOK, gin.H{"tasks": res})
}
/*
func (t *TaskController) UploadTaskImages(c *gin.Context) {
	taskId := c.Param(taskId)

	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no task id specified"})
		return
	}

	task, err := t.taskService.GetTask(taskId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	fileHeader, err := c.FormFile("images")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}(file)

	bytes := make([]byte, 512)

	_, err = file.Read(bytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if contentType := http.DetectContentType(bytes); contentType != "application/zip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not a zip file"})
		return
	}

	destination := fmt.Sprintf("%s", fileHeader.Filename)
	err = c.SaveUploadedFile(fileHeader, destination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filesInZip, err := zip.OpenReader(destination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer func(filesInZip *zip.ReadCloser) {
		err := filesInZip.Close()
		if err != nil {
			log.Print(err)
		}
	}(filesInZip)

	var wg sync.WaitGroup
	resultChan := make(chan string)

	for i, file := range filesInZip.File {
		wg.Add(1)

		go func(wg *sync.WaitGroup, resultChan chan<- string, taskId string, file *zip.File) {
			defer wg.Done()

			err := t.storageService.SaveImage(taskId, file)
			if err != nil {
				log.Print(err)
				return
			}

			resultChan <- taskId

			return
		}(&wg, resultChan, fmt.Sprintf(fmt.Sprintf("%s-%d", taskId, i)), file)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	newSubtasks := make([]*models.Subtask, 0)
	for r := range resultChan {
		newSubtasks = slices.Insert(newSubtasks, len(newSubtasks), &models.Subtask{Id: r})
	}

	task.Subtasks = newSubtasks

	err = t.taskService.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, subtask := range task.Subtasks {
		go func(subtask *models.Subtask) {
			err := t.messageBrokerService.PublishNewSubtask(&models.SubtaskMessage{
				Id:     subtask.Id,
				Labels: task.Labels,
			})
			if err != nil {
				log.Print(err)
				return
			}
		}(subtask)
	}

	c.Status(http.StatusOK)

	defer func() {
		err := os.Remove(destination)
		if err != nil {
			return
		}
	}()
}
*/