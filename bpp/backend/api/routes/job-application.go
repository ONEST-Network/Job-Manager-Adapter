package routes

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	"github.com/gin-gonic/gin"
)

func JobApplicationRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/:id/status", handlers.UpdateJobApplicationStatus(clients))
}
