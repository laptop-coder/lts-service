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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert staff ID
	staffID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert staff id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get staff
	response, err := h.staffService.GetStaffByID(r.Context(), staffID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get staff by id: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. staff ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get staff
	response, err := h.staffService.GetStaffByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get staff by id: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error("failed to parse x-www-form-urlencoded form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get and convert staff ID (user ID)
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert staff id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get and convert staff position ID:
	positionIDFields := r.PostForm["positionId"]
	if len(positionIDFields) != 1 {
		h.log.Error("failed to parse form: positionId must be provided exactly once")
		helpers.FieldExactlyOneError(h.log, w, "positionId")
		return
	}
	// convert to uint64
	positionID64, err := strconv.ParseUint(positionIDFields[0], 10, 16)
	if err != nil {
		h.log.Error("cannot convert position ID from string to uint64")
		helpers.BadRequestFieldError(h.log, w, "positionId")
		return
	}
	// and to uint16
	positionID := uint16(positionID64)
	// Assign position to staff
	if err := h.staffService.AssignPosition(r.Context(), userID, positionID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to assign position to staff: %w", err))
		return
	}
	// Get updated staff
	staff, err := h.staffService.GetStaffByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get staff by staff ID (user ID): %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert staff ID (user ID)
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert staff id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get staff position
	response, err := h.staffService.GetPosition(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get staff position by staff id (user id): %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"staffPosition": response,
	})
}
