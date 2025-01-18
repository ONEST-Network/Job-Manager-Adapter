package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/internal/job"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	jobDb "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/job"
)

// @Summary	Create job
// @Description	Create a job posting
// @Accept		json
// @Produce		json
// @Param request body jobDb.Job true "request body"
// @Success 200 {object}
// @Failure 500 {object} string
// @Router	/job/create	[post]
func CreateJob(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload jobDb.Job

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

// @Summary	Send jobs
// @Description	Send jobs
// @Accept		json
// @Produce		json
// @Param request body request.SearchRequest true "request body"
// @Success 200 {object} response.SearchResponse
// @Failure 500 {object} string
// @Router	/search	[post]
func SendJobs(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		job := job.NewJob(clients)

		payload, ack := job.SendJobsAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go job.SendJobs(payload)
	}
}
