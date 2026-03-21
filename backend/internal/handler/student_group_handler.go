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
		helpers.ErrorResponse(w, "cannot convert groupID from string to uint64", http.StatusBadRequest)
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
			return
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
		helpers.ErrorResponse(w, "cannot convert groupID from string to uint64", http.StatusBadRequest) // TODO: maybe change InternalServerError to BadRequest in the similar places in the whole code
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

func (h *StudentGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert student group ID
	studentGroupID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert student group ID from string to uint64", http.StatusBadRequest)
		return
	}
	studentGroupID := uint16(studentGroupID64)
	// Delete student group
	if err := h.studentGroupService.DeleteStudentGroup(r.Context(), studentGroupID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to delete the student group: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *StudentGroupHandler) AssignAdvisor(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert (to uint16) student group ID
	studentGroupID64, err := strconv.ParseUint(r.PathValue("id"), 10, 16)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert student group ID from string to uint64", http.StatusBadRequest) // TODO: use BadRequest instead of InternalServerError in the whole code like here (when cannot convert parameter)
		return
	}
	studentGroupID := uint16(studentGroupID64)
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert user ID
	userIDFields := r.PostForm["userId"]
	if len(userIDFields) != 1 {
		helpers.ErrorResponse(w, "failed to parse form: userID value must be provided exactly once", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(userIDFields[0])
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Assign advisor
	if err := h.studentGroupService.AssignAdvisor(r.Context(), studentGroupID, userID); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Get updated student group
	studentGroup, err := h.studentGroupService.GetStudentGroupByID(r.Context(), studentGroupID)
	if err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{
		"studentGroup": studentGroup,
		"message":      "advisor assigned successfully",
	}, http.StatusCreated) // TODO: check that in other, e.g., POST requests
	                       // there is 201 instead of 200. The same for other
						   // requests.
}
