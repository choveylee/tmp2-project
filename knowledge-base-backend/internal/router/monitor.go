package router

import (
	"github.com/gin-gonic/gin"

	"dev.choveylee.top/knowledge-base-backend/internal/handler"
)

func registerMonitor(router *gin.Engine) {
	router.GET("/cpu-check", handler.HandleCpuCheck)
	router.GET("/ram-check", handler.HandleRamCheck)
}
