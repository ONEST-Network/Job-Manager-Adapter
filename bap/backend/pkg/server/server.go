package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/handlers"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/middleware"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/api/routes"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/docs"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb"
	dbWorker "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	dbJob "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"
)

func SetupServer(clients *clients.Clients, bppHandler *handlers.OnestBPPHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(middleware.DefaultStructuredLogger())
	server.Use(gin.Recovery())
	server.Use(middleware.ValidateCors())

	docs.SwaggerInfo.Title = "Worker Hub"
	docs.SwaggerInfo.Description = `Worker Hub is acting as seeker platform (BAP) which helps the workers to find jobs`
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	baseRouter := server.Group("/")
	routes.BaseRouter(baseRouter, clients)
	routes.BecknRouter(baseRouter, clients)
	routes.BPPOnestRoutes(baseRouter, bppHandler)
	routes.JobRecommendationRouter(baseRouter, clients)
	
	return server
}

func InitMongoDB() (*dbWorker.Dao, *dbJob.Dao) {
	var err error

	// Initialize mongodb clients
	mongodb.Client, err = mongodb.NewMongoClient()
	if err != nil {
		logrus.Fatalf("[Server]: Failed to connect mongo client, %v", err)
	}

	logrus.Info("[Server]: Connected To MongoDB")

	worker := dbWorker.NewWorkerDao(mongodb.Client.WorkerProfileCollection)
	job := dbJob.NewJobDao(mongodb.Client.JobCollection)


	return worker, job
}
