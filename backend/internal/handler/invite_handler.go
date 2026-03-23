package handler

import (
	"slices"
	"backend/internal/service"
	"backend/internal/permissions"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"net/http"
	"strconv"
)

type InviteHandler struct {
	inviteService       service.InviteService
	inviteServiceConfig service.InviteServiceConfig
	log               logger.Logger
}

func NewInviteHandler(inviteService service.InviteService, inviteServiceConfig service.InviteServiceConfig, log logger.Logger) *InviteHandler {
	return &InviteHandler{
		inviteService:       inviteService,
		inviteServiceConfig: inviteServiceConfig,
		log:               log,
	}
}

func (h *InviteHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	// Get and convert (to uint8) role IDs
	roleIDFields := r.PostForm["roleId"]
	if len(roleIDFields) == 0 {
		helpers.ErrorResponse(w, "failed to parse form: at least one roleId value must be provided", http.StatusBadRequest)
		return
	}
	roleIDs := make([]uint8, len(roleIDFields))
	for i, roleIDString := range roleIDFields {
		roleID64, err := strconv.ParseUint(roleIDString, 10, 8)
		if err != nil {
			h.log.Error("cannot convert role ID from string to uint64")
			helpers.ErrorResponse(w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
			return
		}
		roleID := uint8(roleID64)
		roleIDs[i] = roleID
	}
	// Get user permissions
	userPermissions, ok := r.Context().Value(middleware.UserPermissionsKey).([]string)
	if !ok {
		helpers.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Depending on whether token is for admin (role 2) or user (roles 3-7)
	// registration require different permissions
	// admin:
	if slices.Contains(roleIDs, 2) {
		if !slices.Contains(userPermissions, permissions.TokenInviteAdminCreate) {
			helpers.ErrorResponse(w, "forbidden: you do not have permission to create invite token for admin account registration", http.StatusForbidden)
			return
		}
	}
	// user:
	for _, roleID := range roleIDs {
		if slices.Contains([]uint8{3, 4, 5, 6, 7}, roleID) {
			if !slices.Contains(userPermissions, permissions.TokenInviteUserCreate) {
				helpers.ErrorResponse(w, "forbidden: you do not have permission to create invite token for user account registration", http.StatusForbidden)
				return
			}
			break
		}
	}
	// Generate token
	token, err := h.inviteService.CreateToken(r.Context(), roleIDs)
	if err != nil || token == nil {
		helpers.HandleServiceError(w, err)
		return
	}
	helpers.JsonResponse(w, map[string]interface{}{
		"inviteToken": *token,
	}, http.StatusCreated)
}
