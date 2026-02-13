package handler

import (
	"strconv"
	"fmt"
	"net/http"
	"backend/internal/service"
	"backend/pkg/logger"
)

type StudentGroupHandler struct {
	studentGroupService service.StudentGroupService
	log         logger.Logger
}

func NewStudentGroupHandler(studentGroupService service.StudentGroupService, log logger.Logger) *StudentGroupHandler {
	return &StudentGroupHandler{
		studentGroupService: studentGroupService,
		log:         log,
	}
}

func (h *StudentGroupHandler) GetAdvisorByGroupID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert student group ID.
	// to uint64:
	groupID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		h.log.Error("cannot convert groupID from string to uint64")
		errorResponse(w, "cannot convert groupID from string to uint64", http.StatusInternalServerError)
		return
	}
	// to uint16:
	groupID := uint16(groupID64)
	// Get ID of the group advisor
	response, err := h.studentGroupService.GetAdvisorByGroupID(r.Context(), groupID)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to get student group advisor by group id: %w", err))
		return
	}
	// Return response
	successResponse(w, map[string]interface{}{
		"user": response,
	})
}

