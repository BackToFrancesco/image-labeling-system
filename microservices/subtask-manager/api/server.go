package api

import (
	"context"
	"fabc.it/subtask-manager/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net"
	"net/http"
)

func NewTaskServer(lc fx.Lifecycle, env *config.Env) *gin.Engine {
	router := gin.Default()

	server := &http.Server{Addr: fmt.Sprintf(":%s", env.ServerPort), Handler: router}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				fmt.Println("Failed to start HTTP Server at ", server.Addr)
				return err
			}

			go func() {
				err := server.Serve(ln)
				if err != nil {

				}
			}()

			fmt.Println("Succeeded to start HTTP Server at ", server.Addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := server.Shutdown(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	})

	return router
}
