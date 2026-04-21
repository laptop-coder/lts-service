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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert parent ID
	parentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert parent id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get parent
	response, err := h.parentService.GetParentByID(r.Context(), parentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get parent by id: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. parent ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get parent
	response, err := h.parentService.GetParentByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get parent by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"parent": response,
	})
}

func (h *ParentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert parent ID
	parentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert parent id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get parent students
	response, err := h.parentService.GetParentStudents(r.Context(), parentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get parent students by parent id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"students": response,
	})
}

func (h *ParentHandler) GetStudentsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. parent ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get parent students
	response, err := h.parentService.GetParentStudents(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get parent students by parent id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"students": response, // TODO: think about this messages-wrappers (like "students" here) in the whole code of the API
	})
}

func (h *ParentHandler) GetStudentGroupsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. parent ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get student groups of students assigned to parent
	studentGroups, err := h.parentService.GetStudentGroupsOwn(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student groups of students assigned to parent with ID %s", userID))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentGroups": studentGroups,
	})
}

func (h *ParentHandler) AddStudents(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert parent ID
	parentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert parent id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
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
	// Get and convert student IDs:
	studentIDFields := r.PostForm["studentId"]
	if len(studentIDFields) == 0 {
		h.log.Error("failed to parse form: at least one studentID value must be specified")
		helpers.AtLeastOneFieldError(h.log, w, "studentId")
		return
	}
	studentIDs := make([]uuid.UUID, len(studentIDFields))
	for i, studentIDString := range studentIDFields {
		studentID, err := uuid.Parse(studentIDString)
		if err != nil {
			h.log.Error("cannot convert student id to uuid")
			helpers.BadRequestFieldError(h.log, w, "studentId")
			return
		}
		studentIDs[i] = studentID
	}
	// Add students
	if err := h.parentService.AddStudents(r.Context(), parentID, studentIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated parent
	parent, err := h.parentService.GetParentByID(r.Context(), parentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"parent":  parent,
		"message": "students added successfully",
	})
}

func (h *ParentHandler) AddStudentsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. parent ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
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
	// Get and convert student IDs:
	studentIDFields := r.PostForm["studentId"]
	if len(studentIDFields) == 0 {
		h.log.Error("failed to parse form: at least one studentID value must be specified")
		helpers.AtLeastOneFieldError(h.log, w, "studentId")
		return
	}
	studentIDs := make([]uuid.UUID, len(studentIDFields))
	for i, studentIDString := range studentIDFields {
		studentID, err := uuid.Parse(studentIDString)
		if err != nil {
			h.log.Error("cannot convert student id to uuid")
			helpers.BadRequestFieldError(h.log, w, "studentId")
			return
		}
		studentIDs[i] = studentID
	}
	// Add students
	if err := h.parentService.AddStudents(r.Context(), userID, studentIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated parent
	parent, err := h.parentService.GetParentByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"parent":  parent,
		"message": "students added successfully",
	})
}

func (h *ParentHandler) UnassignStudent(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert parent ID
	parentID, err := uuid.Parse(r.PathValue("parentId"))
	if err != nil {
		h.log.Error("cannot convert parent id to uuid")
		helpers.BadRequestFieldError(h.log, w, "parentId")
		return
	}
	// Get and convert student ID:
	studentID, err := uuid.Parse(r.PathValue("studentId"))
	if err != nil {
		h.log.Error("cannot convert student id to uuid")
		helpers.BadRequestFieldError(h.log, w, "studentId")
		return
	}
	// Unassign student
	if err := h.parentService.UnassignStudent(r.Context(), parentID, studentID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *ParentHandler) UnassignStudentOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. parent ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get and convert student ID:
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert student id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Unassign student
	if err := h.parentService.UnassignStudent(r.Context(), userID, studentID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
