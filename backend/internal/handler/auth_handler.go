package handler

import (
	"backend/internal/permissions"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type AuthHandler struct {
	authService       service.AuthService
	userService       service.UserService
	inviteService     service.InviteService
	authServiceConfig service.AuthServiceConfig
	log               logger.Logger
}

func NewAuthHandler(authService service.AuthService, userService service.UserService, inviteService service.InviteService, authServiceConfig service.AuthServiceConfig, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		userService:       userService,
		inviteService:     inviteService,
		authServiceConfig: authServiceConfig,
		log:               log,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse multipart/formdata form", http.StatusBadRequest)
		return
	}
	// Get fields data
	fieldsData := make(map[string]string)
	for _, s := range []string{"password", "firstName", "lastName", "inviteToken"} {
		formFields := r.PostForm[s]
		if len(formFields) > 1 {
			helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse form: too much %s fields", s), http.StatusBadRequest)
			return
		} else if len(formFields) == 0 {
			helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse form: %s field cannot be empty", s), http.StatusBadRequest)
			return
		}
		trimmed := strings.TrimSpace(formFields[0]) // TODO: add checks like here in the whole code
		if trimmed == "" {
			helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse form: %s field cannot be empty or only whitespace", s), http.StatusBadRequest)
			return
		}
		fieldsData[s] = trimmed
	}
	// Try to get email from the invite token
	var email *string
	email, err := h.inviteService.GetEmail(r.Context(), fieldsData["inviteToken"])
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	if email == nil {
		// Get email from the form
		emailFields := r.PostForm["email"]
		if len(emailFields) != 1 {
			helpers.ErrorResponse(h.log, w, "email must be specified exactly once in the form", http.StatusBadRequest)
			return
		}
		trimmedEmail := strings.TrimSpace(emailFields[0])
		if trimmedEmail == "" {
			helpers.ErrorResponse(h.log, w, "email cannot be empty or only whitespace", http.StatusBadRequest)
			return
		}
		email = &trimmedEmail
	}
	if email == nil {
		helpers.ErrorResponse(h.log, w, "error getting email", http.StatusInternalServerError)
		return
	}
	// Get roles from invite token
	roles, err := h.inviteService.GetRoles(r.Context(), fieldsData["inviteToken"]) // TODO: change ctx to r.Context() in the whole code
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	if len(roles) == 0 {
		helpers.ErrorResponse(h.log, w, "list of the roles cannot be empty", http.StatusInternalServerError) // HTTP 500 because the token was signed by the server
		return
	}
	userRolesDTO := service.UserRolesDTO{}
	roleIDs := make([]uint8, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}
	// Assemble create user DTO
	createUserDTO := service.CreateUserDTO{
		Email:     *email,
		Password:  fieldsData["password"],
		FirstName: fieldsData["firstName"],
		LastName:  fieldsData["lastName"],
		RoleIDs:   roleIDs,
	}
	// Middle name (optional)
	if middleNameFields := r.PostForm["middleName"]; len(middleNameFields) == 1 {
		createUserDTO.MiddleName = &middleNameFields[0]
	} else if len(middleNameFields) != 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much middleName values", http.StatusBadRequest)
		return
	}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint8
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(h.log, w, "cannot convert teacher classroom ID from string to uint64", http.StatusInternalServerError)
			return
		}
		teacherClassroomID := uint8(teacherClassroomID64)
		userRolesDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much teacher classroom id values", http.StatusBadRequest)
		return
	}
	// TeacherSubjectIDs (special)
	if teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]; len(teacherSubjectIDsFields) == 0 {
		// Check if creating user with the teacher role
		if slices.Contains(roleIDs, 5) { // TODO: make smth like enum for roles constants
			helpers.ErrorResponse(h.log, w, "failed to parse form: at least one teacherSubjectId value must be specified", http.StatusBadRequest)
			return
		}
	} else {
		var teacherSubjectIDs = make([]uint8, len(teacherSubjectIDsFields))
		for i, subjectIDString := range teacherSubjectIDsFields {
			subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
			if err != nil {
				helpers.ErrorResponse(h.log, w, "cannot convert teacher subject ID from string to uint64", http.StatusInternalServerError)
				return
			}
			subjectID8 := uint8(subjectID64)
			teacherSubjectIDs[i] = subjectID8
		}
		userRolesDTO.TeacherSubjectIDs = teacherSubjectIDs
	}
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			helpers.ErrorResponse(h.log, w, "cannot convert student group ID from string to uint64", http.StatusInternalServerError)
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userRolesDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) == 0 {
		// Check if creating user with the student role
		if slices.Contains(roleIDs, 7) {
			helpers.ErrorResponse(h.log, w, "failed to parse form: studentGroupId field must be provided exactly once", http.StatusBadRequest)
			return
		}
	} else {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much student group id values", http.StatusBadRequest)
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint8
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(h.log, w, "cannot convert staff position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		staffPositionID := uint8(staffPositionID64)
		userRolesDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) == 0 {
		// Check if creating user with the staff role
		if slices.Contains(roleIDs, 4) {
			helpers.ErrorResponse(h.log, w, "failed to parse form: staffPositionId field must be provided exactly once", http.StatusBadRequest)
			return
		}
	} else {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much staff position id values", http.StatusBadRequest)
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint8
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 8)
		if err != nil {
			helpers.ErrorResponse(h.log, w, "cannot convert institution administrator position ID from string to uint64", http.StatusInternalServerError)
			return
		}
		institutionAdministratorPositionID := uint8(institutionAdministratorPositionID64)
		userRolesDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) == 0 {
		// Check if creating user with the institution administrator role
		if slices.Contains(roleIDs, 3) {
			helpers.ErrorResponse(h.log, w, "failed to parse form: institutionAdministratorPositionId field must be provided exactly once", http.StatusBadRequest)
			return
		}
	} else {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much institution administrator position id values", http.StatusBadRequest)
		return
	}
	// ParentStudentIDs (special)
	if parentStudentIDsFields := r.PostForm["parentStudentId"]; len(parentStudentIDsFields) != 0 {
		var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
		for i, parentStudentIDString := range parentStudentIDsFields {
			parentStudentID, err := uuid.Parse(parentStudentIDString)
			if err != nil {
				helpers.ErrorResponse(h.log, w, "cannot convert student id to uuid", http.StatusBadRequest)
				return
			}
			parentStudentIDs[i] = parentStudentID
		}
		userRolesDTO.ParentStudentIDs = parentStudentIDs
	}
	// Avatar (optional)
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much avatar files", http.StatusBadRequest)
		return
	} else if len(formFiles) == 1 {
		createUserDTO.Avatar = formFiles[0]
	}
	userResponse, err := h.userService.CreateUser(r.Context(), createUserDTO, userRolesDTO)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to create the user: %w", err))
		return
	}
	// Log in automatically
	tokens, _, err := h.authService.Login(r.Context(), createUserDTO.Email, createUserDTO.Password) // don't rewrite userResponse
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to log in automatically: %w", err))
		return
	}
	// Parse created tokens
	parsedAccessToken, err := h.authService.ParseToken(tokens.AccessToken)
	if err != nil || parsedAccessToken == nil {
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil || parsedRefreshToken == nil {
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse refresh token: %s", err.Error()), http.StatusInternalServerError)
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
	// Revoke invite token
	// TODO: move token revoking to user creation transaction. Is it OK now to
	// return error when user was already created and logged in?
	if err := h.inviteService.RevokeToken(r.Context(), fieldsData["inviteToken"]); err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"user": userResponse,
	},
		http.StatusCreated,
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		helpers.ErrorResponse(h.log, w, "failed to parse x-www-form-urlencoded form", http.StatusBadRequest)
		return
	}
	// Get credentials
	email := r.PostForm["email"]
	if len(email) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: email field is required", http.StatusBadRequest)
		return
	} else if len(email) > 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much email values", http.StatusBadRequest)
		return
	}
	password := r.PostForm["password"]
	if len(password) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: password field is required", http.StatusBadRequest)
		return
	} else if len(password) > 1 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: to much password values", http.StatusBadRequest)
		return
	}
	// Log in
	tokens, userResponse, err := h.authService.Login(r.Context(), email[0], password[0])
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to log in with this credentials: %w", err))
		return
	}
	// Parse created tokens
	parsedAccessToken, err := h.authService.ParseToken(tokens.AccessToken)
	if err != nil || parsedAccessToken == nil {
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil || parsedRefreshToken == nil {
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("failed to parse refresh token: %s", err.Error()), http.StatusInternalServerError)
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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Get refresh token from cookies
	refreshToken, err := helpers.GetCookie("jwt_refresh", r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			helpers.ErrorResponse(h.log, w, "not authorized", http.StatusUnauthorized)
			return
		}
		helpers.ErrorResponse(h.log, w, fmt.Sprintf("error reading refresh token cookie: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Revoke JWT refresh
	if err := h.authService.RevokeToken(r.Context(), refreshToken); err != nil {
		helpers.HandleServiceError(h.log, w, err)
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

// TODO: revoke all user tokens after account delete. Now in middleware.Auth
// there is a check that user from JWT field (user ID) exists.
func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusBadRequest)
		return
	}
	// Get roles of user to delete
	roles, err := h.userService.GetUserRoles(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get user roles (user ID: %s): %w", userID, err))
	}
	// Get auth user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Check if auth user has permission to delete this user
	// TODO: rewrite. Add wrapper, e.g. HasRole; add constants (ADMIN, e.g.)
	hasAdminRole := false
	hasUserRole := false
	for _, role := range roles {
		if role.Name == "admin" {
			hasAdminRole = true
		} else if role.Name == "superadmin" {
			helpers.ErrorResponse(h.log, w, "forbidden: you cannot delete user with superadmin role", http.StatusForbidden)
			return
		} else {
			hasUserRole = true
		}
	}
	if hasAdminRole {
		if !slices.Contains(userPermissions, permissions.UserDeleteAnyAdmin) {
			helpers.ErrorResponse(h.log, w, "forbidden: you do not have permission to delete user with admin role", http.StatusForbidden)
			return
		}
	}
	if hasUserRole {
		if !slices.Contains(userPermissions, permissions.UserDeleteAnyUser) {
			helpers.ErrorResponse(h.log, w, "forbidden: you do not have permission to delete user with non-admin role", http.StatusForbidden)
			return
		}
	}
	// Delete user
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to delete the user: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *AuthHandler) DeleteOwnAccount(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodDelete {
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		helpers.ErrorResponse(h.log, w, "cannot convert user id to uuid", http.StatusUnauthorized)
		return
	}
	// Delete user
	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to delete the user: %w", err))
		return
	}
	// Return response
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}
