package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
    dbSearchResponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/searchResponse"
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

		// c.JSON(statusCode, ack)

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
            // Get MongoDB collection where responses are stored
            responseCollection := h.onestService.Clients.SearchReponseClient.Collection
            // Create a context with timeout to prevent indefinite blocking
            ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
            defer cancel()

            parsedRequest, transaction_id, err := builders.BuildBPPSearchJobsRequest(payload)
            if err != nil {
                logrus.Errorf("Failed to parse search job request, %v", err)
                return
            }
            // Create a new document in the SearchResponseClient collection
            searchJobResponse := dbSearchResponse.SearchJobResponse {
                ID: transaction_id,
                TransactionID: transaction_id,
            }

            // Insert the document into the collection
            err = h.onestService.Clients.SearchReponseClient.CreateSearchJobResponse(&searchJobResponse)
            if err != nil {
                logrus.Errorf("Failed to create search response document: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create search response"})
                return
            }

            logrus.Infof("Created search response document with transaction_id: %s", transaction_id)

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
            logrus.Printf("Parsed Seeker Search Request: %+v\n", parsedRequest)
            _, err = h.onestService.Search(c.Request.Context(), parsedRequest)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            

            // Create a pipeline to match only documents with our transaction_id
            pipeline := mongo.Pipeline{
                bson.D{{Key: "$match", Value: bson.D{
                    {Key: "fullDocument.transaction_id", Value: transaction_id},
                    {Key: "updateDescription.updatedFields.jobs_response", Value: bson.D{{Key: "$exists", Value: true}}},
                }}},
            }

            // Set up the change stream options
            opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
            
            // Create the change stream
            changeStream, err := responseCollection.Watch(ctx, pipeline, opts)
            if err != nil {
                logrus.Errorf("Failed to create change stream: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to watch for response"})
                return
            }
            defer changeStream.Close(ctx)
            
            // Wait for a matching document
            logrus.Infof("Waiting for response for transaction ID: %s", parsedRequest.Context.TransactionID)

            var result bson.M
            if changeStream.Next(ctx) {
                if err := changeStream.Decode(&result); err != nil {
                    logrus.Errorf("Error decoding change stream document: %v", err)
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process response"})
                    return
                }
                
                // Extract the full document from the change stream result
                if fullDoc, ok := result["fullDocument"].(bson.M); ok {
                    // Return the jobs_response field to the client
                    if jobsResponse, ok := fullDoc["jobs_response"]; ok {
                        c.JSON(http.StatusOK, jobsResponse)
                    } else {
                        c.JSON(http.StatusOK, fullDoc) // Fallback to returning the entire document
                    }
                    return
                } else {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response format"})
                    return
                }
            }
            
            // If context timed out
            if ctx.Err() != nil {
                c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Response timeout"})
                return
            }

            // c.JSON(http.StatusOK, response)
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
