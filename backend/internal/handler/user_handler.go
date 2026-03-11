package handler

import (
	"strconv"
	"backend/internal/service"
	"backend/internal/repository"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
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

func (h *UserHandler) UpdateOwnProfile(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// DTO (all fields are optional)
	dto := service.UpdateUserDTO{}
	if firstNameFields := r.PostForm["firstName"]; len(firstNameFields) == 1 {
		dto.FirstName = &firstNameFields[0]
	} else if len(firstNameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much firstName values", http.StatusBadRequest)
	}
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) == 1 {
		dto.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much middleName values", http.StatusBadRequest)
	}
	if lastNameFields := r.PostForm["lastName"]; len(lastNameFields) == 1 {
		dto.LastName = &lastNameFields[0]
	} else if len(lastNameFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much lastName values", http.StatusBadRequest)
	}
	// Update user
	userResponse, err := h.userService.UpdateUser(r.Context(), userID, dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the user profile: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": userResponse,
	})
}

func (h *UserHandler) RemoveOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Remove user avatar file
	if err := h.userService.RemoveAvatar(r.Context(), userID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to remove user avatar file: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) UpdateOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPut {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// Parse form
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.ErrorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get avatar file from the request
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much avatar files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: avatar cannot be empty", http.StatusBadRequest)
		return
	}
	// Update avatar file
	if err := h.userService.UpdateAvatar(r.Context(), userID, formFiles[0]); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to update the avatar: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
	}
	// Get user
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get user by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": response,
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse query parameters (for filter)
	roleIDString := r.URL.Query().Get("roleId")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
    // Pre-assemble filter (fill with default values)
	filter := repository.UserFilter {
		Limit: 20,
		Offset: 0,
	}
	// Parse role ID if passed
	if roleIDString != "" {
		// convert to uint64
		roleID64, err := strconv.ParseUint(roleIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert role ID from string to uint64")
			helpers.ErrorResponse(w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
			return
		}
		// and to uint8
		roleID := uint8(roleID64)
		// Add to filter
		filter.RoleID = &roleID
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
	// Get users
	users, err := h.userService.GetUsers(r.Context(), filter)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get users: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"users": users,
	})
}

func (h *UserHandler) GetOwnUser(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(w, "cannot convert own user id to uuid", http.StatusUnauthorized)
		return
	}
	// Get user
	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to get own user by id: %w", err))
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"user": response,
	})
}

