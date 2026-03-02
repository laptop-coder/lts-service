package handler

import (
	"backend/internal/service"
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
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert parent ID
	parentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errorResponse(w, "cannot convert parent id to uuid", http.StatusBadRequest)
	}
	// Get parent
	response, err := h.parentService.GetParentByID(r.Context(), parentID)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to get parent by id: %w", err))
		return
	}
	// Return response
	successResponse(w, map[string]interface{}{
		"parent": response,
	})
}
