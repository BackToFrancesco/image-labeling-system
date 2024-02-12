package api

import (
	"fabc.it/subtask-manager/controllers"
)

type TaskRoutes struct {
	handler    *RequestHandler
	controller *controllers.TaskController
}

func (t *TaskRoutes) setRoutes() {
	api := t.handler.Group("/api")

	api.GET("/ask-images", t.controller.SendImages) //TODO: decide a name
	api.POST("/update-subtask-labels", t.controller.UpdateSubtaskLabels) //TODO: decide a name 
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
