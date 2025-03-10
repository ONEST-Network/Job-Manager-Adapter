package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/onest"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/internal/service"
	builders "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/builders/onest"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
	confirmrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/request"
	initrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/request"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
)

type OnestBPPHandler struct {
	onestService *service.OnestBPPService
}

func NewOnestBPPHandler(onestService *service.OnestBPPService) *OnestBPPHandler {
	return &OnestBPPHandler{
		onestService: onestService,
	}
}

func SearchJobs(clients *clients.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode = http.StatusOK

		onest := onest.NewOnestClient(clients)

		payload, ack := onest.SearchJobsAck(c.Request.Body)
		if ack.Error.Message != "" {
			statusCode = http.StatusBadRequest
			c.JSON(statusCode, ack)
			return
		}

		c.JSON(statusCode, ack)

		go onest.SearchJobs(payload)
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

func (h *OnestBPPHandler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload searchrequest.SeekerSearchPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        // Store transaction ID and message ID in worker profile
        if payload.WorkerID != "" {
            parsedRequest, err := builders.BuildBPPSearchJobsRequest(payload)
            if err != nil {
                logrus.Errorf("Failed to parse search job request, %v", err)
                return
            }
            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
            updateFields := bson.D{{Key: "$set", Value: bson.D{
                {Key: "transaction_id", Value: parsedRequest.Context.TransactionID},
                {Key: "message_id", Value: parsedRequest.Context.MessageID},
            }}}
            
            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                logrus.Errorf("Failed to update worker profile with transaction ID %s: %v",parsedRequest.Context.TransactionID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            fmt.Printf("Parsed Seeker Search Request: %+v\n", parsedRequest)
            response, err := h.onestService.Search(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, response)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Select() gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload selectrequest.SeekerSelectPayload
        if err := c.ShouldBindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            parsedRequest, err := builders.BuildBPPSelectJobRequest(payload, worker.TransactionID, worker.MessageID, payload.BppID, payload.BppURI)
            if err != nil {
                logrus.Errorf("Failed to parse select job request, %v", err)
                return
            }
            response, err := h.onestService.Select(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, response)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Init() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload initrequest.SeekerInitPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            parsedRequest, err := builders.BuildBPPInitJobRequest(payload, worker.TransactionID, worker.MessageID, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse init job request, %v", err)
                return
            }
            response, err := h.onestService.Init(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, response)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Confirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload confirmrequest.SeekerConfirmPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            parsedRequest, err := builders.BuildBPPConfirmJobRequest(payload, worker.TransactionID, worker.MessageID, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse confirm job request, %v", err)
                return
            }
            updateQuery := bson.D{{Key: "id", Value: payload.WorkerID}}
            updateFields := bson.D{{Key: "$set", Value: bson.D{
                {Key: "application_id", Value: map[string]string{parsedRequest.Context.TransactionID: parsedRequest.Message.Order.ID}},
            }}}
            
            if err := h.onestService.Clients.WorkerProfileClient.UpdateWorkerProfile(updateQuery, updateFields); err != nil {
                logrus.Errorf("Failed to update worker profile with application ID %s: %v", parsedRequest.Message.Order.ID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            response, err := h.onestService.Confirm(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, response)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Status() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload statusrequest.SeekerStatusPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            parsedRequest, err := builders.BuildBPPStatusJobRequest(payload, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse status job request, %v", err)
                return
            }
            response, err := h.onestService.Status(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, response)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}

func (h *OnestBPPHandler) Cancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload cancelrequest.SeekerCancelPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        if payload.WorkerID != "" {
            worker, err := h.onestService.Clients.WorkerProfileClient.GetWorkerProfile(payload.WorkerID)
            if err != nil {
                logrus.Errorf("Failed to get worker profile with worker ID %s: %v", payload.WorkerID, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            parsedRequest, err := builders.BuildBPPCancelJobRequest(payload, payload.BppID, payload.BppURI, worker)
            if err != nil {
                logrus.Errorf("Failed to parse cancel job request, %v", err)
                return
            }
            response, err := h.onestService.Cancel(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, response)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "please provide a valid worker ID"})
            return
        }
	}
}
