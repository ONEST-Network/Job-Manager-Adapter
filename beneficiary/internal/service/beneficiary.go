package service

import (
	"fmt"
	"time"
	"github.com/VinVorteX/beneficiary-manager/internal/db"
	"github.com/VinVorteX/beneficiary-manager/internal/logger"
	"github.com/VinVorteX/beneficiary-manager/internal/models"
)

type BeneficiaryService struct {
	db     db.DB
	logger logger.Logger
}

func NewBeneficiaryService(db db.DB, logger logger.Logger) *BeneficiaryService {
	return &BeneficiaryService{
		db:     db,
		logger: logger,
	}
}

func (s *BeneficiaryService) GetSchemes(filter models.SchemeFilter) ([]models.Scheme, error) {
	s.logger.Debug("Getting schemes with filter", map[string]interface{}{
		"provider":   filter.Provider,
		"min_amount": filter.MinAmount,
		"max_amount": filter.MaxAmount,
		"status":     filter.Status,
	})

	if err := filter.Validate(); err != nil {
		s.logger.Error("Invalid filter parameters", err, nil)
		return nil, ErrInvalidFilter
	}

	schemes, err := s.db.GetSchemes(filter)
	if err != nil {
		s.logger.Error("Failed to fetch schemes", err, nil)
		return nil, err
	}

	for _, scheme := range schemes {
		if err := scheme.Validate(); err != nil {
			s.logger.Error("Invalid scheme data in database", err, map[string]interface{}{
				"scheme_id": scheme.ID,
			})
			return nil, fmt.Errorf("invalid scheme data in database: %v", err)
		}
	}

	s.logger.Info("Successfully fetched schemes", map[string]interface{}{
		"count": len(schemes),
	})
	return schemes, nil
}

func (s *BeneficiaryService) SubmitApplication(app models.Application) error {
	s.logger.Debug("Submitting application", map[string]interface{}{
		"scheme_id":    app.SchemeID,
		"applicant_id": app.ApplicantID,
	})

	if err := app.Validate(); err != nil {
		s.logger.Error("Invalid application data", err, nil)
		return ErrInvalidApplication
	}

	scheme, err := s.db.GetSchemeByID(app.SchemeID)
	if err != nil {
		s.logger.Error("Failed to fetch scheme", err, map[string]interface{}{
			"scheme_id": app.SchemeID,
		})
		return err
	}
	if scheme == nil {
		s.logger.Warn("Scheme not found", map[string]interface{}{
			"scheme_id": app.SchemeID,
		})
		return ErrSchemeNotFound
	}

	currentDate := time.Now().Format("2006-01-02")
	if scheme.Status != "active" {
		s.logger.Warn("Attempt to apply for inactive scheme", map[string]interface{}{
			"scheme_id": app.SchemeID,
			"status":    scheme.Status,
		})
		return ErrSchemeInactive
	}

	if scheme.EndDate != "" && scheme.EndDate < currentDate {
		s.logger.Warn("Attempt to apply for expired scheme", map[string]interface{}{
			"scheme_id": app.SchemeID,
			"end_date":  scheme.EndDate,
		})
		return ErrSchemeExpired
	}

	// Set metadata
	app.Status = "pending"
	app.SubmittedAt = time.Now().Format(time.RFC3339)
	app.LastUpdatedAt = app.SubmittedAt

	if err := s.db.SaveApplication(app); err != nil {
		s.logger.Error("Failed to save application", err, map[string]interface{}{
			"scheme_id":    app.SchemeID,
			"applicant_id": app.ApplicantID,
		})
		return err
	}

	s.logger.Info("Application submitted successfully", map[string]interface{}{
		"application_id": app.ID,
		"scheme_id":      app.SchemeID,
		"applicant_id":   app.ApplicantID,
	})
	return nil
}

func (s *BeneficiaryService) GetApplicationStatus(applicationID string) (models.ApplicationStatus, error) {
	status, err := s.db.GetApplicationStatus(applicationID)
	if err != nil {
		return models.ApplicationStatus{}, err
	}
	if status == nil {
		return models.ApplicationStatus{}, ErrApplicationNotFound
	}
	return *status, nil
}
