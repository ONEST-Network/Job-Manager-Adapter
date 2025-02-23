package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/onest"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/service"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
	confirmrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/request"
	initrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/request"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
)

type OnestHandler struct {
    onestService *service.OnestService
}

func NewOnestHandler(onestService *service.OnestService) *OnestHandler {
    return &OnestHandler{
        onestService: onestService,
    }
}


func StoreJobs(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SendJobsAck(c.Request.Body)
		if ack.Error.Message != "" {
            statusCode = http.StatusBadRequest
            c.JSON(statusCode, ack)
            return
        }

		c.JSON(statusCode, ack)

		go onest.StoreJobs(payload)
	}
}

func SendJobFulfillment(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SendJobFulfillmentAck(c.Request.Body)
		if ack.Error.Message != "" {
            statusCode = http.StatusBadRequest
            c.JSON(statusCode, ack)
            return
        }

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.SendJobFulfillment(payload)
	}
}

func InitializeJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.InitializeJobApplicationAck(c.Request.Body)
		if ack.Error.Message != "" {
            statusCode = http.StatusBadRequest
            c.JSON(statusCode, ack)
            return
        }

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.InitializeJobApplication(payload)
	}
}

func ConfirmJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.ConfirmJobApplicationAck(c.Request.Body)
		if ack.Error.Message != "" {
            statusCode = http.StatusBadRequest
            c.JSON(statusCode, ack)
            return
        }

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.ConfirmJobApplication(payload)
	}
}

func JobApplicationStatus(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.JobApplicationStatusAck(c.Request.Body)
		if ack.Error.Message != "" {
            statusCode = http.StatusBadRequest
            c.JSON(statusCode, ack)
            return
        }

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.JobApplicationStatus(payload)
	}
}

func WithdrawJobApplication(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.WithdrawJobApplicationAck(c.Request.Body)
		if ack.Error.Message != "" {
            statusCode = http.StatusBadRequest
            c.JSON(statusCode, ack)
            return
        }

		c.JSON(statusCode, ack)

		// TODO: Implement a message queue to push the payload for processing
		go onest.WithdrawJobApplication(payload)
	}
}

/**
BPP APIs
**/

func (h *OnestHandler) Search() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload searchrequest.SearchRequest
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Set required context fields
        payload.Context = searchrequest.Context{
            Domain:        "jobs",
            Action:        "search",
            Version:       "1.0.0",
            BapID:        "example-bap",
            BapURI:       "https://example.com/callback",
            TransactionID: "tx-001",
            MessageID:     "msg-001",
            Timestamp:     time.Now().Format(time.RFC3339),
        }

        response, err := h.onestService.Search(c.Request.Context(), &payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, response)
    }
}

func (h *OnestHandler) Select() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload selectrequest.SelectRequest
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        payload.Context = selectrequest.Context{
            Domain:        "jobs",
            Action:        "select",
            Version:       "1.0.0",
            TransactionID: "tx-001",
            MessageID:     "msg-001",
            BapURI:       "https://example.com/callback",
        }

        response, err := h.onestService.Select(c.Request.Context(), &payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, response)
    }
}

func (h *OnestHandler) Init() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload initrequest.InitRequest
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        payload.Context = initrequest.Context{
            Domain:        "jobs",
            Action:        "init",
            Version:       "1.0.0",
            TransactionID: "tx-001",
            MessageID:     "msg-001",
            BapURI:       "https://example.com/callback",
        }

        response, err := h.onestService.Init(c.Request.Context(), &payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, response)
    }
}

func (h *OnestHandler) Confirm() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload confirmrequest.ConfirmRequest
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        payload.Context = confirmrequest.Context{
            Domain:        "jobs",
            Action:        "confirm",
            Version:       "1.0.0",
            TransactionID: "tx-001",
            MessageID:     "msg-001",
            BapURI:       "https://example.com/callback",
        }

        response, err := h.onestService.Confirm(c.Request.Context(), &payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, response)
    }
}

func (h *OnestHandler) Status() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload statusrequest.StatusRequest
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        payload.Context = statusrequest.Context{
            Domain:        "jobs",
            Action:        "status",
            Version:       "1.0.0",
            TransactionID: "tx-001",
            MessageID:     "msg-001",
            BapURI:       "https://example.com/callback",
        }

        response, err := h.onestService.Status(c.Request.Context(), &payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, response)
    }
}

func (h *OnestHandler) Cancel() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload cancelrequest.CancelRequest
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        payload.Context = cancelrequest.Context{
            Domain:        "jobs",
            Action:        "cancel",
            Version:       "1.0.0",
            TransactionID: "tx-001",
            MessageID:     "msg-001",
            BapURI:       "https://example.com/callback",
        }

        response, err := h.onestService.Cancel(c.Request.Context(), &payload)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, response)
    }
}
