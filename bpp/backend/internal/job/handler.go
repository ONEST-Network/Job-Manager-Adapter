package job

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	jobDb "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/job"
	jobPayload "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/job"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/utils/random"
)

type Interface interface {
	CreateJob(payload *jobPayload.CreateJobRequest) error
}

type Job struct {
	clients *clients.Clients
}

func NewJob(clients *clients.Clients) Interface {
	return &Job{
		clients: clients,
	}
}

func (j *Job) CreateJob(payload *jobPayload.CreateJobRequest) error {
	logrus.Infof("Received request to create a new job for business: %s", payload.BusinessID)

	business, err := j.clients.BusinessClient.GetBusiness(payload.BusinessID)
	if err != nil {
		return fmt.Errorf("failed to get business with id %s, %v", payload.BusinessID, err)
	}

	var job = &jobDb.Job{
		ID:             random.GetRandomString(7),
		Name:           payload.Name,
		Description:    payload.Description,
		Type:           payload.Type,
		Vacancies:      payload.Vacancies,
		SalaryRange:    payload.SalaryRange,
		ApplicationIDs: payload.ApplicationIDs,
		WorkHours:      payload.WorkHours,
		WorkDays:       payload.WorkDays,
		Eligibility:    payload.Eligibility,
		Location:       payload.Location,
		Business:       *business,
	}

	if err := j.clients.JobClient.CreateJob(job); err != nil {
		return err
	}

	return nil
}
