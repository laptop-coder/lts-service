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
