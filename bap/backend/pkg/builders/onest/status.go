package onest

import (
	"fmt"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
)

func BuildBPPStatusJobRequest(payload statusrequest.SeekerStatusPayload, bppId, bppuri string, worker *workerProfile.WorkerProfile) (*statusrequest.StatusRequest, error) {
	cityCode, err := utils.GetCityCode(payload.Location.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker city code for %s, %v", payload.Location.City, err)
	}
	countryCode, err := utils.GetCountryCode(payload.Location.Country)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker country code for %s, %v", payload.Location.Country, err)
	}
	req := statusrequest.StatusRequest{
		Context: statusrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "status",
			Version:       "2.0.0",
			TransactionID: worker.TransactionID,
			MessageID:     worker.MessageID,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: statusrequest.Location{
				City: statusrequest.City{
					Code: cityCode,
				},
				Country: statusrequest.Country{
					Code: countryCode,
				},
			},
		},
		Message: statusrequest.Message{
			Order: statusrequest.Order{
				ID: worker.ApplicantionID[worker.TransactionID],
			},
		},
	}
	return &req, nil
}