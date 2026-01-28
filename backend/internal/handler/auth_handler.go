package handler

import (
	"backend/internal/service"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	userService service.UserService
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorResponse(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	dto := service.CreateUserDTO{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
	}

	// Middle name (optional)
	if middleName := r.FormValue("middleName"); middleName != "" {
		dto.MiddleName = &middleName
	}

	// Avatar (optional)
	file, fileHeader, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		dto.Avatar = fileHeader
	} else if err != http.ErrMissingFile {
		errorResponse(w, fmt.Sprintf("failed to get avatar: %s", err.Error()), http.StatusBadRequest)
		return
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
