package main

import (
	"fabc.it/subtask-manager/api"
	"fabc.it/subtask-manager/config"
	"fabc.it/subtask-manager/controllers"
	"fabc.it/subtask-manager/datasources"
	"fabc.it/subtask-manager/repositories"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Constructors,
		datasources.Constructors,
		repositories.Constructors,
		controllers.Constructors,
		api.Constructors,
		fx.Invoke(func(engine *gin.Engine) {}),
		fx.Invoke(func(route *api.TaskRoutes) {}),
	)

	app.Run()
}
