package handler

import (
	"strconv"
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
		log:         log,
	}
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	if err := r.ParseForm(); err != nil {
		errorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// UserID
	userIDFields := r.Form["userID"]
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
	if firstNameFields := r.PostForm["firstName"]; len(firstNameFields) == 1 {
		dto.FirstName = &firstNameFields[0]
	} else if len(firstNameFields) != 0 {
		errorResponse(w, "failed to parse form: to much firstName values", http.StatusBadRequest)
	}
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) == 1 {
		dto.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0 {
		errorResponse(w, "failed to parse form: to much middleName values", http.StatusBadRequest)
	}
	if lastNameFields := r.PostForm["lastName"]; len(lastNameFields) == 1 {
		dto.LastName = &lastNameFields[0]
	} else if len(lastNameFields) != 0 {
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

func (h *UserHandler) RemoveAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	if err := r.ParseForm(); err != nil {
		errorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// UserID
	userIDFields := r.Form["userID"]
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
	// Remove user avatar file
	if err := h.userService.RemoveAvatar(r.Context(), userID); err != nil {
		handleServiceError(w, fmt.Errorf("failed to remove user avatar file: %w", err))
		return
	}
	jsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	// UserID
	userIDFields := r.Form["userID"]
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
	// Avatar
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		errorResponse(w, "failed to parse form: to much avatar files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 0 {
		errorResponse(w, "failed to parse form: avatar cannot be empty", http.StatusBadRequest)
		return
	}
	// Update avatar file
	if err := h.userService.UpdateAvatar(r.Context(), userID, formFiles[0]); err != nil {
		handleServiceError(w, fmt.Errorf("failed to update the avatar: %w", err))
		return
	}
	jsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	// UserID
	userIDFields := query["userID"]
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
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to get user by id: %w", err))
		return
	}
	successResponse(w, map[string]interface{}{
		"user": response,
	})
}

func (h *UserHandler) GetStudentGroupAdvisorByGroupID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	// UserID
	groupIDFields := query["groupID"]
	if len(groupIDFields) > 1 {
		errorResponse(w, "failed to parse form: too much groupID fields", http.StatusBadRequest)
		return
	} else if len(groupIDFields) == 0 {
		errorResponse(w, "failed to parse form: groupID field cannot be empty", http.StatusBadRequest)
		return
	}
	// Convert to uint64
	groupID64, err := strconv.ParseUint(groupIDFields[0], 10, 16)
	if err != nil {
		h.log.Error("cannot convert groupID from string to uint64")
		errorResponse(w, "cannot convert groupID from string to uint64", http.StatusInternalServerError)
		return
	}
	// Convert to uint16
	groupID := uint16(groupID64)
	// Get ID of the group advisor
	response, err := h.userService.GetStudentGroupAdvisorByGroupID(r.Context(), groupID)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to get student group advisor by group id: %w", err))
		return
	}
	successResponse(w, map[string]interface{}{
		"user": response,
	})
}

