package clients

import (
	dbBusiness "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/business"
	dbJob "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/job"
	dbJobApplication "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/job-application"
)

type Clients struct {
	JobClient            *dbJob.Dao
	BusinessClient       *dbBusiness.Dao
	JobApplicationClient *dbJobApplication.Dao
}

func NewClients(jobClient *dbJob.Dao, businessClient *dbBusiness.Dao, jobApplicationClient *dbJobApplication.Dao) *Clients {
	return &Clients{
		JobClient:            jobClient,
		BusinessClient:       businessClient,
		JobApplicationClient: jobApplicationClient,
	}
}
