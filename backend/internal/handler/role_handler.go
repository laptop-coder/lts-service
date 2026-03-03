package handler

import (
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"net/http"
	"strconv"
)

type RoleHandler struct {
	roleService service.RoleService
	log         logger.Logger
}

func NewRoleHandler(roleService service.RoleService, log logger.Logger) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		log:         log,
	}
}

func (h *RoleHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Check method
	if r.Method != http.MethodPut {
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
	// Convert role ID to uint64
	roleID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		h.log.Error("cannot convert role ID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// to uint8:
	roleID := uint8(roleID64)
	// Get permission IDs
	permissionIDsFields := r.PostForm["permissionID"]
	if len(permissionIDsFields) == 0 {
		h.log.Error("the list of permissions cannot be empty")
		helpers.ErrorResponse(w, "the list of permissions cannot be empty", http.StatusBadRequest)
		return
	}
	permissionIDs := make([]uint8, len(permissionIDsFields))
	for i, s := range permissionIDsFields {
		val, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			h.log.Error("cannot convert IDs of permissions from string to uint64")
			helpers.ErrorResponse(w, "cannot convert IDs of permissions from string to uint64", http.StatusInternalServerError)
			return
		}
		permissionIDs[i] = uint8(val)
	}
	// Replace old permissions with new ones
	if err := h.roleService.AssignPermissionsToRole(ctx, roleID, permissionIDs); err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Get updated permissions
	permissions, err := h.roleService.GetRolePermissions(ctx, roleID)
	if err != nil {
		helpers.HandleServiceError(w, err)
		return
	}
	// Return response
	helpers.SuccessResponse(w, map[string]interface{}{
		"roleId":      roleID,
		"permissions": permissions,
		"message":     "permissions updated successfully",
	})
}

func (h *RoleHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		helpers.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Convert role ID to uint64
	roleID64, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		h.log.Error("cannot convert role ID from string to uint64")
		helpers.ErrorResponse(w, "cannot convert role ID from string to uint64", http.StatusInternalServerError)
		return
	}
	// to uint8:
	roleID := uint8(roleID64)
	// Get permissions
	permissions, err := h.roleService.GetRolePermissions(r.Context(), roleID)
	if err != nil {
		helpers.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return response
	helpers.SuccessResponse(w, permissions)
}
