package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/internal/job"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	jobPayload "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/job"
	"github.com/gin-gonic/gin"
)

// @Summary	Create job
// @Description	Create a job posting
// @Accept		json
// @Produce		json
// @Param request body jobPayload.CreateJobRequest true "request body"
// @Success 200
// @Failure 500 {object} string
// @Router	/job/create	[post]
func CreateJob(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload jobPayload.CreateJobRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := job.NewJob(clients).CreateJob(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
}
