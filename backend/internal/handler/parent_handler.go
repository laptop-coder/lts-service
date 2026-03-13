package handler

import (
	"backend/internal/service"
	"backend/pkg/middleware"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type ParentHandler struct {
	parentService service.ParentService
	log           logger.Logger
}

func NewParentHandler(parentService service.ParentService, log logger.Logger) *ParentHandler {
	return &ParentHandler{
		parentService: parentService,
		log:           log,
	}
}

func (h *ParentHandler) GetParentByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert parent ID
	parentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert parent id to uuid", http.StatusBadRequest)
	}
	// Get parent
	response, err := h.parentService.GetParentByID(r.Context(), parentID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get parent by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"parent": response,
	})
}

func (h *ParentHandler) GetOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. parent ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get parent
	response, err := h.parentService.GetParentByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get parent by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"parent": response,
	})
}
