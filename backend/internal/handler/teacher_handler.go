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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert teacher id to uuid", http.StatusBadRequest)
	}
	// Get teacher
	response, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get teacher by id: %w", err))
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher
	response, err := h.teacherService.GetTeacherByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get teacher by id: %w", err))
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert teacher id to uuid", http.StatusBadRequest)
	}
	// Get teacher classroom
	response, err := h.teacherService.GetTeacherClassroom(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get teacher classroom by teacher id: %w", err))
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher classroom
	response, err := h.teacherService.GetTeacherClassroom(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get teacher classroom by teacher id: %w", err))
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert teacher id to uuid", http.StatusBadRequest)
	}
	// Get teacher subjects
	response, err := h.teacherService.GetTeacherSubjects(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get teacher subjects by teacher id: %w", err))
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. teacher ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get teacher subjects
	response, err := h.teacherService.GetTeacherSubjects(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get teacher subjects by teacher id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacherSubjects": response,
	})
}

func (h *TeacherHandler) AssignClassroom(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert teacher ID
	teacherID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert teacher id to uuid", http.StatusBadRequest)
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get and convert classroom ID:
	classroomIDFields := r.PostForm["classroomID"]
	if len(classroomIDFields) != 1 {
		helpers.ErrorResponse(w, "failed to parse form: classroomID value must be provided exactly once", http.StatusBadRequest)
	}
	// convert to uint64
	classroomID64, err := strconv.ParseUint(classroomIDFields[0], 10, 8)
	if err != nil {
		h.log.Error("cannot convert classroom ID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert classroom ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// and to uint8
	classroomID := uint8(classroomID64)
	// Assign room
	if err := h.teacherService.AssignClassroom(r.Context(), teacherID, classroomID); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Get updated teacher
	teacher, err := h.teacherService.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"teacher": teacher,
		"message": "classroom assigned successfully",
	})
}
