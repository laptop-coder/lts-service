package handler

import (
	"backend/internal/service"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	userService       service.UserService
	userServiceConfig service.UserServiceConfig
	log               logger.Logger
}

func NewAuthHandler(userService service.UserService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		log:         log,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	fieldsData := make(map[string]string)
	for _, s := range []string{"email", "password", "firstName", "lastName"} {
		formFields := r.PostForm[s]
		if len(formFields) > 1 {
			errorResponse(w, fmt.Sprintf("failed to parse form: too much %s fields", s), http.StatusBadRequest)
			return
		} else if len(formFields) == 0 {
			errorResponse(w, fmt.Sprintf("failed to parse form: %s field cannot be empty", s), http.StatusBadRequest)
			return
		}
		fieldsData[s] = formFields[0]
	}
	dto := service.CreateUserDTO{
		Email:     fieldsData["email"],
		Password:  fieldsData["password"],
		FirstName: fieldsData["firstName"],
		LastName:  fieldsData["lastName"],
	}
	if roleIDs := r.PostForm["roleID"]; len(roleIDs) == 0 { // TODO: maybe this check won't work
		h.log.Error("the list of roles cannot be empty")
		errorResponse(w, "the list of roles cannot be empty", http.StatusBadRequest)
		return
	} else {
		uints := make([]uint8, len(roleIDs))
		for i, s := range roleIDs {
			val, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				h.log.Error("cannot convert IDs of roles string to uint64")
				errorResponse(w, "cannot convert IDs of roles string to uint64", http.StatusInternalServerError)
				return
			}
			uints[i] = uint8(val)
		}
		dto.RoleIDs = uints
	}
	// Middle name (optional)
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) == 1 {
		dto.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0 {
		errorResponse(w, "failed to parse form: to much middleName values", http.StatusBadRequest)
		return
	}
	// Avatar (optional)
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		errorResponse(w, "failed to parse form: to much avatar files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 1 {
		dto.Avatar = formFiles[0]
	}
	userResponse, err := h.userService.CreateUser(r.Context(), dto)
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to create the user: %w", err))
		return
	}
	// TODO: automatically login user here. Get token, set cookie
	jsonResponse(w, map[string]interface{}{
		"user": userResponse,
	},
		http.StatusCreated,
	)
}

func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
	}
	// Delete user
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		handleServiceError(w, fmt.Errorf("failed to delete the user: %w", err))
		return
	}
	// Return response
	jsonResponse(w, map[string]interface{}{}, http.StatusNoContent)

}
