package handler

import (
	"backend/internal/service"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
	log         logger.Logger
}

func NewUserHandler(userService service.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log: log,
	}
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	if err := r.ParseForm(); err != nil {
		errorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// UserID
	// TODO: replace FormValue in the whole code. FormValue takes data not only
	// from POST-params, but also from GET-params
	userIDFields := r.PostForm["userID"]
	if len(userIDFields) > 1 {
		errorResponse(w, "failed to parse form: too much userID fields", http.StatusBadRequest)
		return
	} else if len(userIDFields) == 0 {
		errorResponse(w, "failed to parse form: userID field cannot be empty", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(userIDFields[0])
	if err != nil {
		errorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
	}
	// DTO
	dto := service.UpdateUserDTO{}
	// All fields are optional
	if firstNameFields := r.PostForm["firstName"]; len(firstNameFields) <= 1 {
		dto.FirstName = &firstNameFields[0]
	} else if len(firstNameFields) != 0{
		errorResponse(w, "failed to parse form: to much firstName values", http.StatusBadRequest)
	}
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) <= 1 {
		dto.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0{
		errorResponse(w, "failed to parse form: to much middleName values", http.StatusBadRequest)
	}
	if lastNameFields := r.PostForm["lastName"]; len(lastNameFields) <= 1 {
		dto.LastName = &lastNameFields[0]
	} else if len(lastNameFields) != 0{
		errorResponse(w, "failed to parse form: to much lastName values", http.StatusBadRequest)
	}
	userResponse, err := h.userService.UpdateUser(r.Context(), userID, dto)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to update the user profile: %w", err))
		return
	}
	successResponse(w, map[string]interface{}{
		"user": userResponse,
	})
}

