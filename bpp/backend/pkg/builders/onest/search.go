package onest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/config"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/database/mongodb/job"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/search/request"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/types/payload/search/response"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bpp/backend/pkg/utils"
)

func BuildSearchJobsResponse(clients *clients.Clients, payload *request.SearchRequest, jobs []job.Job) (*response.SearchResponse, error) {
	res := response.SearchResponse{
		Context: response.Context{
			Domain:        payload.Context.Domain,
			Action:        "on_search",
			Version:       payload.Context.Version,
			BapID:         payload.Context.BapID,
			BapURI:        payload.Context.BapURI,
			BppID:         config.Config.BppId,
			BppURI:        config.Config.BppUri,
			TransactionID: payload.Context.TransactionID,
			MessageID:     payload.Context.MessageID,
			Location: response.Location{
				City: response.ContextCity{
					Code: payload.Context.Location.City.Code,
				},
				Country: response.Country{
					Code: payload.Context.Location.Country.Code,
				},
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			TTL:       "PT30S",
		},
		Message: response.Message{
			Catalog: response.Catalog{
				Descriptor: response.CatalogDescriptor{
					Name: "BPP",
				},
				Providers: []response.Providers{
					{
						ID: "1",
						Descriptor: response.ProvidersDescriptor{
							Name:      "BPP",
							ShortDesc: "BPP",
						},
						Fulfillments: []response.Fulfillments{
							{
								ID:   "F1",
								Type: "lead & recruitment",
							},
						},
						Locations: []response.Locations{},
						Items:     []response.Items{},
					},
				},
			},
		},
	}

	for i, job := range jobs {
		business, err := clients.BusinessClient.GetBusiness(job.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s business, %v", job.BusinessID, err)
		}

		businessStateCode, err := utils.GetStateCode(business.Location.State)
		if err != nil {
			return nil, fmt.Errorf("failed to get business state code for %s, %v", business.Location.State, err)
		}

		businessCityCode, err := utils.GetCityCode(business.Location.City)
		if err != nil {
			return nil, fmt.Errorf("failed to get business city code for %s, %v", business.Location.State, err)
		}

		jobStateCode, err := utils.GetStateCode(job.Location.State)
		if err != nil {
			return nil, fmt.Errorf("failed to get job state code for %s, %v", job.Location.State, err)
		}

		jobCityCode, err := utils.GetCityCode(job.Location.City)
		if err != nil {
			return nil, fmt.Errorf("failed to get job city code for %s, %v", job.Location.State, err)
		}

		res.Message.Catalog.Providers[0].Locations = append(res.Message.Catalog.Providers[0].Locations, response.Locations{
			ID:      fmt.Sprintf("L%d", i+1),
			Address: job.Location.Address,
			Street:  job.Location.Street,
			City: response.City{
				Name: job.Location.City,
				Code: jobCityCode,
			},
			State: response.State{
				Name: job.Location.State,
				Code: jobStateCode,
			},
			GPS: fmt.Sprintf("%f,%f", job.Location.Coordinates.Latitude, job.Location.Coordinates.Longitute),
		})

		item := response.Items{
			ID: job.ID,
			Descriptor: response.ItemsDescriptor{
				Name:     job.Name,
				LongDesc: job.Description,
			},
			Quantity: response.Quantity{
				Available: response.Available{
					Count: job.Vacancies,
				},
			},
			LocationIds: []string{fmt.Sprintf("L%d", i+1)},
			FulfillmentIds: []string{
				"F1",
			},
			Creator: response.Creator{
				Descriptor: response.CreatorDescriptor{
					Name:     business.Name,
					LongDesc: business.Description,
				},
				Address: business.Location.Address,
				State: response.State{
					Name: business.Location.State,
					Code: businessStateCode,
				},
				City: response.City{
					Name: business.Location.City,
					Code: businessCityCode,
				},
				Contact: response.Contact{
					Email: business.Email,
					Phone: business.Phone,
				},
			},
			Tags: []response.Tags{
				*addTimingTag(&job),
				*addSalaryRange(&job),
			},
		}

		if academicEligibility := addAcademicEligibility(&job); academicEligibility != nil {
			item.Tags = append(item.Tags, *academicEligibility)
		}

		for _, doc := range addDocumentsRequired(&job) {
			item.Tags = append(item.Tags, *doc)
		}

		if jobRequirements := addJobRequirements(&job); jobRequirements != nil {
			item.Tags = append(item.Tags, *jobRequirements)
		}

		res.Message.Catalog.Providers[0].Items = append(res.Message.Catalog.Providers[0].Items, item)
	}

	return &res, nil
}

func addTimingTag(j *job.Job) *response.Tags {
	return &response.Tags{
		Descriptor: response.TagsDescriptor{
			Code: "TIMING",
			Name: "Timing",
		},
		List: []response.List{
			{
				Descriptor: response.TagsDescriptor{
					Code: "DAY_FROM",
					Name: "Day from",
				},
				Value: strconv.Itoa(j.WorkDays.Start),
			},
			{
				Descriptor: response.TagsDescriptor{
					Code: "DAY_TO",
					Name: "Day to",
				},
				Value: strconv.Itoa(j.WorkDays.End),
			},
			{
				Descriptor: response.TagsDescriptor{
					Code: "TIME_FROM",
					Name: "Time from",
				},
				Value: j.WorkHours.Start,
			},
			{
				Descriptor: response.TagsDescriptor{
					Code: "TIME_TO",
					Name: "Time to",
				},
				Value: j.WorkHours.Start,
			},
		},
	}
}

func addAcademicEligibility(j *job.Job) *response.Tags {
	if j.Eligibility.AcademicQualification == job.AcademicQualificationNone {
		return nil
	}

	return &response.Tags{
		Descriptor: response.TagsDescriptor{
			Code: "ACADEMIC_ELIGIBILITY",
			Name: "Academic Eligibility",
		},
		List: []response.List{
			{
				Descriptor: response.TagsDescriptor{
					Code: "COURSE_Level",
					Name: "Level of the course",
				},
				Value: string(j.Eligibility.AcademicQualification),
			},
			{
				Descriptor: response.TagsDescriptor{
					Code: "MANDATORY_ELIGIBILITY",
					Name: "Name of the course",
				},
				Value: "true",
			},
		},
	}
}

func addDocumentsRequired(j *job.Job) []*response.Tags {
	if j.Eligibility.DocumentsRequired == nil {
		return nil
	}

	var tags []*response.Tags

	for _, doc := range j.Eligibility.DocumentsRequired {
		tags = append(tags, &response.Tags{
			Descriptor: response.TagsDescriptor{
				Code: "DOCUMENT_NAME",
				Name: "Name of the document",
			},
			List: []response.List{
				{
					Descriptor: response.TagsDescriptor{
						Code: "DOCUMENT_NAME",
						Name: "Name of the document",
					},
					Value: string(doc),
				},
				{
					Descriptor: response.TagsDescriptor{
						Code: "MANDATORY_DOCUMENT",
						Name: "Mandatory DOCUMENT",
					},
					Value: "true",
				},
			},
		})
	}

	return tags
}

func addJobRequirements(j *job.Job) *response.Tags {
	if j.Eligibility.YearsOfExperience == 0 {
		return nil
	}

	return &response.Tags{
		Descriptor: response.TagsDescriptor{
			Code: "JOB_REQUIREMENTS",
			Name: "Job requirements",
		},
		List: []response.List{
			{
				Descriptor: response.TagsDescriptor{
					Code: "REQ_EXPERIENCE",
					Name: "Required work experience in years",
				},
				Value: fmt.Sprintf("P%dY", j.Eligibility.YearsOfExperience),
			},
		},
	}
}

func addSalaryRange(j *job.Job) *response.Tags {
	return &response.Tags{
		Descriptor: response.TagsDescriptor{
			Code: "SALARY_INFO",
			Name: "Salary information",
		},
		List: []response.List{
			{
				Descriptor: response.TagsDescriptor{
					Code: "GROSS_MIN",
					Name: "Minimum gross pay",
				},
				Value: strconv.Itoa(j.SalaryRange.Min),
			},
			{
				Descriptor: response.TagsDescriptor{
					Code: "GROSS_MAX",
					Name: "Maximum gross pay",
				},
				Value: strconv.Itoa(j.SalaryRange.Max),
			},
		},
	}
}
