package main

import (
	"fmt"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/log"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/proxy"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/server"
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
	businessClient, jobClient, jobApplicationClient, initJobApplication := server.InitMongoDB()

	// Set up clients
	clients := clients.NewClients(jobClient, businessClient, jobApplicationClient, initJobApplication)

	// initialize the server
	server := server.SetupServer(clients)

	// start the server
	if err := server.Run(fmt.Sprintf(":%s", config.Config.HTTPPort)); err != nil {
		logrus.Fatal(err)
	}
}
