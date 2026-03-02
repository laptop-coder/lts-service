package handler

import (
	"backend/internal/service"
	"backend/pkg/logger"
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
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert student ID
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errorResponse(w, "cannot convert student id to uuid", http.StatusBadRequest)
	}
	// Get student
	response, err := h.studentService.GetStudentByID(r.Context(), studentID)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to get student by id: %w", err))
		return
	}
	// Return response
	successResponse(w, map[string]interface{}{
		"student": response,
	})
}
