package models

type Application struct {
	ID            string                 `json:"id"`
	SchemeID      string                 `json:"scheme_id"`
	ApplicantID   string                 `json:"applicant_id"`
	Status        string                 `json:"status"`
	Credentials   map[string]interface{} `json:"credentials"`
	SubmittedAt   string                 `json:"submitted_at"`
	LastUpdatedAt string                 `json:"last_updated_at"`
}

type ApplicationStatus struct {
	ApplicationID string `json:"application_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	UpdatedAt     string `json:"updated_at"`
}
