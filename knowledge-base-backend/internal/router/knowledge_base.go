package router

import (
	"github.com/gin-gonic/gin"

	"dev.choveylee.top/knowledge-base-backend/internal/handler"
)

func registerKnowledgeBase(router *gin.RouterGroup) {
	router.GET("/knowledge-bases", handler.HandleListKnowledgeBases)

	router.GET("/knowledge-bases/:id", handler.HandleGetKnowledgeBase)
	router.POST("/knowledge-bases", handler.HandleCreateKnowledgeBase)
	router.PUT("/knowledge-bases/:id", handler.HandleUpdateKnowledgeBase)
	router.DELETE("/knowledge-bases/:id", handler.HandleDeleteKnowledgeBase)
}
