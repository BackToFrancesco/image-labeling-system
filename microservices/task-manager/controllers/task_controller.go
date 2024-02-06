package controllers

import (
	"archive/zip"
	"fabc.it/task-manager/domain"
	"fabc.it/task-manager/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type TaskController struct {
	taskService domain.TaskService
}

func NewTaskController(
	taskService domain.TaskService,
) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

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

func (t *TaskController) UploadTaskImages(c *gin.Context) {
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

	for _, file := range filesInZip.File {
		fmt.Println(file.Name)
	}

	defer func() {
		err := os.Remove(destination)
		if err != nil {
			return
		}
	}()
}
