package routes

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	"github.com/gin-gonic/gin"
)

func BecknRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/search", handlers.SendJobs(clients))
	router.POST("/select", handlers.SendJobFulfillment(clients))
	router.POST("/init", handlers.InitializeJobApplication(clients))
	router.POST("/confirm", handlers.ConfirmJobApplication(clients))
	router.POST("/status", handlers.JobApplicationStatus(clients))
	router.POST("/cancel", handlers.WithdrawJobApplication(clients))
}
