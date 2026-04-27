package middleware

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/env"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://10.0.6.100:%s", env.GetStringRequired("FRONTEND_PORT")))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRolesKey contextKey = "user_roles"
const UserPermissionsKey contextKey = "user_permissions"

func Auth(authService service.AuthService, authServiceConfig service.AuthServiceConfig, jwtRepo repository.JWTRepository, db *gorm.DB, log logger.Logger, allowUnauthorized bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get JWT access from cookies
			jwtAccess, err := helpers.GetCookie("jwt_access", r)
			if err != nil {
				ctx := r.Context()
				// Try to refresh access-token
				// Get refresh token
				jwtRefresh, err := helpers.GetCookie("jwt_refresh", r)
				if err != nil {
					if allowUnauthorized {
						ctx = context.WithValue(ctx, UserIDKey, uuid.Nil)
						ctx = context.WithValue(ctx, UserRolesKey, []string{})
						ctx = context.WithValue(ctx, UserPermissionsKey, []string{})
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					helpers.ErrorResponse(log, w, "unauthorized", http.StatusUnauthorized)
					log.Error(fmt.Sprintf("Failed to get refresh token from cookies: %s", err.Error()))
					return
				}
				// Generate new token pair
				tokens, err := authService.RefreshToken(ctx, jwtRefresh)
				// TODO: in the whole backend code fix situations like this.
				// When pointer is returned check if it nil.
				if err != nil || tokens == nil {
					if allowUnauthorized {
						ctx = context.WithValue(ctx, UserIDKey, uuid.Nil)
						ctx = context.WithValue(ctx, UserRolesKey, []string{})
						ctx = context.WithValue(ctx, UserPermissionsKey, []string{})
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					helpers.ErrorResponse(log, w, "unauthorized", http.StatusUnauthorized)
					log.Error(fmt.Sprintf("Failed to refresh access token: %s", err.Error()))
					return
				}
				// Parse new tokens
				parsedAccessToken, err := authService.ParseToken(tokens.AccessToken)
				if err != nil || parsedAccessToken == nil {
					if allowUnauthorized {
						ctx = context.WithValue(ctx, UserIDKey, uuid.Nil)
						ctx = context.WithValue(ctx, UserRolesKey, []string{})
						ctx = context.WithValue(ctx, UserPermissionsKey, []string{})
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					helpers.ErrorResponse(log, w, "unauthorized", http.StatusUnauthorized)
					log.Error(fmt.Sprintf("Failed to parse access token: %s", err.Error()))
					return
				}
				parsedRefreshToken, err := authService.ParseToken(tokens.RefreshToken)
				if err != nil || parsedRefreshToken == nil {
					if allowUnauthorized {
						ctx = context.WithValue(ctx, UserIDKey, uuid.Nil)
						ctx = context.WithValue(ctx, UserRolesKey, []string{})
						ctx = context.WithValue(ctx, UserPermissionsKey, []string{})
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					helpers.ErrorResponse(log, w, "unauthorized", http.StatusUnauthorized)
					log.Error(fmt.Sprintf("Failed to parse refresh token: %s", err.Error()))
					return
				}
				// Update loaded JWT access
				jwtAccess = tokens.AccessToken
				// Save tokens to cookies
				http.SetCookie(w, &http.Cookie{
					Name:     "jwt_access",
					Value:    tokens.AccessToken,
					Path:     "/",
					Expires:  parsedAccessToken.RegisteredClaims.ExpiresAt.Time,
					HttpOnly: true,
					Secure:   authServiceConfig.CookieSecure,
				})
				log.Debug("Added JWT access to the cookies")
				http.SetCookie(w, &http.Cookie{
					Name:     "jwt_refresh",
					Value:    tokens.RefreshToken,
					Path:     "/",
					Expires:  parsedRefreshToken.RegisteredClaims.ExpiresAt.Time,
					HttpOnly: true,
					Secure:   authServiceConfig.CookieSecure,
				})
				log.Debug("Added JWT refresh to the cookies")
				http.SetCookie(w, &http.Cookie{
					Name:     "authorized",
					Value:    "true",
					Path:     "/",
					Expires:  parsedRefreshToken.RegisteredClaims.ExpiresAt.Time,
					HttpOnly: false,
					Secure:   authServiceConfig.CookieSecure,
				})
				log.Debug("Authorized value is set to true in cookies")
			}
			// Validate token
			claims, err := authService.ParseToken(jwtAccess)
			if err != nil || claims == nil {
				if allowUnauthorized {
					ctx := r.Context()
					ctx = context.WithValue(ctx, UserIDKey, uuid.Nil)
					ctx = context.WithValue(ctx, UserRolesKey, []string{})
					ctx = context.WithValue(ctx, UserPermissionsKey, []string{})
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				helpers.ErrorResponse(log, w, "invalid token", http.StatusUnauthorized)
				log.Error(fmt.Sprintf("Failed to validate access token: %s", err.Error()))
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
				if allowUnauthorized {
					ctx = context.WithValue(ctx, UserIDKey, uuid.Nil)
					ctx = context.WithValue(ctx, UserRolesKey, []string{})
					ctx = context.WithValue(ctx, UserPermissionsKey, []string{})
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				helpers.ErrorResponse(log, w, "failed to load user (by user ID from JWT access from cookies)", http.StatusUnauthorized)
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

func RequireRoles(log logger.Logger, all bool, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get roles from context
			userRoles, ok := r.Context().Value(UserRolesKey).([]string)
			if !ok {
				log.Error("Forbidden: failed to get roles from context")
				helpers.ErrorResponse(log, w, "forbidden", http.StatusForbidden)
				return
			}
			if all {
				// Check if user has all required roles
				for _, role := range requiredRoles {
					if !slices.Contains(userRoles, role) {
						log.Error("Forbidden: user must have all required roles")
						helpers.ErrorResponse(log, w, "forbidden", http.StatusForbidden)
						return
					}
				}
				next.ServeHTTP(w, r)
				return
			} else {
				// Check if user has at least one required role
				for _, userRole := range userRoles {
					if slices.Contains(requiredRoles, userRole) {
						next.ServeHTTP(w, r)
						return
					}
				}
				log.Error("Forbidden: user must have at least one required role")
				helpers.ErrorResponse(log, w, "forbidden", http.StatusForbidden)
			}
		})
	}
}

func RequirePermissions(log logger.Logger, all bool, requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user permissions from context
			userPermissions, ok := r.Context().Value(UserPermissionsKey).([]string)
			if !ok {
				log.Error("Forbidden: failed to get permissions from context")
				helpers.ErrorResponse(log, w, "forbidden", http.StatusForbidden)
				return
			}
			if all {
				// Check if user has all required permissions
				for _, permission := range requiredPermissions {
					if !slices.Contains(userPermissions, permission) {
						log.Error("Forbidden: user must have all required permissions")
						helpers.ErrorResponse(log, w, "forbidden", http.StatusForbidden)
						return
					}
				}
				next.ServeHTTP(w, r)
				return
			} else {
				// Check if user has at least one required permission
				for _, userPermission := range userPermissions {
					if slices.Contains(requiredPermissions, userPermission) {
						next.ServeHTTP(w, r)
						return
					}
				}
				log.Error("Forbidden: user must have at least one required permission")
				helpers.ErrorResponse(log, w, "forbidden", http.StatusForbidden)
			}
		})
	}
}
