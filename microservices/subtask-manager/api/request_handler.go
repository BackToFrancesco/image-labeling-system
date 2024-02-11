package api

import "github.com/gin-gonic/gin"

// RequestHandler is a wrapper of gin.Engine
type RequestHandler struct {
	*gin.Engine
}

func NewRequestHandler(engine *gin.Engine) *RequestHandler {
	return &RequestHandler{engine}
}
