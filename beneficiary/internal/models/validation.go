package models

import (
	"fmt"
	"regexp"
	"time"
)

var (
	// Regex for basic ID format validation
	idRegex = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
)

// ValidateSchemeFilter validates the scheme filter parameters
func (f *SchemeFilter) Validate() error {
	if f.MinAmount < 0 {
		return fmt.Errorf("min_amount cannot be negative")
	}
	if f.MaxAmount < 0 {
		return fmt.Errorf("max_amount cannot be negative")
	}
	if f.MaxAmount > 0 && f.MinAmount > f.MaxAmount {
		return fmt.Errorf("min_amount cannot be greater than max_amount")
	}
	if f.Status != "" && !isValidStatus(f.Status) {
		return fmt.Errorf("invalid status value: %s", f.Status)
	}
	return nil
}

// ValidateApplication validates the application data
func (a *Application) Validate() error {
	if !idRegex.MatchString(a.SchemeID) {
		return fmt.Errorf("invalid scheme_id format")
	}
	if !idRegex.MatchString(a.ApplicantID) {
		return fmt.Errorf("invalid applicant_id format")
	}
	if a.Credentials == nil {
		return fmt.Errorf("credentials cannot be nil")
	}
	if len(a.Credentials) == 0 {
		return fmt.Errorf("credentials cannot be empty")
	}
	
	// Validate required credentials
	if err := validateCredentials(a.Credentials); err != nil {
		return fmt.Errorf("invalid credentials: %v", err)
	}
	
	return nil
}

// ValidateScheme validates the scheme data
func (s *Scheme) Validate() error {
	if !idRegex.MatchString(s.ID) {
		return fmt.Errorf("invalid id format")
	}
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if s.Amount < 0 {
		return fmt.Errorf("amount cannot be negative")
	}
	
	// Validate dates
	if err := validateDates(s.StartDate, s.EndDate); err != nil {
		return err
	}
	
	return nil
}

// Helper functions

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"active":    true,
		"inactive":  true,
		"pending":   true,
		"approved":  true,
		"rejected":  true,
		"completed": true,
	}
	return validStatuses[status]
}

func validateCredentials(creds map[string]interface{}) error {
	requiredFields := []string{"aadhar", "pan"} // Add required fields as needed
	
	for _, field := range requiredFields {
		if _, exists := creds[field]; !exists {
			return fmt.Errorf("missing required credential: %s", field)
		}
	}
	return nil
}

func validateDates(start, end string) error {
	if start == "" {
		return fmt.Errorf("start_date is required")
	}
	
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("invalid start_date format: use YYYY-MM-DD")
	}
	
	if end != "" {
		endDate, err := time.Parse("2006-01-02", end)
		if err != nil {
			return fmt.Errorf("invalid end_date format: use YYYY-MM-DD")
		}
		if endDate.Before(startDate) {
			return fmt.Errorf("end_date cannot be before start_date")
		}
	}
	
	return nil
}

// IsValidID checks if the provided ID matches the required format
func IsValidID(id string) bool {
	return idRegex.MatchString(id)
} 