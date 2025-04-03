package mock

import (
	"github.com/VinVorteX/beneficiary-manager/internal/models"
)

type MockDB struct {
	schemes      []models.Scheme
	applications map[string]models.Application
}

func NewMockDB() *MockDB {
	return &MockDB{
		applications: make(map[string]models.Application),
	}
}

func (m *MockDB) GetSchemes(filter models.SchemeFilter) ([]models.Scheme, error) {
	var filtered []models.Scheme
	for _, s := range m.schemes {
		if matchesFilter(s, filter) {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}

func (m *MockDB) GetSchemeByID(id string) (*models.Scheme, error) {
	for _, s := range m.schemes {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, nil
}

func (m *MockDB) SaveApplication(app models.Application) error {
	m.applications[app.ID] = app
	return nil
}

func (m *MockDB) GetApplicationStatus(id string) (*models.ApplicationStatus, error) {
	if app, exists := m.applications[id]; exists {
		return &models.ApplicationStatus{
			ApplicationID: app.ID,
			Status:       app.Status,
			UpdatedAt:    app.LastUpdatedAt,
		}, nil
	}
	return nil, nil
}

func (m *MockDB) AddScheme(scheme models.Scheme) {
	m.schemes = append(m.schemes, scheme)
}

func matchesFilter(s models.Scheme, f models.SchemeFilter) bool {
	if f.Provider != "" && f.Provider != s.Provider {
		return false
	}
	if f.MinAmount > 0 && s.Amount < f.MinAmount {
		return false
	}
	if f.MaxAmount > 0 && s.Amount > f.MaxAmount {
		return false
	}
	if f.Status != "" && f.Status != s.Status {
		return false
	}
	return true
}

// ConfigurePool implements the DB interface
func (m *MockDB) ConfigurePool(maxOpen, maxIdle int) error {
	return nil
} 