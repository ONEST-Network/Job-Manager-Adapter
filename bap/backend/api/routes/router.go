package routes

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/gin-gonic/gin"
)

func BecknRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/on_search", handlers.StoreJobs(clients))
	router.POST("/on_select", handlers.SendJobFulfillment(clients))
	router.POST("/on_init", handlers.InitializeJobApplication(clients))
	router.POST("/on_confirm", handlers.ConfirmJobApplication(clients))
	router.POST("/on_status", handlers.JobApplicationStatus(clients))
	router.POST("/on_cancel", handlers.WithdrawJobApplication(clients))
}


func RegisterOnestRoutes(router *gin.Engine, handler *handlers.OnestHandler) {
    group := router.Group("/onest")
    {
        group.POST("/search", handler.Search())
        group.POST("/select", handler.Select())
        group.POST("/init", handler.Init())
        group.POST("/confirm", handler.Confirm())
        group.POST("/status", handler.Status())
        group.POST("/cancel", handler.Cancel())
    }
}