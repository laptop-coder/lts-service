package handler

import (
	"backend/internal/permissions"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"slices"
	"strconv"
)

type StudentGroupHandler struct {
	teacherService      service.TeacherService
	studentGroupService service.StudentGroupService
	log                 logger.Logger
}

func NewStudentGroupHandler(teacherService service.TeacherService, studentGroupService service.StudentGroupService, log logger.Logger) *StudentGroupHandler {
	return &StudentGroupHandler{
		teacherService:      teacherService,
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
	}, http.StatusCreated)
	// TODO: check that in other, e.g., POST requests
	// there is 201 instead of 200. The same for other
	// requests.
}

func (h *StudentGroupHandler) UnassignAdvisor(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
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
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		helpers.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Check if user unassigning himself
	if slices.Contains(userPermissions, permissions.StudentGroupAdvisorUnassignOwn) && !slices.Contains(userPermissions, permissions.StudentGroupAdvisorUnassignAny) {
		// Get and convert user ID
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			helpers.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// Get teacher
		teacher, err := h.teacherService.GetTeacherByID(r.Context(), userID)
		if err != nil || teacher == nil {
			helpers.HandleServiceError(w, fmt.Errorf("failed to find the teacher by ID: %w", err))
			return
		}
		// Check if the teacher is advisor of the student group
		isAdvisor := false
		if len(teacher.StudentGroups) > 0 {
			for _, group := range teacher.StudentGroups {
				if studentGroupID == group.ID {
					isAdvisor = true
				}
			}
		}
		if !isAdvisor {
			helpers.ErrorResponse(w, "forbidden: you do not have permission to unassign advisor from this student group", http.StatusForbidden)
			return
		}
	}
	// Unassign advisor
	if err := h.studentGroupService.UnassignAdvisor(r.Context(), studentGroupID); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *StudentGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
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
	// DTO
	dto := service.CreateStudentGroupDTO{}
	// Name (required field)
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = nameFields[0]
	} else {
		helpers.ErrorResponse(w, "failed to parse form: name must be provided exactly once", http.StatusBadRequest)
		return
	}
	// Advisor ID (optional field)
	if advisorIDFields := r.PostForm["advisorId"]; len(advisorIDFields) == 1 {
		advisorID, err := uuid.Parse(advisorIDFields[0])
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert advisor id to uuid", http.StatusBadRequest)
			return
		}
		dto.GroupAdvisorID = &advisorID
	} else if len(advisorIDFields) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much group advisor ID values", http.StatusBadRequest)
		return
	}
	// Create student group
	groupResponse, err := h.studentGroupService.CreateStudentGroup(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to create the student group: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{
		"studentGroup": groupResponse,
	},
    http.StatusCreated)
}

func (h *StudentGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPatch {
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
	// Get and convert (to uint16) student group ID:
	groupIDString := r.PathValue("id")
	// convert to uint16
	groupID64, err := strconv.ParseUint(groupIDString, 10, 16)
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert student group ID from string to uint64", http.StatusBadRequest)
		return
	}
	groupID := uint16(groupID64)
	// DTO (all fields are optional)
	dto := service.UpdateStudentGroupDTO{}
	if nameFields := r.PostForm["name"]; len(nameFields) == 1 {
		dto.Name = &nameFields[0]
	} else if len(nameFields) > 1{
		helpers.ErrorResponse(w, "failed to parse form: too much name fields", http.StatusBadRequest)
		return
	}
	if advisorIDFields := r.PostForm["advisorId"]; len(advisorIDFields) == 1 {
		advisorID, err := uuid.Parse(advisorIDFields[0])
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert advisor id to uuid", http.StatusBadRequest)
			return
		}
		dto.GroupAdvisorID = &advisorID
	} else if len(advisorIDFields) > 1{
		helpers.ErrorResponse(w, "failed to parse form: to much group advisor ID values", http.StatusBadRequest)
		return
	}
	// Update student group
	groupResponse, err := h.studentGroupService.UpdateStudentGroup(r.Context(), groupID, dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the student group: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentGroup": groupResponse,
	})
}
