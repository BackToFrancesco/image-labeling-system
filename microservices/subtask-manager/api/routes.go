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

	api.GET("/subtasks", t.controller.GetSubtasks)
	api.PATCH("/subtasks/:subtaskId", t.controller.UpdateSubtaskLabel)
}

func NewSubtaskRoutes(
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
