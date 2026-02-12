package handler

import (
	"backend/internal/service"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	userService service.UserService
	userServiceConfig service.UserServiceConfig
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
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
