package api

import (
	"fabc.it/subtask-manager/controllers"
)

type SubtaskRoutes struct {
	handler    *RequestHandler
	controller *controllers.SubtaskController
}

func (t *SubtaskRoutes) setRoutes() {
	api := t.handler.Group("/api")

	api.GET("/ask-images", t.controller.GetSubtasks) //TODO: decide a name
	api.POST("/update-subtask-label", t.controller.UpdateSubtaskLabel) //TODO: decide a name 
}

func NewTaskRoutes(
	handler *RequestHandler,
	controller *controllers.SubtaskController,
) *SubtaskRoutes {
	routes := &SubtaskRoutes{
		handler:    handler,
		controller: controller,
	}

	routes.setRoutes()

	return routes
}
