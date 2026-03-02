package handler

import (
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/helpers"
	"errors"
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.ErrorResponse(w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	fieldsData := make(map[string]string)
	for _, s := range []string{"email", "password", "firstName", "lastName"} {
		formFields := r.PostForm[s]
		if len(formFields) > 1 {
			helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: too much %s fields", s), http.StatusBadRequest)
			return
		} else if len(formFields) == 0 {
			helpers.ErrorResponse(w, fmt.Sprintf("failed to parse form: %s field cannot be empty", s), http.StatusBadRequest)
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
		helpers.ErrorResponse(w, "the list of roles cannot be empty", http.StatusBadRequest)
		return
	} else {
		uints := make([]uint8, len(roleIDs))
		for i, s := range roleIDs {
			val, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				h.log.Error("cannot convert IDs of roles string to uint64")
				helpers.ErrorResponse(w, "cannot convert IDs of roles string to uint64", http.StatusInternalServerError)
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
		helpers.ErrorResponse(w, "failed to parse form: to much middleName values", http.StatusBadRequest)
		return
	}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomID"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint8
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher classroom ID from string to uint64", http.StatusInternalServerError)
			return
		}
		teacherClassroomID := uint8(teacherClassroomID64)
		dto.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much teacher classroom id values", http.StatusBadRequest)
		return
	}
	// TeacherSubjectIDs (special)
	teacherSubjectIDsFields := r.PostForm["teacherSubjectIDs"]
	var teacherSubjectIDs = make([]uint8, len(teacherSubjectIDsFields))
	for i, subjectIDString := range teacherSubjectIDsFields {
		subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert teacher subject ID from string to uint64", http.StatusInternalServerError)
			return
		}
		subjectID8 := uint8(subjectID64)
		teacherSubjectIDs[i] = subjectID8
	}
	dto.TeacherSubjectIDs = teacherSubjectIDs
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupID"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert student group ID from string to uint64", http.StatusInternalServerError)
			return
		}
		studentGroupID := uint16(studentGroupID64)
		dto.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much student group id values", http.StatusBadRequest)
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionID"]; len(staffPositionIDFields) == 1 {
		// Convert to uint8
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert staff position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		staffPositionID := uint8(staffPositionID64)
		dto.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much staff position id values", http.StatusBadRequest)
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionID"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint8
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert institution administrator position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		institutionAdministratorPositionID := uint8(institutionAdministratorPositionID64)
		dto.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) != 0 {
		helpers.ErrorResponse(w, "failed to parse form: to much institution administrator position id values", http.StatusBadRequest)
		return
	}
	// ParentStudentIDs (special)
	parentStudentIDsFields := r.PostForm["parentStudentIDs"]
	var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
	for i, parentStudentIDString := range parentStudentIDsFields {
		parentStudentID, err := uuid.Parse(parentStudentIDString)
		if err != nil {
			helpers.ErrorResponse(w, "cannot convert student id to uuid", http.StatusBadRequest)
			return
		}
		parentStudentIDs[i] = parentStudentID
	}
	dto.ParentStudentIDs = parentStudentIDs
	// Avatar (optional)
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much avatar files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 1 {
		dto.Avatar = formFiles[0]
	}
	userResponse, err := h.userService.CreateUser(r.Context(), dto)
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to create the user: %w", err))
		return
	}
	// Log in automatically
	tokens, _, err := h.authService.Login(r.Context(), dto.Email, dto.Password) // don't rewrite userResponse
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to log in automatically: %w", err))
		return
	}
	// Parse created tokens
	parsedAccessToken, err := h.authService.ParseToken(tokens.AccessToken)
	if err != nil {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse refresh token: %s", err.Error()), http.StatusInternalServerError)
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
	helpers.JsonResponse(w, map[string]interface{}{
		"user": userResponse,
	},
		http.StatusCreated,
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
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
	// Get credentials
	email := r.PostForm["email"]
	if len(email) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: email field is required", http.StatusBadRequest)
		return
	} else if len(email) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much email values", http.StatusBadRequest)
		return
	}
	password := r.PostForm["password"]
	if len(password) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: password field is required", http.StatusBadRequest)
		return
	} else if len(password) > 1 {
		helpers.ErrorResponse(w, "failed to parse form: to much password values", http.StatusBadRequest)
		return
	}
	// Log in
	tokens, userResponse, err := h.authService.Login(r.Context(), email[0], password[0])
	if err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to log in with this credentials: %w", err))
		return
	}
	// Parse created tokens
	parsedAccessToken, err := h.authService.ParseToken(tokens.AccessToken)
	if err != nil {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil {
		helpers.ErrorResponse(w, fmt.Sprintf("failed to parse refresh token: %s", err.Error()), http.StatusInternalServerError)
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
	helpers.JsonResponse(w, map[string]interface{}{
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
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Get refresh token from cookies
	refreshToken, err := helpers.GetCookie("jwt_refresh", r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			helpers.ErrorResponse(w, "not authorized", http.StatusUnauthorized)
			return
		}
		helpers.ErrorResponse(w, fmt.Sprintf("error reading refresh token cookie: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Revoke JWT refresh
	if err := h.authService.RevokeToken(r.Context(), refreshToken); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt_access",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	h.log.Debug("Cleared JWT access cookie")
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt_refresh",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	h.log.Debug("Cleared JWT refresh cookie")
	http.SetCookie(w, &http.Cookie{
		Name:   "authorized",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	h.log.Debug("Cleared authorized cookie")
	h.log.Info("User logged out successfully")
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(w, "cannot convert user id to uuid", http.StatusBadRequest)
	}
	// Delete user
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		helpers.HandleServiceError(w, fmt.Errorf("failed to delete the user: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
