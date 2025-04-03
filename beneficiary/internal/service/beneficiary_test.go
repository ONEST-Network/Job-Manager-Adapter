package service

import (
	"testing"
	"time"
	"github.com/VinVorteX/beneficiary-manager/internal/db/mock"
	loggerMock "github.com/VinVorteX/beneficiary-manager/internal/logger/mock"
	"github.com/VinVorteX/beneficiary-manager/internal/models"
)

func TestGetSchemes(t *testing.T) {
	mockDB := mock.NewMockDB()
	mockLogger := loggerMock.NewMockLogger()
	service := NewBeneficiaryService(mockDB, mockLogger)

	// Add test schemes
	mockDB.AddScheme(models.Scheme{
		ID:        "scheme1",
		Name:      "Test Scheme 1",
		Provider:  "gov",
		Amount:    1000,
		Status:    "active",
		StartDate: time.Now().Format("2006-01-02"),
	})

	tests := []struct {
		name          string
		filter        models.SchemeFilter
		expectedCount int
		expectError   bool
	}{
		{
			name:          "No filter",
			filter:        models.SchemeFilter{},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "Invalid amount filter",
			filter: models.SchemeFilter{
				MinAmount: 2000,
				MaxAmount: 1000,
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemes, err := service.GetSchemes(tt.filter)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(schemes) != tt.expectedCount {
				t.Errorf("expected %d schemes, got %d", tt.expectedCount, len(schemes))
			}
		})
	}
}

func TestSubmitApplication(t *testing.T) {
	mockDB := mock.NewMockDB()
	mockLogger := loggerMock.NewMockLogger()
	service := NewBeneficiaryService(mockDB, mockLogger)

	// Add active scheme
	mockDB.AddScheme(models.Scheme{
		ID:        "scheme1",
		Status:    "active",
		StartDate: time.Now().Format("2006-01-02"),
	})

	tests := []struct {
		name        string
		application models.Application
		expectError bool
	}{
		{
			name: "Valid application",
			application: models.Application{
				ID:          "app1",
				SchemeID:    "scheme1",
				ApplicantID: "user1",
				Credentials: map[string]interface{}{
					"aadhar": "1234",
					"pan":    "ABCD1234",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid scheme ID",
			application: models.Application{
				SchemeID:    "nonexistent",
				ApplicantID: "user1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SubmitApplication(tt.application)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
