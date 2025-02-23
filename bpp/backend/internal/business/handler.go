package business

import (
	"fmt"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	businessDb "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/business"
	businessPayload "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/business"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/utils/random"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Interface interface {
	AddBusiness(business *businessPayload.AddBusinessRequest) (string, error)
	ListJobs(businessID string) ([]businessPayload.ListJobsResponse, error)
}

type Business struct {
	clients *clients.Clients
}

func NewBusiness(clients *clients.Clients) Interface {
	return &Business{
		clients: clients,
	}
}

func (b *Business) AddBusiness(payload *businessPayload.AddBusinessRequest) (string, error) {
	logrus.Infof("[Request]: Received request to add a new business: %s", payload.Name)

	var business = &businessDb.Business{
		ID:             random.GetRandomString(7),
		Name:           payload.Name,
		Phone:          payload.Phone,
		Email:          payload.Email,
		PictureURLs:    payload.PictureURLs,
		Description:    payload.Description,
		GSTIndexNumber: payload.GSTIndexNumber,
		Location:       payload.Location,
		Industry:       payload.Industry,
	}

	if err := b.clients.BusinessClient.CreateBusiness(business); err != nil {
		logrus.Errorf("Failed to create a new business, %v", err)
		return "", fmt.Errorf("failed to create a new business, %v", err)
	}

	return business.ID, nil
}

func (b *Business) ListJobs(businessID string) ([]businessPayload.ListJobsResponse, error) {
	logrus.Infof("[Request]: Received request to get jobs for business: %s", businessID)

	var query = bson.D{{Key: "business.id", Value: businessID}}

	jobs, err := b.clients.JobClient.ListJobs(query)
	if err != nil {
		logrus.Errorf("Failed to get jobs for business %s, %v", businessID, err)
		return nil, fmt.Errorf("failed to get jobs for business, %v", err)
	}

	var listJobsResponse []businessPayload.ListJobsResponse
	for _, job := range jobs {
		listJobsResponse = append(listJobsResponse, businessPayload.ListJobsResponse{
			ID:             job.ID,
			Name:           job.Name,
			Description:    job.Description,
			Type:           job.Type,
			Vacancies:      job.Vacancies,
			SalaryRange:    job.SalaryRange,
			ApplicationIDs: job.ApplicationIDs,
			WorkHours:      job.WorkHours,
			WorkDays:       job.WorkDays,
			Eligibility:    job.Eligibility,
			Location:       job.Location,
		})
	}

	return listJobsResponse, nil
}
