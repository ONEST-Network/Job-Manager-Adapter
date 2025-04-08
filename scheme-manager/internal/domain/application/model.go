package application

import "time"

// ApplicantInfo contains information about the applicant
type ApplicantInfo struct {
	Name           string    `json:"name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Gender         string    `json:"gender"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Country        string    `json:"country"`
	ZipCode        string    `json:"zip_code"`
	EducationLevel string    `json:"education_level"`
	FamilyIncome   float64   `json:"family_income"`
}

// Document represents a submitted document
type Document struct {
	Type       string     `json:"type"`
	URL        string     `json:"url"`
	UploadedAt time.Time  `json:"uploaded_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	Status     string     `json:"status"` // Pending, Verified, Rejected
}

// Application represents an application for a scheme
type Application struct {
	ID                string        `json:"id"`
	SchemeID          string        `json:"scheme_id"`
	ApplicantID       string        `json:"applicant_id"`
	ApplicantInfo     ApplicantInfo `json:"applicant_info"`
	Documents         []Document    `json:"documents"`
	Status            string        `json:"status"`
	EligibilityStatus string        `json:"eligibility_status"`
	AppliedAt         time.Time     `json:"applied_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	ProcessedAt       *time.Time    `json:"processed_at,omitempty"`
}
