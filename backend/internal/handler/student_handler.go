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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert student ID
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert student id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get student
	response, err := h.studentService.GetStudentByID(r.Context(), studentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student by id: %w", err))
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. student ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get student
	response, err := h.studentService.GetStudentByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"student": response,
	})
}

func (h *StudentHandler) GetClassroom(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert student ID
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert student id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get student classroom
	response, err := h.studentService.GetStudentClassroom(r.Context(), studentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student classroom by student id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentClassroom": response,
	})
}

func (h *StudentHandler) GetClassroomOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. student ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get student classroom
	response, err := h.studentService.GetStudentClassroom(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student classroom by student id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentClassroom": response,
	})
}

func (h *StudentHandler) GetAdvisor(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert student ID
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert student id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get student advisor
	response, err := h.studentService.GetStudentAdvisor(r.Context(), studentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student advisor by student id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentAdvisor": response,
	})
}

func (h *StudentHandler) GetAdvisorOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. student ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get student advisor
	response, err := h.studentService.GetStudentAdvisor(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student advisor by student id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentAdvisor": response,
	})
}

func (h *StudentHandler) GetParents(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert student ID
	studentID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert student id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get student parents
	response, err := h.studentService.GetStudentParents(r.Context(), studentID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student parents by student id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentParents": response,
	})
}

func (h *StudentHandler) GetParentsOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. student ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get student parents
	response, err := h.studentService.GetStudentParents(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student parents by student id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentParents": response,
	})
}

func (h *StudentHandler) GetStudentGroupOwn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID (i.e. student ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
		return
	}
	// Get student
	student, err := h.studentService.GetStudentByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get student by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"studentGroup": student.StudentGroup,
	})
}
