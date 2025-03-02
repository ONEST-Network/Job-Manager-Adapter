package onest

import (
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/business"
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"

	searchresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/response"
	searchresponseack "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/response-ack"

	selectresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/response"
	selectresponseack "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/response-ack"

	initresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/response"
	initresponseack "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/response-ack"

	confirmresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/response"
	confirmresponseack "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/response-ack"

	statusresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/response"
	statusresponseack "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/response-ack"

	cancelresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/response"
	cancelresponseack "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/response-ack"
)

type Interface interface {
	// search api handlers
	SearchJobsAck(body io.ReadCloser) (*searchresponse.SearchResponse, *searchresponseack.SearchResponseAck)
	SearchJobs(payload *searchresponse.SearchResponse)
	// select api handlers
	SendJobFulfillmentAck(body io.ReadCloser) (*selectresponse.SelectResponse, *selectresponseack.SelectResponseAck)
	SendJobFulfillment(payload *selectresponse.SelectResponse)
	// init api handlers
	InitializeJobApplicationAck(body io.ReadCloser) (*initresponse.InitResponse, *initresponseack.InitResponseAck)
	InitializeJobApplication(payload *initresponse.InitResponse)
	// confirm api handlers
	ConfirmJobApplicationAck(body io.ReadCloser) (*confirmresponse.ConfirmResponse, *confirmresponseack.ConfirmResponseAck)
	ConfirmJobApplication(response *confirmresponse.ConfirmResponse)
	// status api handlers
	JobApplicationStatusAck(body io.ReadCloser) (*statusresponse.StatusResponse, *statusresponseack.StatusResponseAck)
	JobApplicationStatus(payload *statusresponse.StatusResponse)
	// cancel api handlers
	WithdrawJobApplicationAck(body io.ReadCloser) (*cancelresponse.CancelResponse, *cancelresponseack.CancelResponseAck)
	WithdrawJobApplication(payload *cancelresponse.CancelResponse)
}

type Onest struct {
	clients *clients.Clients
}

func NewOnestClient(clients *clients.Clients) Interface {
	return &Onest{
		clients: clients,
	}
}

func (o *Onest) SearchJobsAck(body io.ReadCloser) (*searchresponse.SearchResponse, *searchresponseack.SearchResponseAck) {
    var (
        payload      searchresponse.SearchResponse
        payloadError *searchresponseack.Error
        status       = "ACK"
    )

    if err := json.NewDecoder(body).Decode(&payload); err != nil {
        payloadError = &searchresponseack.Error{
            Code:    "10000",
            Paths:   "",
            Message: err.Error(),
        }
    }

    if payloadError != nil {
        status = "NACK"
    }

    var errorResponse searchresponseack.Error
    if payloadError != nil {
        errorResponse = *payloadError
    }

    return &payload, &searchresponseack.SearchResponseAck{
        Message: searchresponseack.Message{
            Ack: searchresponseack.Ack{
                Status: status,
            },
        },
        Error: errorResponse,
    }
}

