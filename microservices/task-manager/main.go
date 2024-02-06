package main

import (
	"fabc.it/task-manager/api"
	"fabc.it/task-manager/config"
	"fabc.it/task-manager/controllers"
	"fabc.it/task-manager/datasources"
	"fabc.it/task-manager/repositories"
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
