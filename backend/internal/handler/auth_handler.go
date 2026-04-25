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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// TODO: add cleaning of temporary data (ParseMultipartForm, r.MultipartForm, etc)
	// TODO: check if r.MultipartForm == nil and r.PostForm == nil (in all handlers)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error("failed to parse multipart/formdata form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get fields data
	fieldsData := make(map[string]string)
	for _, s := range []string{"password", "firstName", "lastName", "inviteToken"} {
		formFields := r.PostForm[s]
		if len(formFields) > 1 {
			h.log.Error(fmt.Sprintf("failed to parse form: too much %s fields", s))
			helpers.TooManyFieldsError(h.log, w, s)
			return
		} else if len(formFields) == 0 {
			h.log.Error(fmt.Sprintf("failed to parse form: %s field cannot be empty", s))
			helpers.FieldRequiredError(h.log, w, s)
			return
		}
		trimmed := strings.TrimSpace(formFields[0]) // TODO: add checks like here in the whole code
		if trimmed == "" {
			h.log.Error(fmt.Sprintf("failed to parse form: %s field cannot be empty or only whitespace", s))
			helpers.FieldRequiredError(h.log, w, s)
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
			h.log.Error("email must be specified exactly once in the form")
			helpers.FieldExactlyOneError(h.log, w, "email")
			return
		}
		trimmedEmail := strings.TrimSpace(emailFields[0])
		if trimmedEmail == "" {
			h.log.Error("email cannot be empty or only whitespace")
			helpers.FieldRequiredError(h.log, w, "email")
			return
		}
		email = &trimmedEmail
	}
	if email == nil {
		h.log.Error("error getting email")
		helpers.InternalError(h.log, w)
		return
	}
	// Get roles from invite token
	roles, err := h.inviteService.GetRoles(r.Context(), fieldsData["inviteToken"]) // TODO: change ctx to r.Context() in the whole code
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	if roles == nil || len(roles) == 0 { // TODO: check if roles == nil in the whole code
		h.log.Error("list of the roles cannot be empty")
		helpers.InternalError(h.log, w) // HTTP 500 because the token was signed by the server
		return
	}

	userExtensionsDTO := service.UserExtensionsDTO{}
	roleIDs := make([]uint16, len(roles))
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
		h.log.Error("failed to parse form: too many middleName values")
		helpers.TooManyFieldsError(h.log, w, "middleName")
		return
	}
	// TeacherClassroomID (special)
	if teacherClassroomIDFields := r.PostForm["teacherClassroomId"]; len(teacherClassroomIDFields) == 1 {
		// Convert to uint16
		teacherClassroomID64, err := strconv.ParseUint(teacherClassroomIDFields[0], 10, 8)
		if err != nil {
			h.log.Error("cannot convert teacher classroom ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "teacherClassroomId")
			return
		}
		teacherClassroomID := uint16(teacherClassroomID64)
		userExtensionsDTO.TeacherClassroomID = &teacherClassroomID
	} else if len(teacherClassroomIDFields) != 0 {
		h.log.Error("failed to parse form: too many teacher classroom id values")
		helpers.TooManyFieldsError(h.log, w, "teacherClassroomId")
		return
	}
	// TeacherSubjectIDs (special)
	if teacherSubjectIDsFields := r.PostForm["teacherSubjectId"]; len(teacherSubjectIDsFields) == 0 {
		// Check if creating user with the teacher role
		if slices.Contains(roleIDs, 5) { // TODO: make smth like enum for roles constants
			h.log.Error("failed to parse form: at least one teacherSubjectId value must be specified")
			helpers.AtLeastOneFieldError(h.log, w, "teacherSubjectId")
			return
		}
	} else {
		var teacherSubjectIDs = make([]uint16, len(teacherSubjectIDsFields))
		for i, subjectIDString := range teacherSubjectIDsFields {
			subjectID64, err := strconv.ParseUint(subjectIDString, 10, 8)
			if err != nil {
				h.log.Error("cannot convert teacher subject ID from string to uint64")
				helpers.BadRequestFieldError(h.log, w, "teacherSubjectId")
				return
			}
			subjectID8 := uint16(subjectID64)
			teacherSubjectIDs[i] = subjectID8
		}
		userExtensionsDTO.TeacherSubjectIDs = teacherSubjectIDs
	}
	// TeacherStudentGroupIDs (special)
	if teacherStudentGroupIDsFields := r.PostForm["teacherStudentGroupId"]; len(teacherStudentGroupIDsFields) != 0 {
		var teacherStudentGroupIDs = make([]uint16, len(teacherStudentGroupIDsFields))
		for i, groupIDString := range teacherStudentGroupIDsFields {
			groupID64, err := strconv.ParseUint(groupIDString, 10, 16)
			if err != nil {
				h.log.Error("cannot convert teacher student group ID from string to uint64")
				helpers.BadRequestFieldError(h.log, w, "teacherStudentGroupId")
				return
			}
			groupID16 := uint16(groupID64)
			teacherStudentGroupIDs[i] = groupID16
		}
		userExtensionsDTO.TeacherStudentGroupIDs = teacherStudentGroupIDs
	}
	// StudentGroupID (special)
	if studentGroupIDFields := r.PostForm["studentGroupId"]; len(studentGroupIDFields) == 1 {
		// Convert to uint16
		studentGroupID64, err := strconv.ParseUint(studentGroupIDFields[0], 10, 16)
		if err != nil {
			h.log.Error("cannot convert student group ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "studentGroupId")
			return
		}
		studentGroupID := uint16(studentGroupID64)
		userExtensionsDTO.StudentGroupID = &studentGroupID
	} else if len(studentGroupIDFields) == 0 {
		// Check if creating user with the student role
		if slices.Contains(roleIDs, 7) {
			h.log.Error("failed to parse form: studentGroupId field must be provided exactly once")
			helpers.FieldExactlyOneError(h.log, w, "studentGroupId")
			return
		}
	} else {
		h.log.Error("failed to parse form: too many student group id values")
		helpers.TooManyFieldsError(h.log, w, "studentGroupId")
		return
	}
	// StaffPositionID (special)
	if staffPositionIDFields := r.PostForm["staffPositionId"]; len(staffPositionIDFields) == 1 {
		// Convert to uint16
		staffPositionID64, err := strconv.ParseUint(staffPositionIDFields[0], 10, 8)
		if err != nil {
			h.log.Error("cannot convert staff position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "staffPositionId")
			return
		}
		staffPositionID := uint16(staffPositionID64)
		userExtensionsDTO.StaffPositionID = &staffPositionID
	} else if len(staffPositionIDFields) == 0 {
		// Check if creating user with the staff role
		if slices.Contains(roleIDs, 4) {
			h.log.Error("failed to parse form: staffPositionId field must be provided exactly once")
			helpers.FieldExactlyOneError(h.log, w, "staffPositionId")
			return
		}
	} else {
		h.log.Error("failed to parse form: too many staff position id values")
		helpers.TooManyFieldsError(h.log, w, "staffPositionId")
		return
	}
	// InstitutionAdministratorPositionID (special)
	if institutionAdministratorPositionIDFields := r.PostForm["institutionAdministratorPositionId"]; len(institutionAdministratorPositionIDFields) == 1 {
		// Convert to uint16
		institutionAdministratorPositionID64, err := strconv.ParseUint(institutionAdministratorPositionIDFields[0], 10, 8)
		if err != nil {
			h.log.Error("cannot convert institution administrator position ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "institutionAdministratorPositionId")
			return
		}
		institutionAdministratorPositionID := uint16(institutionAdministratorPositionID64)
		userExtensionsDTO.InstitutionAdministratorPositionID = &institutionAdministratorPositionID
	} else if len(institutionAdministratorPositionIDFields) == 0 {
		// Check if creating user with the institution administrator role
		if slices.Contains(roleIDs, 3) {
			h.log.Error("failed to parse form: institutionAdministratorPositionId field must be provided exactly once")
			helpers.FieldExactlyOneError(h.log, w, "institutionAdministratorPositionId")
			return
		}
	} else {
		h.log.Error("failed to parse form: too many institution administrator position id values")
		helpers.BadRequestFieldError(h.log, w, "institutionAdministratorPositionId")
		return
	}
	// ParentStudentIDs (special)
	if parentStudentIDsFields := r.PostForm["parentStudentId"]; len(parentStudentIDsFields) != 0 {
		var parentStudentIDs = make([]uuid.UUID, len(parentStudentIDsFields))
		for i, parentStudentIDString := range parentStudentIDsFields {
			parentStudentID, err := uuid.Parse(parentStudentIDString)
			if err != nil {
				h.log.Error("cannot convert parent student id to uuid")
				helpers.BadRequestFieldError(h.log, w, "parentStudentId")
				return
			}
			parentStudentIDs[i] = parentStudentID
		}
		userExtensionsDTO.ParentStudentIDs = parentStudentIDs
	}
	// Avatar (optional)
	formFiles := r.MultipartForm.File["avatar"]
	if len(formFiles) > 1 {
		h.log.Error("failed to parse form: too many avatar files")
		helpers.TooManyFieldsError(h.log, w, "avatar")
		return
	} else if len(formFiles) == 1 {
		createUserDTO.Avatar = formFiles[0]
	}
	userResponse, err := h.userService.CreateUser(r.Context(), createUserDTO, userExtensionsDTO)
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
		h.log.Error(fmt.Sprintf("failed to parse access token: %s", err.Error()))
		helpers.InternalError(h.log, w)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil || parsedRefreshToken == nil {
		h.log.Error(fmt.Sprintf("failed to parse refresh token: %s", err.Error()))
		helpers.InternalError(h.log, w)
		return
	}
	// TODO: in the whole code check if parsedAccessToken.RegisteredClaims.ExpiresAt == nil. The same for refresh
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Parse form
	if err := r.ParseForm(); err != nil {
		h.log.Error("failed to parse x-www-form-urlencoded form")
		helpers.BadRequestError(h.log, w)
		return
	}
	// Get credentials
	email := r.PostForm["email"]
	if len(email) == 0 {
		h.log.Error("failed to parse form: email field is required")
		helpers.FieldRequiredError(h.log, w, "email")
		return
	} else if len(email) > 1 {
		h.log.Error("failed to parse form: too many email values")
		helpers.TooManyFieldsError(h.log, w, "email")
		return
	}
	password := r.PostForm["password"]
	if len(password) == 0 {
		h.log.Error("failed to parse form: password field is required")
		helpers.FieldRequiredError(h.log, w, "password")
		return
	} else if len(password) > 1 {
		h.log.Error("failed to parse form: too many password values")
		helpers.TooManyFieldsError(h.log, w, "password")
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
		h.log.Error(fmt.Sprintf("failed to parse access token: %s", err.Error()))
		helpers.InternalError(h.log, w)
		return
	}
	parsedRefreshToken, err := h.authService.ParseToken(tokens.RefreshToken)
	if err != nil || parsedRefreshToken == nil {
		h.log.Error(fmt.Sprintf("failed to parse refresh token: %s", err.Error()))
		helpers.InternalError(h.log, w)
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Restrictions
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	// Get refresh token from cookies
	refreshToken, err := helpers.GetCookie("jwt_refresh", r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.log.Warn("logout attempt without refresh token")
			helpers.UnauthorizedError(h.log, w)
			return
		}
		h.log.Error(fmt.Sprintf("error reading refresh token cookie: %s", err.Error()))
		helpers.InternalError(h.log, w)
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.log.Error("cannot convert user id to uuid")
		helpers.BadRequestFieldError(h.log, w, "id")
		return
	}
	// Get roles of user to delete
	roles, err := h.userService.GetUserRoles(r.Context(), userID)
	if err != nil {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("failed to get user roles (user ID: %s): %w", userID, err))
		return
	}
	// Get auth user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
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
			h.log.Error("forbidden: you cannot delete user with superadmin role")
			helpers.ForbiddenError(h.log, w)
			return
		} else {
			hasUserRole = true
		}
	}
	if hasAdminRole {
		if !slices.Contains(userPermissions, permissions.UserDeleteAnyAdmin) {
			h.log.Error("forbidden: you do not have permission to delete user with admin role")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	if hasUserRole {
		if !slices.Contains(userPermissions, permissions.UserDeleteAnyUser) {
			h.log.Error("forbidden: you do not have permission to delete user with non-admin role")
			helpers.ForbiddenError(h.log, w)
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
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get and convert user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		h.log.Error("failed to get userID from context and convert it to UUID")
		helpers.InternalError(h.log, w)
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
