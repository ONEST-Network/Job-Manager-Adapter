package models

type Scheme struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Provider    string   `json:"provider"`
	Criteria    []string `json:"criteria"`
	Amount      float64  `json:"amount"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Status      string   `json:"status"`
}

type SchemeFilter struct {
	Provider  string  `query:"provider"`
	MinAmount float64 `query:"min_amount"`
	MaxAmount float64 `query:"max_amount"`
	Status    string  `query:"status"`
}
