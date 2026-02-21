package handler

import (
	"errors"
	"backend/internal/service"
	"backend/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService       service.AuthService
	userService       service.UserService
	authServiceConfig service.AuthServiceConfig
	log               logger.Logger
}

func NewAuthHandler(authService service.AuthService, userService service.UserService, authServiceConfig service.AuthServiceConfig, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		userService:       userService,
		authServiceConfig: authServiceConfig,
		log:               log,
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
	// Log in automatically
	tokens, _, err := h.authService.Login(r.Context(), dto.Email, dto.Password) // don't rewrite userResponse
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to log in automatically: %w", err))
		return
	}
	// Parse created tokens
	parsedAccessToken, err := h.authService.ParseToken(tokens.AccessToken)
	if err != nil {
		errorResponse(w, fmt.Sprintf("failed to parse access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil {
		errorResponse(w, fmt.Sprintf("failed to parse refresh token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_access",
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  parsedAccessToken.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   h.authServiceConfig.CookieSecure,
	})
	h.log.Debug("Added JWT access to the cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_refresh",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Expires:  parsedRefreshToken.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   h.authServiceConfig.CookieSecure,
	})
	h.log.Debug("Added JWT refresh to the cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "authorized",
		Value:    "true",
		Path:     "/",
		Expires:  parsedRefreshToken.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: false,
		Secure:   h.authServiceConfig.CookieSecure,
	})
	h.log.Debug("Authorized value is set to true in cookies")
	h.log.Info("User logged in successfully")
	jsonResponse(w, map[string]interface{}{
		"user": userResponse,
	},
		http.StatusCreated,
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		errorResponse(w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get credentials
	email := r.PostForm["email"]
	if len(email) == 0 {
		errorResponse(w, "failed to parse form: email field is required", http.StatusBadRequest)
	} else if len(email) > 1 {
		errorResponse(w, "failed to parse form: to much email values", http.StatusBadRequest)
	}
	password := r.PostForm["password"]
	if len(password) == 0 {
		errorResponse(w, "failed to parse form: password field is required", http.StatusBadRequest)
	} else if len(password) > 1 {
		errorResponse(w, "failed to parse form: to much password values", http.StatusBadRequest)
	}
	// Log in
	tokens, userResponse, err := h.authService.Login(r.Context(), email[0], password[0])
	if err != nil {
		handleServiceError(w, fmt.Errorf("failed to log in with this credentials: %w", err))
		return
	}
	// Parse created tokens
	parsedAccessToken, err := h.authService.ParseToken(tokens.AccessToken)
	if err != nil {
		errorResponse(w, fmt.Sprintf("failed to parse access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil {
		errorResponse(w, fmt.Sprintf("failed to parse refresh token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_access",
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  parsedAccessToken.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   h.authServiceConfig.CookieSecure,
	})
	h.log.Debug("Added JWT access to the cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_refresh",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Expires:  parsedRefreshToken.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   h.authServiceConfig.CookieSecure,
	})
	h.log.Debug("Added JWT refresh to the cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "authorized",
		Value:    "true",
		Path:     "/",
		Expires:  parsedRefreshToken.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: false,
		Secure:   h.authServiceConfig.CookieSecure,
	})
	h.log.Debug("Authorized value is set to true in cookies")
	h.log.Info("User logged in successfully")
	jsonResponse(w, map[string]interface{}{
		"user": userResponse,
	},
		http.StatusOK,
	)
}

// TODO: LogoutAll (i.e. increase version of new tokens, so any other token with
// old version will be considered revoked)

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Get refresh token from cookies
	refreshToken, err := getCookie("jwt_refresh", r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			errorResponse(w, "not authorized", http.StatusUnauthorized)
			return
		}
		errorResponse(w, fmt.Sprintf("error reading refresh token cookie: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Revoke JWT refresh
	if err := h.authService.RevokeToken(r.Context(), refreshToken); err != nil {
		handleServiceError(w, err)
	}
	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_access",
		Value:    "",
		Path:     "/",
		MaxAge: -1,
	})
	h.log.Debug("Cleared JWT access cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_refresh",
		Value:    "",
		Path:     "/",
		MaxAge: -1,
	})
	h.log.Debug("Cleared JWT refresh cookie")
	http.SetCookie(w, &http.Cookie{
		Name:     "authorized",
		Value:    "",
		Path:     "/",
		MaxAge: -1,
	})
	h.log.Debug("Cleared authorized cookie")
	h.log.Info("User logged out successfully")
	jsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
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
