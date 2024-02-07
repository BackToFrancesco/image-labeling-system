package api

import (
	"fabc.it/task-manager/controllers"
)

type TaskRoutes struct {
	handler    *RequestHandler
	controller *controllers.TaskController
}

func (t *TaskRoutes) setRoutes() {
	api := t.handler.Group("/api")

	api.POST("/tasks", t.controller.CreateNewTask)
	api.POST("/tasks/:taskId/upload", t.controller.UploadTaskImages)
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
