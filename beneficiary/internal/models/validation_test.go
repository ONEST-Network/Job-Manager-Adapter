package models

import (
	"testing"
	"time"
)

func TestSchemeValidation(t *testing.T) {
	tests := []struct {
		name        string
		scheme      Scheme
		expectError bool
	}{
		{
			name: "Valid scheme",
			scheme: Scheme{
				ID:        "scheme1",
				Name:      "Test Scheme",
				Provider:  "gov",
				Amount:    1000,
				StartDate: time.Now().Format("2006-01-02"),
				Status:    "active",
			},
			expectError: false,
		},
		{
			name: "Invalid ID",
			scheme: Scheme{
				ID:        "",
				Name:      "Test Scheme",
				Provider:  "gov",
				StartDate: time.Now().Format("2006-01-02"),
			},
			expectError: true,
		},
		{
			name: "Invalid date format",
			scheme: Scheme{
				ID:        "scheme1",
				Name:      "Test Scheme",
				Provider:  "gov",
				StartDate: "invalid-date",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.scheme.Validate()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestApplicationValidation(t *testing.T) {
	tests := []struct {
		name        string
		application Application
		expectError bool
	}{
		{
			name: "Valid application",
			application: Application{
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
			name: "Missing credentials",
			application: Application{
				SchemeID:    "scheme1",
				ApplicantID: "user1",
				Credentials: nil,
			},
			expectError: true,
		},
		{
			name: "Missing required credential",
			application: Application{
				SchemeID:    "scheme1",
				ApplicantID: "user1",
				Credentials: map[string]interface{}{
					"aadhar": "1234",
					// Missing PAN
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.application.Validate()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
