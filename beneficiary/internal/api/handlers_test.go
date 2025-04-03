package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/VinVorteX/beneficiary-manager/internal/db/mock"
	loggerMock "github.com/VinVorteX/beneficiary-manager/internal/logger/mock"
	"github.com/VinVorteX/beneficiary-manager/internal/models"
	"github.com/VinVorteX/beneficiary-manager/internal/service"
)

func setupTestHandler() *Handler {
	mockDB := mock.NewMockDB()
	mockLogger := loggerMock.NewMockLogger()
	beneficiaryService := service.NewBeneficiaryService(mockDB, mockLogger)
	return NewHandler(beneficiaryService)
}

func TestGetSchemes(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			queryParams:    "?provider=gov&min_amount=1000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid amount filter",
			queryParams:    "?min_amount=-100",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/schemes"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.GetSchemes(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSubmitApplication(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		application    models.Application
		expectedStatus int
	}{
		{
			name: "Valid application",
			application: models.Application{
				SchemeID:    "scheme1",
				ApplicantID: "user1",
				Credentials: map[string]interface{}{
					"aadhar": "1234",
					"pan":    "ABCD1234",
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid application",
			application: models.Application{
				SchemeID: "", // Missing required field
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.application)
			req := httptest.NewRequest("POST", "/applications", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.SubmitApplication(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
