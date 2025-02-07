package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/internal/onest"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
)

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

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SendJobsAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.SendJobs(payload)
	}
}

// @Summary	Send job fulfillment
// @Description	Send job fulfillment
// @Accept		json
// @Produce		json
// @Param request body request.SelectRequest true "request body"
// @Success 200 {object} response.SelectResponse
// @Failure 500 {object} string
// @Router	/select	[post]
func SendJobFulfillment(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SendJobFulfillmentAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.SendJobFulfillment(payload)
	}
}

// @Summary	Initialize job application
// @Description	Initialize job application
// @Accept		json
// @Produce		json
// @Param request body request.InitRequest true "request body"
// @Success 200 {object} response.InitResponse
// @Failure 500 {object} string
// @Router	/init	[post]
func InitializeJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.InitializeJobApplicationAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.InitializeJobApplication(payload)
	}
}

// @Summary	Confirm job application submission
// @Description	Confirm job application submission
// @Accept		json
// @Produce		json
// @Param request body request.ConfirmRequest true "request body"
// @Success 200 {object} response.ConfirmResponse
// @Failure 500 {object} string
// @Router	/confirm	[post]
func ConfirmJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, initJobApplication, ack := onest.ConfirmJobApplicationAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.ConfirmJobApplication(payload, initJobApplication)
	}
}

// @Summary	Send job application current status
// @Description	Send job application current status
// @Accept		json
// @Produce		json
// @Param request body request.StatusRequest true "request body"
// @Success 200 {object} response.StatusResponse
// @Failure 500 {object} string
// @Router	/status	[post]
func JobApplicationStatus(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.JobApplicationStatusAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.JobApplicationStatus(payload)
	}
}

// @Summary	Cancel job application
// @Description	Cancel job application
// @Accept		json
// @Produce		json
// @Param request body request.CancelRequest true "request body"
// @Success 200 {object} response.CancelResponse
// @Failure 500 {object} string
// @Router	/cancel	[post]
func WithdrawJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.WithdrawJobApplicationAck(c.Request.Body)
		if ack.Error != nil {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.WithdrawJobApplication(payload)
	}
}
