package api

import (
	//"fmt"
	"fabc.it/subtask-manager/controllers"
	//"github.com/gin-gonic/gin"
)

type TaskRoutes struct {
	handler    *RequestHandler
	controller *controllers.TaskController
}

func (t *TaskRoutes) setRoutes() {
	api := t.handler.Group("/api")
	/*
	api.POST("/tasks", func(c *gin.Context) {
        // Here, you can add any pre-processing logic before creating a new task
        // For example, validating the request body, checking for permissions, etc.
        fmt.Println("Pre-processing before creating a new task")

        // Call the actual controller method to create a new task
        //t.controller.CreateNewTask(c)

        // Optionally, you can add any post-processing logic here
        fmt.Println("Post-processing after creating a new task")
    })
	api.GET("/tasks", func(c *gin.Context) {
        // Here, you can add any pre-processing logic before creating a new task
        // For example, validating the request body, checking for permissions, etc.
        fmt.Println("Pre-processing before creating a new task")

        // Call the actual controller method to create a new task
        //t.controller.CreateNewTask(c)

        // Optionally, you can add any post-processing logic here
        fmt.Println("Post-processing after creating a new task")
    })*/
	api.GET("/ask-images", t.controller.SendImages) //fare qualcosa
	//api.POST("/tasks", t.controller.CreateNewTask) //fare qualcosa
	//api.POST("/tasks/:taskId/upload", t.controller.UploadTaskImages) //fare qualcose
	
}

func NewTaskRoutes(
	handler *RequestHandler,
	controller *controllers.TaskController,
) *TaskRoutes {
	routes := &TaskRoutes{
		handler:    handler,
		controller: controller,
	}

	routes.setRoutes()

	return routes
}
