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

type StudentHandler struct {
	studentService service.StudentService
	log            logger.Logger
}

func NewStudentHandler(studentService service.StudentService, log logger.Logger) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
		log:            log,
	}
}

func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert student ID
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert student id to uuid", http.StatusBadRequest)
	}
	// Get student
	response, err := h.studentService.GetStudentByID(r.Context(), studentID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get student by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"student": response,
	})
}

func (h *StudentHandler) GetOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID (i.e. student ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get student
	response, err := h.studentService.GetStudentByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get student by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"student": response,
	})
}
