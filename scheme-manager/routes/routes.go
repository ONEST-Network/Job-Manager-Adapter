package routes

import (
	"github.com/Aerospace-prog/scheme-manager/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")

	api.POST("/schemes", controllers.PushScheme)
	api.POST("/applications", controllers.CreateApplication) 
	api.GET("/applications", controllers.GetApplications)
	api.GET("/applications/:id", controllers.GetApplicationByID)
	api.PUT("/applications/:id/status", controllers.UpdateApplicationStatus)
}