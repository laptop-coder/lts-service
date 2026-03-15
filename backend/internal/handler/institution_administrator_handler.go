package handler

import (
	"strconv"
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
	log                             logger.Logger
}

func NewInstitutionAdministratorHandler(institutionAdministratorService service.InstitutionAdministratorService, log logger.Logger) *InstitutionAdministratorHandler {
	return &InstitutionAdministratorHandler{
		institutionAdministratorService: institutionAdministratorService,
		log:                             log,
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

func (h *InstitutionAdministratorHandler) AssignPosition(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert institution administrator ID (user ID)
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert institution administrator id to uuid", http.StatusBadRequest)
	}
	// Get and convert institution administrator position ID:
	positionIDFields := r.PostForm["positionId"]
	if len(positionIDFields) != 1 {
		helpers.ErrorResponse(w, "failed to parse form: positionId must be provided exactly once", http.StatusBadRequest)
	}
	// convert to uint64
	positionID64, err := strconv.ParseUint(positionIDFields[0], 10, 8)
	if err != nil {
		h.log.Error("cannot convert position ID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert position ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// and to uint8
	positionID := uint8(positionID64)
	// Assign position to institution administrator
	if err := h.institutionAdministratorService.AssignPosition(r.Context(), userID, positionID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to assign position to institution administrator: %w", err))
		return
	}
	// Get updated institution administrator
	institutionAdministrator, err := h.institutionAdministratorService.GetInstitutionAdministratorByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get institution administrator by institution administrator ID (user ID): %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"institutionAdministrator": institutionAdministrator,
	})
}

func (h *InstitutionAdministratorHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert institution administrator ID (user ID)
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert institution administrator id to uuid", http.StatusBadRequest)
	}
	// Get institution administrator position
	response, err := h.institutionAdministratorService.GetPosition(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get institution administrator position by institution administrator id (user id): %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"institutionAdministratorPosition": response,
	})
}
