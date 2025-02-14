package main

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/log"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
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

}
