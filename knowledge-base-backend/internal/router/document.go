package router

import (
	"github.com/gin-gonic/gin"

	"dev.choveylee.top/knowledge-base-backend/internal/handler"
)

func registerDocument(router *gin.RouterGroup) {
	router.GET("/documents", handler.HandleListDocuments)

	router.GET("/documents/:id", handler.HandleGetDocument)
	router.POST("/documents", handler.HandleCreateDocument)
	router.PUT("/documents/:id", handler.HandleUpdateDocument)
	router.DELETE("/documents/:id", handler.HandleDeleteDocument)
}
