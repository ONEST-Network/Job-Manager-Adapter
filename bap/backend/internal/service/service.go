package service

import (
	"context"
	"fmt"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
	cancelrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/request"
	cancelresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/cancel/response"
	confirmrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/request"
	confirmresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/confirm/response"
	initrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/request"
	initresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/init/response"
	searchrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/request"
	searchresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/response"
	selectrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/request"
	selectresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/select/response"
	statusrequest "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/request"
	statusresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/status/response"
)

type OnestService struct {
    Clients *clients.Clients
}

func NewOnestService(clients *clients.Clients) *OnestService {
    return &OnestService{
        Clients: clients,
    }
}

func (s *OnestService) Search(ctx context.Context, req *searchrequest.SearchRequest) (*searchresponse.SearchResponse, error) {
    // Make API call to BPP search endpoint
    var response searchresponse.SearchResponse
    err := s.Clients.ApiClient.ApiCall(req, "/search", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to search jobs: %w", err)
    }
    return &response, nil
}

func (s *OnestService) Select(ctx context.Context, req *selectrequest.SelectRequest) (*selectresponse.SelectResponse, error) {
    // Make API call to BPP select endpoint
    var response selectresponse.SelectResponse
    err := s.Clients.ApiClient.ApiCall(req, "/select", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to select job: %w", err)
    }
    return &response, nil
}

func (s *OnestService) Init(ctx context.Context, req *initrequest.InitRequest) (*initresponse.InitResponse, error) {
    // Make API call to BPP init endpoint
    var response initresponse.InitResponse
    err := s.Clients.ApiClient.ApiCall(req, "/init", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to initialize job application: %w", err)
    }
    return &response, nil
}

func (s *OnestService) Confirm(ctx context.Context, req *confirmrequest.ConfirmRequest) (*confirmresponse.ConfirmResponse, error) {
    // Make API call to BPP confirm endpoint
    var response confirmresponse.ConfirmResponse
    err := s.Clients.ApiClient.ApiCall(req, "/confirm", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to confirm job application: %w", err)
    }
    return &response, nil
}

func (s *OnestService) Status(ctx context.Context, req *statusrequest.StatusRequest) (*statusresponse.StatusResponse, error) {
    // Make API call to BPP status endpoint
    var response statusresponse.StatusResponse
    err := s.Clients.ApiClient.ApiCall(req, "/status", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to get job application status: %w", err)
    }
    return &response, nil
}

func (s *OnestService) Cancel(ctx context.Context, req *cancelrequest.CancelRequest) (*cancelresponse.CancelResponse, error) {
    // Make API call to BPP cancel endpoint
    var response cancelresponse.CancelResponse
    err := s.Clients.ApiClient.ApiCall(req, "/cancel", &response, "POST")
    if err != nil {
        return nil, fmt.Errorf("failed to cancel job application: %w", err)
    }
    return &response, nil
}