func (o *Onest) SearchJobs(payload *searchresponse.SearchResponse) {
    var jobs []job.Job
    
    // Extract jobs from each provider in the catalog
    for _, provider := range payload.Message.Catalog.Providers {
        for _, item := range provider.Items {
            // Extract work hours
            workHours := job.WorkHours{
                Start: "",
                End:   "",
            }
            // Extract work days
            workDays := job.WorkDays{
                Start: 0,
                End:   0,
            }
            
            // Process timing tags
            for _, tag := range item.Tags {
                if tag.Descriptor.Code == "TIMING" {
                    for _, list := range tag.List {
                        switch list.Descriptor.Code {
                        case "TIME_FROM":
                            workHours.Start = list.Value
                        case "TIME_TO":
                            workHours.End = list.Value
                        case "DAY_FROM":
                            start, _ := strconv.Atoi(list.Value)
                            workDays.Start = start
                        case "DAY_TO":
                            end, _ := strconv.Atoi(list.Value)
                            workDays.End = end
                        }
                    }
                }
            }

            // Extract salary range
            salaryRange := job.SalaryRange{
                Min: 0,
                Max: 0,
            }
            for _, tag := range item.Tags {
                if tag.Descriptor.Code == "SALARY_INFO" {
                    for _, list := range tag.List {
                        switch list.Descriptor.Code {
                        case "GROSS_MIN":
                            min, _ := strconv.Atoi(list.Value)
                            salaryRange.Min = min
                        case "GROSS_MAX":
                            max, _ := strconv.Atoi(list.Value)
                            salaryRange.Max = max
                        }
                    }
                }
            }

            // Extract eligibility criteria
            eligibility := job.Eligibility{
                AcademicQualification: job.AcademicQualificationNone,
                YearsOfExperience:    0,
                DocumentsRequired:     []job.Document{},
            }
            
            for _, tag := range item.Tags {
                switch tag.Descriptor.Code {
                case "ACADEMIC_ELIGIBILITY":
                    for _, list := range tag.List {
                        if list.Descriptor.Code == "COURSE_Level" {
                            eligibility.AcademicQualification = job.AcademicQualification(list.Value)
                        }
                    }
                case "JOB_REQUIREMENTS":
                    for _, list := range tag.List {
                        if list.Descriptor.Code == "REQ_EXPERIENCE" {
                            re := regexp.MustCompile(`P(\d+)Y`)
                            if matches := re.FindStringSubmatch(list.Value); len(matches) > 1 {
                                years, _ := strconv.Atoi(matches[1])
                                eligibility.YearsOfExperience = years
                            }
                        }
                    }
                case "DOCUMENT_NAME":
                    for _, list := range tag.List {
                        if list.Descriptor.Code == "DOCUMENT_NAME" {
                            eligibility.DocumentsRequired = append(eligibility.DocumentsRequired, job.Document(list.Value))
                        }
                    }
                }
            }

            // Extract location
            location := job.Location{}
            for _, loc := range provider.Locations {
                if contains(item.LocationIds, loc.ID) {
                    coords := strings.Split(loc.GPS, ",")
                    if len(coords) == 2 {
                        lat, _ := strconv.ParseFloat(coords[0], 64)
                        lng, _ := strconv.ParseFloat(coords[1], 64)
                        location = job.Location{
                            Address: loc.Address,
                            Street:  loc.Street,
                            City:    loc.City.Name,
                            State:   loc.State.Name,
                            Coordinates: job.Coordinates{
                                Type:        "Point",
                                Coordinates: []float64{lng, lat},
                            },
                        }
                        break
                    }
                }
            }

            // Create business details
            business := business.Business{
                Name:        item.Creator.Descriptor.Name,
                Description: item.Creator.Descriptor.LongDesc,
                Location: business.Location{
                    Address: item.Creator.Address,
                    City:    item.Creator.City.Name,
                    State:   item.Creator.State.Name,
                },
                Email: item.Creator.Contact.Email,
                Phone: item.Creator.Contact.Phone,
            }

            // Create job entry
            job := job.Job{
                ID:          item.ID,
                Name:        item.Descriptor.Name,
                Description: item.Descriptor.LongDesc,
                Type:        job.JobType("FULL_TIME"), // Set default or extract from tags if available
                Vacancies:   item.Quantity.Available.Count,
                SalaryRange: salaryRange,
                Business:    business,
                WorkHours:   workHours,
                WorkDays:    workDays,
                Eligibility: eligibility,
                Location:    location,
            }

            jobs = append(jobs, job)
        }
    }

    // Store jobs in database
    // if err := o.clients.JobClient.CreateJobs(jobs); err != nil {
    //     logrus.Errorf("Failed to store jobs: %v", err)
    //     return
    // }
}

