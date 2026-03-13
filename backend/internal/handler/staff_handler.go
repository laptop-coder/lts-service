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

type StaffHandler struct {
	staffService service.StaffService
	log          logger.Logger
}

func NewStaffHandler(staffService service.StaffService, log logger.Logger) *StaffHandler {
	return &StaffHandler{
		staffService: staffService,
		log:          log,
	}
}

func (h *StaffHandler) GetStaffByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert staff ID
	staffID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert staff id to uuid", http.StatusBadRequest)
	}
	// Get staff
	response, err := h.staffService.GetStaffByID(r.Context(), staffID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get staff by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staff": response,
	})
}

func (h *StaffHandler) GetOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. staff ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get staff
	response, err := h.staffService.GetStaffByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get staff by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staff": response,
	})
}
