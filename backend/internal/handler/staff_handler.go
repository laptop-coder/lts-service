package handler

import (
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
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
		return
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

func (h *StaffHandler) AssignPosition(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert staff ID (user ID)
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert staff id to uuid", http.StatusBadRequest)
		return
	}
	// Get and convert staff position ID:
	positionIDFields := r.PostForm["positionId"]
	if len(positionIDFields) != 1 {
		helpers.ErrorResponse(w, "failed to parse form: positionId must be provided exactly once", http.StatusBadRequest)
		return
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
	// Assign position to staff
	if err := h.staffService.AssignPosition(r.Context(), userID, positionID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to assign position to staff: %w", err))
		return
	}
	// Get updated staff
	staff, err := h.staffService.GetStaffByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get staff by staff ID (user ID): %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staff": staff,
	})
}

func (h *StaffHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert staff ID (user ID)
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert staff id to uuid", http.StatusBadRequest)
		return
	}
	// Get staff position
	response, err := h.staffService.GetPosition(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get staff position by staff id (user id): %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staffPosition": response,
	})
}