func (j *Onest) SendJobFulfillmentAck(body io.ReadCloser) (*selectresponse.SelectResponse, *selectresponseack.SelectResponseAck) {
	var (
		payload  selectresponse.SelectResponse
		getError = func(message, paths, code string) *selectresponseack.SelectResponseAck {
			if code == "" {
				code = "10000"
			}

			return &selectresponseack.SelectResponseAck{
				Message: selectresponseack.Message{
					Ack: selectresponseack.Ack{
						Status: "NACK",
					},
				},
				Error: selectresponseack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if payload.Message.Order.Items == nil {
		return nil, getError("No items found", ".message.order.items", "30004")
	}

	jobs, err := j.clients.JobClient.ListJobs(bson.D{{Key: "id", Value: payload.Message.Order.Items[0].ID}})
	if err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if jobs == nil {
		return nil, getError("No job found for id: "+payload.Message.Order.Items[0].ID, "", "30004")
	}

	if jobs[0].Vacancies == 0 {
		return nil, getError("No vacancies available for job: "+jobs[0].ID, "", "40002")
	}

	return &payload, &selectresponseack.SelectResponseAck{
		Message: selectresponseack.Message{
			Ack: selectresponseack.Ack{
				Status: "ACK",
			},
		},
	}
}

func (j *Onest) SendJobFulfillment(payload *selectresponse.SelectResponse) {
	// Create a slice to store job IDs
    var jobIDs []string
    for _, item := range payload.Message.Order.Items {
        jobIDs = append(jobIDs, item.ID)
    }

    // Get jobs from database
    jobs, err := j.clients.JobClient.ListJobs(bson.D{{Key: "id", Value: bson.D{{Key: "$in", Value: jobIDs}}}})
    if err != nil {
        logrus.Errorf("Failed to fetch jobs: %v", err)
        return
    }
	// Log job details
    for _, job := range jobs {
        logrus.Infof("Job Details: ID=%s, Name=%s, Description=%s, Vacancies=%d", 
            job.ID, 
            job.Name, 
            job.Description, 
            job.Vacancies,
        )
        logrus.Infof("Salary Range: %d-%d", 
            job.SalaryRange.Min, 
            job.SalaryRange.Max,
        )
        logrus.Infof("Location: City=%s, State=%s, Address=%s", 
            job.Location.City, 
            job.Location.State, 
            job.Location.Address,
        )
        logrus.Infof("Business: Name=%s, Contact: Email=%s, Phone=%s",
            job.Business.Name,
            job.Business.Email,
            job.Business.Phone,
        )
        logrus.Infof("Work Schedule: Hours=%s-%s, Days=%d-%d",
            job.WorkHours.Start,
            job.WorkHours.End,
            job.WorkDays.Start,
            job.WorkDays.End,
        )
        logrus.Infof("Requirements: Academic=%s, Experience=%d years",
            job.Eligibility.AcademicQualification,
            job.Eligibility.YearsOfExperience,
        )
        if len(job.Eligibility.DocumentsRequired) > 0 {
            logrus.Infof("Required Documents: %v", job.Eligibility.DocumentsRequired)
        }
        logrus.Info("------------------------")
    }
}

func (j *Onest) InitializeJobApplicationAck(body io.ReadCloser) (*initresponse.InitResponse, *initresponseack.InitResponseAck) {
	var (
		payload  initresponse.InitResponse
		getError = func(message, paths, code string) *initresponseack.InitResponseAck {
			if code == "" {
				code = "10000"
			}

			return &initresponseack.InitResponseAck{
				Message: initresponseack.Message{
					Ack: initresponseack.Ack{
						Status: "NACK",
					},
				},
				Error: initresponseack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if payload.Message.Order.Fulfillments == nil {
		return nil, getError("no fulfillments found", ".message.order.fulfillments", "")
	}

	if payload.Message.Order.Items == nil {
		return nil, getError("no items found", ".message.order.items", "")
	}

    return &payload, &initresponseack.InitResponseAck{
        Message: initresponseack.Message{
            Ack: initresponseack.Ack{
                Status: "ACK",
            },
        },
        Error: initresponseack.Error{},
    }
}

func (j *Onest) InitializeJobApplication(payload *initresponse.InitResponse) {
	// Get applicant details from the first fulfillment
    fulfillment := payload.Message.Order.Fulfillments[0]
    
    // Pretty print application details
    logrus.Info("========== Job Application Details ==========")
    
    // Applicant Information
    logrus.Infof("Applicant Information:\n"+
        "Name: %s\n"+
        "Gender: %s\n"+
        "Age: %s\n"+
        "Contact: Phone=%s, Email=%s",
        fulfillment.Customer.Person.Name,
        fulfillment.Customer.Person.Gender,
        fulfillment.Customer.Person.Age,
        fulfillment.Customer.Contact.Phone,
        fulfillment.Customer.Contact.Email,
    )
	// Skills and Languages
    if len(fulfillment.Customer.Person.Skills) > 0 {
        var skills []string
        for _, skill := range fulfillment.Customer.Person.Skills {
            skills = append(skills, skill.Name)
        }
        logrus.Infof("Skills: %v", skills)
    }

    if len(fulfillment.Customer.Person.Languages) > 0 {
        var languages []string
        for _, lang := range fulfillment.Customer.Person.Languages {
            languages = append(languages, lang.Name)
        }
        logrus.Infof("Languages: %v", languages)
    }

    // Credentials/Documents
    if len(fulfillment.Customer.Person.Creds) > 0 {
        logrus.Info("Documents Submitted:")
        for _, cred := range fulfillment.Customer.Person.Creds {
            logrus.Infof("- Type: %s, URL: %s",
                cred.Type,
                cred.URL,
            )
        }
    }

	// Application Status
    logrus.Infof("Application Status:\n"+
        "State: %s\n"+
        "Updated At: %s",
        fulfillment.State.Descriptor.Code,
        fulfillment.State.UpdatedAt,
    )

    // Job Details
    logrus.Info("Applied Jobs:")
    for _, item := range payload.Message.Order.Items {
        logrus.Infof("Job ID: %s", item.ID)
    }

    // Provider Details
    logrus.Infof("Provider ID: %s", payload.Message.Order.Provider.ID)
    
    logrus.Info("===========================================")
}

func (j *Onest) ConfirmJobApplicationAck(body io.ReadCloser) (*confirmresponse.ConfirmResponse, *confirmresponseack.ConfirmResponseAck) {
	var (
		payload  confirmresponse.ConfirmResponse
		getError = func(message, paths, code string) *confirmresponseack.ConfirmResponseAck {
			if code == "" {
				code = "10000"
			}

			return &confirmresponseack.ConfirmResponseAck{
				Message: confirmresponseack.Message{
					Ack: confirmresponseack.Ack{
						Status: "NACK",
					},
				},
				Error: confirmresponseack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	if payload.Message.Order.Fulfillments == nil {
		logrus.Errorf("No fulfillments found")
		return nil, getError("no fulfillments found", ".message.order.fulfillments", "30004")
	}

	if payload.Message.Order.Items == nil {
		logrus.Errorf("No items found")
		return nil, getError("no items found", ".message.order.items", "30004")
	}

    return &payload, &confirmresponseack.ConfirmResponseAck{
        Message: confirmresponseack.Message{
            Ack: confirmresponseack.Ack{
                Status: "ACK",
            },
        },
        Error: confirmresponseack.Error{},
    }
}

func (j *Onest) ConfirmJobApplication(response *confirmresponse.ConfirmResponse) {
	logrus.Info("========== Job Application Confirmation ==========")
    
    // Transaction Details
    logrus.Infof("Transaction Details:\n"+
        "Transaction ID: %s\n"+
        "Order ID: %s\n"+
        "Provider ID: %s",
        response.Context.TransactionID,
        response.Message.Order.ID,
        response.Message.Order.Provider.ID,
    )
    
    // Application Status
    logrus.Infof("\nApplication Status:\n"+
        "State: %s\n"+
        "Updated At: %s",
        response.Message.Order.Fulfillments[0].State.Descriptor.Code,
        response.Message.Order.Fulfillments[0].State.UpdatedAt,
    )
	// Applicant Details
    fulfillment := response.Message.Order.Fulfillments[0]
    logrus.Infof("\nApplicant Information:\n"+
        "Name: %s\n"+
        "Gender: %s\n"+
        "Age: %s\n"+
        "Contact: Phone=%s, Email=%s",
        fulfillment.Customer.Person.Name,
        fulfillment.Customer.Person.Gender,
        fulfillment.Customer.Person.Age,
        fulfillment.Customer.Contact.Phone,
        fulfillment.Customer.Contact.Email,
    )

	// Skills
    if len(fulfillment.Customer.Person.Skills) > 0 {
        var skills []string
        for _, skill := range fulfillment.Customer.Person.Skills {
            skills = append(skills, skill.Name)
        }
        logrus.Infof("\nSkills: %v", skills)
    }
    
    // Languages
    if len(fulfillment.Customer.Person.Languages) > 0 {
        var languages []string
        for _, lang := range fulfillment.Customer.Person.Languages {
            languages = append(languages, lang.Name)
        }
        logrus.Infof("Languages: %v", languages)
    }
    
    // Documents
    if len(fulfillment.Customer.Person.Creds) > 0 {
        logrus.Info("\nSubmitted Documents:")
        for _, cred := range fulfillment.Customer.Person.Creds {
            logrus.Infof("- Type: %s, URL: %s",
                cred.Type,
                cred.URL,
            )
        }
    }

	// Applied Jobs
    logrus.Info("\nApplied Jobs:")
    for _, item := range response.Message.Order.Items {
        logrus.Infof("Job ID: %s", item.ID)
    }
    
    logrus.Info("============================================")
}

func (j *Onest) JobApplicationStatusAck(body io.ReadCloser) (*statusresponse.StatusResponse, *statusresponseack.StatusResponseAck) {
	var (
		payload  statusresponse.StatusResponse
		getError = func(message, paths, code string) *statusresponseack.StatusResponseAck {
			if code == "" {
				code = "10000"
			}

			return &statusresponseack.StatusResponseAck{
				Message: statusresponseack.Message{
					Ack: statusresponseack.Ack{
						Status: "NACK",
					},
				},
				Error: statusresponseack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	return &payload, &statusresponseack.StatusResponseAck{
		Message: statusresponseack.Message{
			Ack: statusresponseack.Ack{
				Status: "ACK",
			},
		},
		Error: statusresponseack.Error{},
	}
}

func (j *Onest) JobApplicationStatus(payload *statusresponse.StatusResponse) {
	logrus.Info("========== Job Application Status ==========")
    // Transaction Details
    logrus.Infof("Transaction Details:\n"+
        "Transaction ID: %s\n"+
        "Message ID: %s\n"+
        "Domain: %s\n"+
        "Action: %s\n"+
        "Timestamp: %s",
        payload.Context.TransactionID,
        payload.Context.MessageID,
        payload.Context.Domain,
        payload.Context.Action,
        payload.Context.Timestamp,
    )

    // Order Details
    order := payload.Message.Order
    logrus.Infof("\nOrder Information:\n"+
        "Order ID: %s\n"+
        "Provider ID: %s",
        order.ID,
        order.Provider.ID,
    )

    // Job Details
    if len(order.Items) > 0 {
        logrus.Info("\nJob Details:")
        for _, item := range order.Items {
            logrus.Infof("Job ID: %s", item.ID)
            if len(item.Tags) > 0 {
                for _, tag := range item.Tags {
                    logrus.Infof("Tag Code: %s", tag.Descriptor.Code)
                    for _, list := range tag.List {
                        logrus.Infof("- %s: %s", list.Code, list.Value)
                    }
                }
            }
        }
    }

    // Fulfillment Status
    if len(order.Fulfillments) > 0 {
        fulfillment := order.Fulfillments[0]
        logrus.Infof("\nApplication Status:\n"+
            "ID: %s\n"+
            "Type: %s\n"+
            "State: %s",
            fulfillment.ID,
            fulfillment.Type,
            fulfillment.State.Descriptor.Code,
        )
    }

	logrus.Info("===========================================")
}

func (j *Onest) WithdrawJobApplicationAck(body io.ReadCloser) (*cancelresponse.CancelResponse, *cancelresponseack.CancelResponseAck) {
	var (
		payload  cancelresponse.CancelResponse
		getError = func(message, paths, code string) *cancelresponseack.CancelResponseAck {
			if code == "" {
				code = "10000"
			}

			return &cancelresponseack.CancelResponseAck{
				Message: cancelresponseack.Message{
					Ack: cancelresponseack.Ack{
						Status: "NACK",
					},
				},
				Error: cancelresponseack.Error{
					Code:    code,
					Paths:   paths,
					Message: message,
				},
			}
		}
	)

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, getError(err.Error(), "", "")
	}

	return &payload, &cancelresponseack.CancelResponseAck{
		Message: cancelresponseack.Message{
			Ack: cancelresponseack.Ack{
				Status: "ACK",
			},
		},
		Error: cancelresponseack.Error{},
	}
}

func (j *Onest) WithdrawJobApplication(payload *cancelresponse.CancelResponse) {
	logrus.Info("========== Job Application Withdrawal ==========")
    
    // Transaction Details
    logrus.Infof("Transaction Details:\n"+
        "Transaction ID: %s\n"+
        "Message ID: %s\n"+
        "Domain: %s\n"+
        "Action: %s\n"+
        "Timestamp: %s",
        payload.Context.TransactionID,
        payload.Context.MessageID,
        payload.Context.Domain,
        payload.Context.Action,
        payload.Context.Timestamp,
    )

    // Order Details
    order := payload.Message.Order
    logrus.Infof("\nOrder Information:\n"+
        "Order ID: %s\n"+
        "Status: %s\n"+
        "Provider ID: %s",
        order.ID,
        order.Status,
        order.Provider.ID,
    )

	// Job Details
    if len(order.Items) > 0 {
        item := order.Items[0]
        logrus.Infof("\nJob Details:\n"+
            "Job ID: %s\n"+
            "Application Period: %s to %s",
            item.ID,
            item.Time.Range.Start,
            item.Time.Range.End,
        )
    }

    // Fulfillment Status
    if len(order.Fulfillments) > 0 {
        for _, fulfillment := range order.Fulfillments {
            if fulfillment.Type != "" {
                logrus.Infof("\nFulfillment Details:\n"+
                    "ID: %s\n"+
                    "Type: %s\n"+
                    "State: %s",
                    fulfillment.ID,
                    fulfillment.Type,
                    fulfillment.State.Descriptor.Code,
                )
            }
        }
    }

    logrus.Info("===========================================")
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}

