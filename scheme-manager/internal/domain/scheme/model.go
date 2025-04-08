package scheme

import "time"

// EligibilityCriteria defines the criteria for a scheme
type EligibilityCriteria struct {
	MinAge             *int     `json:"min_age,omitempty"`
	MaxAge             *int     `json:"max_age,omitempty"`
	Gender             []string `json:"gender,omitempty"`
	Location           []string `json:"location,omitempty"`
	EducationLevel     []string `json:"education_level,omitempty"`
	MaxFamilyIncome    *float64 `json:"max_family_income,omitempty"`
	RequiredDocuments  []string `json:"required_documents,omitempty"`
	AdditionalCriteria []string `json:"additional_criteria,omitempty"`
}

// Scheme represents a financial support or scholarship scheme
type Scheme struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	OrganizationID      string              `json:"organization_id"`
	Amount              float64             `json:"amount"`
	Currency            string              `json:"currency"`
	StartDate           time.Time           `json:"start_date"`
	EndDate             time.Time           `json:"end_date"`
	ApplicationDeadline time.Time           `json:"application_deadline"`
	MaxBeneficiaries    *int                `json:"max_beneficiaries,omitempty"`
	EligibilityCriteria EligibilityCriteria `json:"eligibility_criteria"`
	Status              string              `json:"status"` // Active, Inactive, Draft, Closed
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}
