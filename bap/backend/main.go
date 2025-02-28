package main

import (
	"fmt"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/log"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/proxy"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/server"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	// init logrus logger
	log.InitLogger()

	// parse env variables
	if err := envconfig.Process("", &config.Config); err != nil {
		logrus.Fatalf("Failed to parse ENVs, %v", err)
	}

	// set proxy envs, if provided
	proxy.SetProxyENVs()

	// Initialize mongodb clients
	seekerClient, jobClient := server.InitMongoDB()

	// Set up clients
	clients := clients.NewClients(jobClient, seekerClient)

	// initialize the server
	server := server.SetupServer(clients)

	// Initialize job sync from BPP side
	jobSync := utils.NewJobSync(clients, 5*time.Minute)
    jobSync.Start()
    defer jobSync.Stop()

	// start the server
	if err := server.Run(fmt.Sprintf(":%s", config.Config.HTTPPort)); err != nil {
		logrus.Fatal(err)
	}
}

// Write an api in bap/backend/main.go to call this /recommend_jobs with required params