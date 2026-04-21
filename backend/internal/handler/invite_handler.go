package handler

import (
	"backend/internal/permissions"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type InviteHandler struct {
	inviteService       service.InviteService
	inviteServiceConfig service.InviteServiceConfig
	log                 logger.Logger
}

func NewInviteHandler(inviteService service.InviteService, inviteServiceConfig service.InviteServiceConfig, log logger.Logger) *InviteHandler {
	return &InviteHandler{
		inviteService:       inviteService,
		inviteServiceConfig: inviteServiceConfig,
		log:                 log,
	}
}

func (h *InviteHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert (to uint8) role IDs
	roleIDFields := r.PostForm["roleId"]
	if len(roleIDFields) == 0 {
		h.log.Error("failed to parse form: at least one roleId value must be provided")
		helpers.AtLeastOneFieldError(h.log, w, "roleId")
		return
	}
	roleIDs := make([]uint8, len(roleIDFields))
	for i, roleIDString := range roleIDFields {
		roleID64, err := strconv.ParseUint(roleIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert role ID from string to uint64")
			helpers.BadRequestFieldError(h.log, w, "roleId")
			return
		}
		roleID := uint8(roleID64)
		roleIDs[i] = roleID
	}
	// Get email
	emailFields := r.PostForm["email"]
	var email *string
	if len(emailFields) == 1 {
		trimmed := strings.TrimSpace(emailFields[0])
		if trimmed != "" {
			email = &trimmed
		}
	} else if len(emailFields) > 1 {
		h.log.Error("too many email fields")
		helpers.TooManyFieldsError(h.log, w, "email")
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Depending on whether token is for admin (role 2) or user (roles 3-7)
	// registration require different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.TokenInviteAdminCreate) {
			h.log.Error("forbidden: you do not have permission to create invite token for admin account registration")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint8{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.TokenInviteUserCreate) {
				h.log.Error("forbidden: you do not have permission to create invite token for user account registration")
				helpers.ForbiddenError(h.log, w)
				return
			}
			break
		}
	}
	// Generate token
	token, err := h.inviteService.CreateToken(r.Context(), roleIDs, email)
	if err != nil || token == nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"inviteToken": *token,
	}, http.StatusCreated)
}

func (h *InviteHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get token
	token := r.PathValue("token")
	// Parse token
	claims, err := h.inviteService.ParseToken(token)
	if err != nil || claims == nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	// Get roleIDs
	roleIDsInt := claims.RoleIDs
	if len(roleIDsInt) == 0 {
		h.log.Error("list of role IDs cannot be nil")
		helpers.InternalError(h.log, w) // HTTP 500 because token was signed by server
		return
	}
	// Convert role IDs from int to uint8
	roleIDs := make([]uint8, len(roleIDsInt))
	for i, roleID := range roleIDsInt {
		roleIDs[i] = uint8(roleID)
	}
	// TODO: move this code to service layer (the same for the whole code: move
	// business logic to the service layer):
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		h.log.Error("failed to get user permissions from the context")
		helpers.InternalError(h.log, w)
		return
	}
	// Depending on whether token is for admin (role 2) or user (roles 3-7)
	// registration require different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.TokenInviteAdminDelete) {
			h.log.Error("forbidden: you do not have permission to revoke invite token for admin account registration")
			helpers.ForbiddenError(h.log, w)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint8{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.TokenInviteUserDelete) {
				h.log.Error("forbidden: you do not have permission to revoke invite token for user account registration")
				helpers.ForbiddenError(h.log, w)
				return
			}
			break
		}
	}
	// Revoke token
	err = h.inviteService.RevokeToken(r.Context(), token)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{}, http.StatusNoContent)
}

func (h *InviteHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get token
	token := r.PathValue("token")
	// Get roles
	roles, err := h.inviteService.GetRoles(r.Context(), token)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	if len(roles) == 0 {
		helpers.HandleServiceError(h.log, w, fmt.Errorf("list of roles cannot be nil"))
		return
	}
	helpers.SuccessResponse(w, map[string]interface{}{
		"roles": roles,
	})
}

func (h *InviteHandler) GetEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.MethodNotAllowedError(h.log, w)
		return
	}
	// Get token
	token := r.PathValue("token")
	// Get email
	email, err := h.inviteService.GetEmail(r.Context(), token)
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.SuccessResponse(w, map[string]interface{}{
		"email": email,
	})
}

func (h *InviteHandler) MakeStudentInviteRequest(w http.ResponseWriter, r *http.Request) {
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
	// Get email
	emailFields := r.PostForm["email"]
	var email *string
	if len(emailFields) == 1 {
		trimmed := strings.TrimSpace(emailFields[0])
		if trimmed != "" {
			email = &trimmed
		} else {
			h.log.Error("email cannot be empty or only whitespace")
			helpers.FieldRequiredError(h.log, w, "email")
			return
		}
	} else if len(emailFields) != 1 {
		h.log.Error("email must be provided exactly once")
		helpers.FieldExactlyOneError(h.log, w, "email")
		return
	}
	// Make request
	err := h.inviteService.MakeInviteRequest(r.Context(), email, []uint8{7}) // TODO: change "7" to the constant
	if err != nil {
		helpers.HandleServiceError(h.log, w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"message": "the email with the registration link has been sent",
	}, http.StatusAccepted)
}
