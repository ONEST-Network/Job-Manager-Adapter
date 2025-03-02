package onest

import (
	"fmt"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
)

func BuildBPPSelectJobRequest(payload selectrequest.SeekerSelectPayload, transactionId, messageId, bppId, bppuri string) (*selectrequest.SelectRequest, error) {
	cityCode, err := utils.GetCityCode(payload.Location.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker city code for %s, %v", payload.Location.City, err)
	}
	countryCode, err := utils.GetCountryCode(payload.Location.Country)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker country code for %s, %v", payload.Location.Country, err)
	}
	req := selectrequest.SelectRequest{
		Context: selectrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "select",
			Version:       "2.0.0",
			TransactionID: transactionId,
			MessageID:     messageId,
			BppID:         bppId,
			BppURI:        bppuri,
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC(),
			TTL:           "PT30S",
			Location: selectrequest.Location{
				City: selectrequest.City{
					Code: cityCode,
				},
				Country: selectrequest.Country{
					Code: countryCode,
				},
			},
		},
		Message: selectrequest.Message{
			Order: selectrequest.Order{
				Provider: selectrequest.Provider{
					ID: payload.ProviderID,
				},
				Fulfillments: []selectrequest.Fulfillments{
					{
						ID:   "F1",
						Type: "lead & recruitment",
					},
				},
				Items: []selectrequest.Items{
					{
						ID:             payload.JobID,
					},
				},
			},
		},
	}
	return &req, nil
}
