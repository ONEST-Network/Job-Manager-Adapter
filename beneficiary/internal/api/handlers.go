package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/VinVorteX/beneficiary-manager/internal/models"
	"github.com/VinVorteX/beneficiary-manager/internal/service"
)

type Handler struct {
	service *service.BeneficiaryService
}

func NewHandler(service *service.BeneficiaryService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSchemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filter := models.SchemeFilter{
		Provider:  r.URL.Query().Get("provider"),
		MinAmount: parseFloat(r.URL.Query().Get("min_amount")),
		MaxAmount: parseFloat(r.URL.Query().Get("max_amount")),
		Status:    r.URL.Query().Get("status"),
	}

	if err := filter.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	schemes, err := h.service.GetSchemes(filter)
	if err != nil {
		switch err {
		case service.ErrInvalidFilter:
			respondWithError(w, http.StatusBadRequest, "Invalid filter parameters")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch schemes")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, schemes)
}

func (h *Handler) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var app models.Application
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := app.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid application data: %v", err))
		return
	}

	if err := h.service.SubmitApplication(app); err != nil {
		switch err {
		case service.ErrInvalidApplication:
			respondWithError(w, http.StatusBadRequest, "Invalid application data")
		case service.ErrSchemeNotFound:
			respondWithError(w, http.StatusNotFound, "Scheme not found")
		case service.ErrSchemeInactive:
			respondWithError(w, http.StatusBadRequest, "Scheme is not active")
		case service.ErrSchemeExpired:
			respondWithError(w, http.StatusBadRequest, "Scheme has expired")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to submit application")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message":        "Application submitted successfully",
		"application_id": app.ID,
	})
}

func (h *Handler) GetApplicationStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	applicationID := r.URL.Query().Get("application_id")
	if applicationID == "" {
		respondWithError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	if !models.IsValidID(applicationID) {
		respondWithError(w, http.StatusBadRequest, "Invalid application ID format")
		return
	}

	status, err := h.service.GetApplicationStatus(applicationID)
	if err != nil {
		switch err {
		case service.ErrApplicationNotFound:
			respondWithError(w, http.StatusNotFound, "Application not found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch application status")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, status)
}

// Implement other handler methods...
