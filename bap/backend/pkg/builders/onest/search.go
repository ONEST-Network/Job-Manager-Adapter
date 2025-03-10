package onest

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils"
)

func BuildBPPSearchJobsRequest(payload searchrequest.SeekerSearchPayload) (*searchrequest.SearchRequest, error) {
	cityCode, err := utils.GetCityCode(payload.Location.City)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker city code for %s, %v", payload.Location.City, err)
	}
	stateCode, err := utils.GetStateCode(payload.Location.State)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker city code for %s, %v", payload.Location.City, err)
	}
	countryCode, err := utils.GetCountryCode(payload.Location.Country)
	if err != nil {
		return nil, fmt.Errorf("failed to get seeker country code for %s, %v", payload.Location.Country, err)
	}
	req := searchrequest.SearchRequest{
		Context: searchrequest.Context{
			Domain:        "ONDC:ONEST10",
			Action:        "search",
			Version:       "2.0.0",
			TransactionID: uuid.New().String(),
            MessageID:     uuid.New().String(),
			BapID:         config.Config.BapId,
			BapURI:        config.Config.BapUri,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			TTL:           "PT30S",
			Location: searchrequest.Location{
				City: searchrequest.City{
					Code: cityCode,
				},
				Country: searchrequest.Country{
					Code: countryCode,
				},
			},
		},
		Message: searchrequest.Message{
			Intent: searchrequest.Intent{
				Item: searchrequest.Item{
					Descriptor: searchrequest.ItemDescriptor{
						Name: payload.Role,
					},
				},
				Provider: searchrequest.Provider{
					Descriptor: searchrequest.ProviderDescriptor{
						Name: payload.Provider,
					},
					Locations: []searchrequest.ProviderLocations{
						{
							City: searchrequest.ProviderCity {
								Code: cityCode,
							},
							State: searchrequest.ProviderState {
								Code: stateCode,
							},
							Coordinates: searchrequest.Coordinates{
								Latitude:  payload.Location.Coordinates.Latitude,
								Longitude: payload.Location.Coordinates.Longitude,
							},
						},
					},
				},
				Tags: []searchrequest.Tags{
					{
						Descriptor: searchrequest.TagsDescriptor{
							Code: "JOB_DETAILS",
						},
						List: []searchrequest.List{
							{
								Descriptor: searchrequest.TagsDescriptor{
									Code: "INDUSTRY_TYPE",
								},
								Value: payload.Category,
							},
							{
								Descriptor: searchrequest.TagsDescriptor{
									Code: "JOB_TYPE",
								},
								Value: payload.EmploymentType,
							},
						},
					},
				},
			},
		},
	}
	return &req, nil
}
