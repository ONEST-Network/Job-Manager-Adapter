package db

import (
	"github.com/VinVorteX/beneficiary-manager/internal/models"
)

// DB defines the interface for database operations
type DB interface {
	GetSchemes(filter models.SchemeFilter) ([]models.Scheme, error)
	GetSchemeByID(id string) (*models.Scheme, error)
	SaveApplication(app models.Application) error
	GetApplicationStatus(applicationID string) (*models.ApplicationStatus, error)
	ConfigurePool(maxOpen, maxIdle int) error
}
