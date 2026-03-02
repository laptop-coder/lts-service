package middleware

import (
	"gorm.io/gorm"
	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"context"
	"net/http"
	"slices"
	"time"
)

func Logging(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRolesKey contextKey = "user_roles"
const UserPermissionsKey contextKey = "user_permissions"

func Auth(authService service.AuthService, db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get JWT access from cookies
			cookie, err := helpers.GetCookie("jwt_access", r)
			if err != nil {
				helpers.ErrorResponse(w, "unauthorized", http.StatusUnauthorized) // TODO: maybe improve error handling
				return
			}
			// Validate token
			claims, err := authService.ParseToken(cookie)
			if err != nil {
				helpers.ErrorResponse(w, "invalid token", http.StatusUnauthorized) // TODO: maybe improve error handling
				return
			}
			// Get user permissions
			ctx := r.Context()
			var user model.User
			err = db.WithContext(ctx).
				Preload("Roles").
				Preload("Roles.Permissions").
				First(&user, "id = ?", claims.UserID).Error
			if err != nil {
				helpers.ErrorResponse(w, "failed to load user (by user ID from JWT access from cookies)", http.StatusInternalServerError)
				return
			}
			// Collect permissions
			var permissions []string
			permissionsMap := make(map[string]bool)
			for _, role := range user.Roles {
				for _, permission := range role.Permissions {
					if !permissionsMap[permission.Name] { // remove duplicates
						permissionsMap[permission.Name] = true
						permissions = append(permissions, permission.Name)
					}
				}
			}
			// Collect roles
			var roles []string
			for _, role := range user.Roles {
				roles = append(roles, role.Name)
			}
			// Put user ID, roles and permissions to the context
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRolesKey, roles)
			ctx = context.WithValue(ctx, UserPermissionsKey, permissions)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
