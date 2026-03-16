package handler

import (
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type StudentGroupHandler struct {
	studentGroupService service.StudentGroupService
	log                 logger.Logger
}

func NewStudentGroupHandler(studentGroupService service.StudentGroupService, log logger.Logger) *StudentGroupHandler {
	return &StudentGroupHandler{
		studentGroupService: studentGroupService,
		log:                 log,
	}
}

func (h *StudentGroupHandler) GetAdvisorByGroupID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert student group ID.
	// to uint64:
	groupID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		h.log.Error("cannot convert groupID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert groupID from string to uint64", http.StatusInternalServerError)
		return
	}
	// to uint16:
	groupID := uint16(groupID64)
	// Get student group advisor
	response, err := h.studentGroupService.GetAdvisorByGroupID(r.Context(), groupID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get student group advisor by group id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": response,
	})
}

func (h *StudentGroupHandler) GetStudentGroups(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse query parameters (for filter)
	groupAdvisorIDString := r.URL.Query().Get("groupAdvisorID")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	// Pre-assemble filter (fill with default values)
	filter := repository.StudentGroupFilter{
		Limit:  20,
		Offset: 0,
	}
	// Parse group advisor ID if passed
	if groupAdvisorIDString != "" {
		groupAdvisorID, err := uuid.Parse(groupAdvisorIDString)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert group advisor id from string to uuid", http.StatusBadRequest)
		}
		// Add to filter
		filter.GroupAdvisorID = &groupAdvisorID
	}
	// Parse limit if passed
	if limitString != "" {
		if limit, err := strconv.Atoi(limitString); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100 // max value
			}
			filter.Limit = limit
		} else {
			h.log.Error("invalid limit")
			helpers.ErrorResponse(w, "invalid limit", http.StatusBadRequest)
			return
		}
	}
	// Parse offset if passed
	if offsetString != "" {
		if offset, err := strconv.Atoi(offsetString); err == nil && offset >= 0 {
			filter.Offset = offset
		} else {
			h.log.Error("invalid offset")
			helpers.ErrorResponse(w, "invalid offset", http.StatusBadRequest)
			return
		}
	}
	// Get student groups
	studentGroups, err := h.studentGroupService.GetStudentGroups(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get student groups: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentGroups": studentGroups,
	})
}

func (h *StudentGroupHandler) GetStudentGroupByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert student group ID.
	// to uint64:
	groupID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		h.log.Error("cannot convert groupID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert groupID from string to uint64", http.StatusInternalServerError)
		return
	}
	// to uint16:
	groupID := uint16(groupID64)
	// Get student group
	response, err := h.studentGroupService.GetStudentGroupByID(r.Context(), groupID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get student group by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentGroup": response,
	})
}
