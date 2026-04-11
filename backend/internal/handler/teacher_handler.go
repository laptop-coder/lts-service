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

type TeacherHandler struct {
	teacherService service.TeacherService
	log            logger.Logger
}

func NewTeacherHandler(teacherService service.TeacherService, log logger.Logger) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
		log:            log,
	}
}

func (h *TeacherHandler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Get teacher
	response, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": response,
	})
}

func (h *TeacherHandler) GetOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher
	response, err := h.teacherService.GetTeacherByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": response,
	})
}

func (h *TeacherHandler) GetClassroom(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Get teacher classroom
	response, err := h.teacherService.GetTeacherClassroom(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher classroom by teacher id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacherClassroom": response,
	})
}

func (h *TeacherHandler) GetClassroomOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher classroom
	response, err := h.teacherService.GetTeacherClassroom(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher classroom by teacher id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacherClassroom": response,
	})
}

func (h *TeacherHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Get teacher subjects
	response, err := h.teacherService.GetTeacherSubjects(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher subjects by teacher id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacherSubjects": response,
	})
}

func (h *TeacherHandler) GetSubjectsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher subjects
	response, err := h.teacherService.GetTeacherSubjects(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher subjects by teacher id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacherSubjects": response,
	})
}

func (h *TeacherHandler) AssignClassroom(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert classroom ID:
	classroomIDFields := r.PostForm["classroomId"]
	if len(classroomIDFields) != 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: classroomID value must be provided exactly once", http.StatusBadRequest)
		return
	}
	// convert to uint64
	classroomID64, err := strconv.ParseUint(classroomIDFields[0], 10, 8)
	if err != nil {
		h.log.Error("cannot convert classroom ID from string to uint64")
		helpers.ErrorResponse(h.log, w, "cannot convert classroom ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// and to uint8
	classroomID := uint8(classroomID64)
	// Assign room
	if err := h.teacherService.AssignClassroom(r.Context(), teacherID, classroomID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "classroom assigned successfully",
	})
}

func (h *TeacherHandler) AssignClassroomOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert classroom ID:
	classroomIDFields := r.PostForm["classroomId"]
	if len(classroomIDFields) != 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: classroomID value must be provided exactly once", http.StatusBadRequest)
		return
	}
	// convert to uint64
	classroomID64, err := strconv.ParseUint(classroomIDFields[0], 10, 8)
	if err != nil {
		h.log.Error("cannot convert classroom ID from string to uint64")
		helpers.ErrorResponse(h.log, w, "cannot convert classroom ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// and to uint8
	classroomID := uint8(classroomID64)
	// Assign room
	if err := h.teacherService.AssignClassroom(r.Context(), userID, classroomID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "classroom assigned successfully",
	})
}

func (h *TeacherHandler) UnassignClassroom(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Unassign room
	if err := h.teacherService.UnassignClassroom(r.Context(), teacherID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *TeacherHandler) UnassignClassroomOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Unassign room
	if err := h.teacherService.UnassignClassroom(r.Context(), userID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *TeacherHandler) AssignSubjects(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert subject IDs:
	subjectIDFields := r.PostForm["subjectId"]
	if len(subjectIDFields) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: subjectID value cannot be empty", http.StatusBadRequest)
		return
	}
	subjectIDs := make([]uint8, len(subjectIDFields))
	for i, subjectIDString := range subjectIDFields {
		// convert to uint64
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert subject ID from string to uint64")
			helpers.ErrorResponse(h.log, w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		// and to uint8
		subjectID := uint8(subjectID64)
		subjectIDs[i] = subjectID
	}
	// Assign subjects
	if err := h.teacherService.AssignSubjects(r.Context(), teacherID, subjectIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "subjects assigned successfully",
	})
}

func (h *TeacherHandler) AssignSubjectsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert subject IDs:
	subjectIDFields := r.PostForm["subjectId"]
	if len(subjectIDFields) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: subjectID value cannot be empty", http.StatusBadRequest)
		return
	}
	subjectIDs := make([]uint8, len(subjectIDFields))
	for i, subjectIDString := range subjectIDFields {
		// convert to uint64
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert subject ID from string to uint64")
			helpers.ErrorResponse(h.log, w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		// and to uint8
		subjectID := uint8(subjectID64)
		subjectIDs[i] = subjectID
	}
	// Assign subjects
	if err := h.teacherService.AssignSubjects(r.Context(), userID, subjectIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "subjects assigned successfully",
	})
}

func (h *TeacherHandler) AddSubjects(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert subject IDs:
	subjectIDFields := r.PostForm["subjectId"]
	if len(subjectIDFields) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: subjectID value cannot be empty", http.StatusBadRequest)
		return
	}
	subjectIDs := make([]uint8, len(subjectIDFields))
	for i, subjectIDString := range subjectIDFields {
		// convert to uint64
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert subject ID from string to uint64")
			helpers.ErrorResponse(h.log, w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		// and to uint8
		subjectID := uint8(subjectID64)
		subjectIDs[i] = subjectID
	}
	// Add subjects
	if err := h.teacherService.AddSubjects(r.Context(), teacherID, subjectIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "subjects added successfully",
	})
}

func (h *TeacherHandler) AddSubjectsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert subject IDs:
	subjectIDFields := r.PostForm["subjectId"]
	if len(subjectIDFields) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: subjectID value cannot be empty", http.StatusBadRequest)
		return
	}
	subjectIDs := make([]uint8, len(subjectIDFields))
	for i, subjectIDString := range subjectIDFields {
		// convert to uint64
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert subject ID from string to uint64")
			helpers.ErrorResponse(h.log, w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		// and to uint8
		subjectID := uint8(subjectID64)
		subjectIDs[i] = subjectID
	}
	// Assign subjects
	if err := h.teacherService.AddSubjects(r.Context(), userID, subjectIDs); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "subjects assigned successfully",
	})
}

func (h *TeacherHandler) UnassignSubject(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert teacher id to uuid", http.StatusBadRequest)
		return
	}
	// Get and convert subject ID:
	subjectIDString := r.PathValue("subjectId")
	// convert to uint64
	subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
	if err != nil {
		h.log.Error("cannot convert subject ID from string to uint64")
		helpers.ErrorResponse(h.log, w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// and to uint8
	subjectID := uint8(subjectID64)
	// Unassign subject
	if err := h.teacherService.UnassignSubject(r.Context(), teacherID, subjectID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *TeacherHandler) UnassignSubjectOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get and convert subject ID:
	subjectIDString := r.PathValue("id")
	// convert to uint64
	subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
	if err != nil {
		h.log.Error("cannot convert subject ID from string to uint64")
		helpers.ErrorResponse(h.log, w, "cannot convert subject ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// and to uint8
	subjectID := uint8(subjectID64)
	// Unassign subject
	if err := h.teacherService.UnassignSubject(r.Context(), userID, subjectID); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *TeacherHandler) GetStudentGroupsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get teacher by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentGroups": teacher.StudentGroups,
	})
}
