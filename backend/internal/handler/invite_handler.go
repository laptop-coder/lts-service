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
	// Get and convert (to uint8) role IDs
	roleIDFields := r.PostForm["roleId"]
	if len(roleIDFields) == 0 {
		helpers.ErrorResponse(h.log, w, "failed to parse form: at least one roleId value must be provided", http.StatusBadRequest)
		return
	}
	roleIDs := make([]uint8, len(roleIDFields))
	for i, roleIDString := range roleIDFields {
		roleID64, err := strconv.ParseUint(roleIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert role ID from string to uint64")
			helpers.ErrorResponse(h.log, w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
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
		helpers.ErrorResponse(h.log, w, "too much email fields", http.StatusBadRequest)
		return
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Depending on whether token is for admin (role 2) or user (roles 3-7)
	// registration require different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.TokenInviteAdminCreate) {
			helpers.ErrorResponse(h.log, w, "forbidden: you do not have permission to create invite token for admin account registration", http.StatusForbidden)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint8{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.TokenInviteUserCreate) {
				helpers.ErrorResponse(h.log, w, "forbidden: you do not have permission to create invite token for user account registration", http.StatusForbidden)
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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
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
		helpers.ErrorResponse(h.log, w, "list of role IDs cannot be nil", http.StatusInternalServerError) // HTTP 500 because token was signed by server
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
		helpers.ErrorResponse(h.log, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Depending on whether token is for admin (role 2) or user (roles 3-7)
	// registration require different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.TokenInviteAdminDelete) {
			helpers.ErrorResponse(h.log, w, "forbidden: you do not have permission to revoke invite token for admin account registration", http.StatusForbidden)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint8{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.TokenInviteUserDelete) {
				helpers.ErrorResponse(h.log, w, "forbidden: you do not have permission to revoke invite token for user account registration", http.StatusForbidden)
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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
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
		helpers.ErrorResponse(h.log, w, "method not allowed", http.StatusMethodNotAllowed)
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
	// Get email
	emailFields := r.PostForm["email"]
	var email *string
	if len(emailFields) == 1 {
		trimmed := strings.TrimSpace(emailFields[0])
		if trimmed != "" {
			email = &trimmed
		} else {
			helpers.ErrorResponse(h.log, w, "email cannot be empty or only whitespace", http.StatusBadRequest)
			return
		}
	} else if len(emailFields) != 1 {
		helpers.ErrorResponse(h.log, w, "email must be provided exactly once", http.StatusBadRequest)
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
