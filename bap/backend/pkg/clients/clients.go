package clients

import (
	apiclient "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/api-client"
	dbJob "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"
	dbWorker "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	dbInitJobApplication "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/init-job-application"
)

type Clients struct {
	ApiClient           apiclient.Interface
	WorkerProfileClient *dbWorker.Dao
	JobClient           *dbJob.Dao
	InitJobApplicationClient *dbInitJobApplication.Dao
}

func NewClients(jobClient *dbJob.Dao, workerProfileClient *dbWorker.Dao) *Clients {
	return &Clients{
		ApiClient:           apiclient.NewAPIClient(),
		WorkerProfileClient: workerProfileClient,
		JobClient:                jobClient,
	}
}
