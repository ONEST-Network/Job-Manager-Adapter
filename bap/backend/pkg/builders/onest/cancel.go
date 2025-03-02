package onest


import (
	"fmt"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
)

func BuildBPPCancelJobRequest(payload cancelrequest.SeekerCancelPayload, bppId, bppuri string, worker *workerProfile.WorkerProfile) (*cancelrequest.CancelRequest, error) {
	cityCode, err := utils.GetCityCode(payload.Location.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker city code for %s, %v", payload.Location.City, err)
	}
	countryCode, err := utils.GetCountryCode(payload.Location.Country)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker country code for %s, %v", payload.Location.Country, err)
	}
	req := cancelrequest.CancelRequest{
		Context: cancelrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "cancel",
			Version:       "2.0.0",
			TransactionID: worker.TransactionID,
			MessageID:     worker.MessageID,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: cancelrequest.Location{
				City: cancelrequest.City{
					Code: cityCode,
				},
				Country: cancelrequest.Country{
					Code: countryCode,
				},
			},
		},
		Message: cancelrequest.Message{
			OrderID: worker.ApplicantionID[worker.TransactionID],
		},
	}
	return &req, nil
}