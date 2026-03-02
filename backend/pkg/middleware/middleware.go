package middleware

import (
	"backend/internal/service"
	"backend/pkg/helpers"
	"backend/pkg/logger"
	"context"
	"net/http"
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

func Auth(authService service.AuthService) func(http.Handler) http.Handler {
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
			// Put user ID and roles to the context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRolesKey, claims.Roles)
			//
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
