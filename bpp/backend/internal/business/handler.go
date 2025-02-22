package business

import (
	"fmt"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	businessDb "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/business"
	businessPayload "github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/business"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/utils/random"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	AddBusiness(business *businessPayload.AddBusinessRequest) (string, error)
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
