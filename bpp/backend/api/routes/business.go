package routes

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	"github.com/gin-gonic/gin"
)

func BusinessRouter(router *gin.RouterGroup, clients *clients.Clients) {
	router.POST("/add", handlers.AddBusiness(clients))
}
