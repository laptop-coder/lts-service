package handler

import (
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type InstitutionAdministratorHandler struct {
	institutionAdministratorService service.InstitutionAdministratorService
	log          logger.Logger
}

func NewInstitutionAdministratorHandler(institutionAdministratorService service.InstitutionAdministratorService, log logger.Logger) *InstitutionAdministratorHandler {
	return &InstitutionAdministratorHandler{
		institutionAdministratorService: institutionAdministratorService,
		log:          log,
	}
}

func (h *InstitutionAdministratorHandler) GetInstitutionAdministratorByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert institutionAdministrator ID
	institutionAdministratorID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert institutionAdministrator id to uuid", http.StatusBadRequest)
	}
	// Get institutionAdministrator
	response, err := h.institutionAdministratorService.GetInstitutionAdministratorByID(r.Context(), institutionAdministratorID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get institutionAdministrator by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"institutionAdministrator": response,
	})
}

func (h *InstitutionAdministratorHandler) GetOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. institutionAdministrator ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get institutionAdministrator
	response, err := h.institutionAdministratorService.GetInstitutionAdministratorByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get institutionAdministrator by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"institutionAdministrator": response,
	})
}
