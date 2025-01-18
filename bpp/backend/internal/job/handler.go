package job

import (
	"encoding/json"
	"io"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/builders/job"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	jobDb "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/job"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/search/request"
	requestack "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/search/request-ack"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/search/response"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/utils/random"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Interface interface {
	CreateJob(payload *jobDb.Job) error
	SendJobsAck(body io.ReadCloser) (*request.SearchRequest, *requestack.SearchRequestAck)
	SendJobs(payload *request.SearchRequest) (*response.SearchResponse, error)
}

type Job struct {
	clients *clients.Clients
}

func NewJob(clients *clients.Clients) Interface {
	return &Job{
		clients: clients,
	}
}

func (j *Job) CreateJob(payload *jobDb.Job) error {
	logrus.Infof("Received request to create a new job for business: %s", payload.BusinessID)

	// populate job id
	payload.ID = random.GetRandomString(7)

	if err := j.clients.JobClient.CreateJob(payload); err != nil {
		return err
	}

	return nil
}

func (j *Job) SendJobsAck(body io.ReadCloser) (*request.SearchRequest, *requestack.SearchRequestAck) {
	var (
		payload      request.SearchRequest
		payloadError *requestack.Error
		status       = "ACK"
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		payloadError = &requestack.Error{
			Code:    "10000",
			Paths:   "",
			Message: err.Error(),
		}
	}

	if payloadError != nil {
		status = "NACK"
	}

	return &payload, &requestack.SearchRequestAck{
		Message: requestack.Message{
			Ack: requestack.Ack{
				Status: status,
			},
		},
		Error: payloadError,
	}
}

func (j *Job) SendJobs(payload *request.SearchRequest) (*response.SearchResponse, error) {
	var (
		searchTerm = payload.Message.Intent.Item.Descriptor.Name
		query      = bson.D{}
	)

	if searchTerm != "" {
		query = bson.D{
			{"$or", bson.A{
				bson.D{{"name", bson.D{{"$regex", searchTerm}, {"$options", "i"}}}},
				bson.D{{"description", bson.D{{"$regex", searchTerm}, {"$options", "i"}}}},
			}},
		}
	}

	jobs, err := j.clients.JobClient.ListJobs(query)
	if err != nil {
		logrus.Errorf("Error listing jobs: %v", err)
		return nil, err
	}

	return job.BuildJobsResponse(j.clients, payload, jobs)
}
