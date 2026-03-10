package main

import (
	"backend/internal/handler"
	"backend/internal/permissions"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

func SetupRoutes(
	mux *http.ServeMux,
	log logger.Logger,
	authMiddleware func(http.Handler) http.Handler,
	requireRoles func(bool, ...string) func(http.Handler) http.Handler,
	requirePermissions func(bool, ...string) func(http.Handler) http.Handler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	studentGroupHandler *handler.StudentGroupHandler,
	roomHandler *handler.RoomHandler,
	subjectHandler *handler.SubjectHandler,
	studentHandler *handler.StudentHandler,
	parentHandler *handler.ParentHandler,
	roleHandler *handler.RoleHandler,
) {
	// Public routes (no auth required)
	mux.HandleFunc("POST /api/v1/users", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/posts/public", postHandler.GetPostsPublic)
	mux.HandleFunc("/health", healthHandler)

	// Secure routes (auth required)

	// User
	mux.Handle("PATCH /api/v1/users/me", authMiddleware(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnProfile))))
	mux.Handle("PUT /api/v1/users/me/avatar", authMiddleware(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.UpdateOwnAvatar))))
	mux.Handle("DELETE /api/v1/users/me/avatar", authMiddleware(requirePermissions(false, permissions.UserUpdateOwn)(http.HandlerFunc(userHandler.RemoveOwnAvatar))))
	mux.Handle("GET /api/v1/users/{id}", authMiddleware(requirePermissions(false, permissions.UserReadOther)(http.HandlerFunc(userHandler.GetUserByID))))
	mux.Handle("GET /api/v1/users", authMiddleware(requirePermissions(false, permissions.UserReadAll)(http.HandlerFunc(userHandler.GetUsers))))
	// Student groups
	mux.Handle("GET /api/v1/student_groups/{id}/advisor", authMiddleware(requirePermissions(false, permissions.StudentGroupAdvisorRead)(http.HandlerFunc(studentGroupHandler.GetAdvisorByGroupID))))
	// Auth
	mux.Handle("DELETE /api/v1/users/{id}", authMiddleware(requirePermissions(false, permissions.UserDeleteAny)(http.HandlerFunc(authHandler.DeleteAccount))))
	mux.Handle("DELETE /api/v1/users/me", authMiddleware(requirePermissions(false, permissions.UserDeleteOwn)(http.HandlerFunc(authHandler.DeleteOwnAccount))))
	mux.Handle("POST /api/v1/auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))
	// Posts
	mux.Handle("POST /api/v1/posts", authMiddleware(requirePermissions(false, permissions.PostCreate)(http.HandlerFunc(postHandler.Create))))
	mux.Handle("DELETE /api/v1/posts/{id}", authMiddleware(requirePermissions(false, permissions.PostDeleteAny, permissions.PostDeleteOwn)(http.HandlerFunc(postHandler.Delete))))
	mux.Handle("DELETE /api/v1/posts/{id}/photo", authMiddleware(requirePermissions(false, permissions.PostPhotoDeleteAny, permissions.PostPhotoDeleteOwn)(http.HandlerFunc(postHandler.RemovePhoto))))
	mux.Handle("PATCH /api/v1/posts/{id}", authMiddleware(requirePermissions(false, permissions.PostUpdateAny, permissions.PostUpdateOwn)(http.HandlerFunc(postHandler.Update))))
	mux.Handle("GET /api/v1/posts", authMiddleware(requirePermissions(false, permissions.PostReadAny)(http.HandlerFunc(postHandler.GetPosts))))
	mux.Handle("GET /api/v1/users/me/posts", authMiddleware(requirePermissions(false, permissions.PostReadOwn)(http.HandlerFunc(postHandler.GetOwnPosts))))
	mux.Handle("PATCH /api/v1/posts/{id}/verify", authMiddleware(requirePermissions(false, permissions.PostVerify)(http.HandlerFunc(postHandler.Verify))))
	mux.Handle("PATCH /api/v1/posts/{id}/return", authMiddleware(requirePermissions(false, permissions.PostMarkReturnedAny, permissions.PostMarkReturnedOwn)(http.HandlerFunc(postHandler.ReturnToOwner))))
	// Rooms
	mux.Handle("POST /api/v1/rooms", authMiddleware(requirePermissions(false, permissions.RoomCreate)(http.HandlerFunc(roomHandler.Create))))
	mux.Handle("DELETE /api/v1/rooms/{id}", authMiddleware(requirePermissions(false, permissions.RoomDelete)(http.HandlerFunc(roomHandler.Delete))))
	mux.Handle("PATCH /api/v1/rooms/{id}", authMiddleware(requirePermissions(false, permissions.RoomUpdate)(http.HandlerFunc(roomHandler.Update))))
	mux.Handle("GET /api/v1/rooms", authMiddleware(requirePermissions(false, permissions.RoomRead)(http.HandlerFunc(roomHandler.GetRooms))))
	// Subjects
	mux.Handle("POST /api/v1/subjects", authMiddleware(requirePermissions(false, permissions.SubjectCreate)(http.HandlerFunc(subjectHandler.Create))))
	mux.Handle("DELETE /api/v1/subjects/{id}", authMiddleware(requirePermissions(false, permissions.SubjectDelete)(http.HandlerFunc(subjectHandler.Delete))))
	mux.Handle("PATCH /api/v1/subjects/{id}", authMiddleware(requirePermissions(false, permissions.SubjectUpdate)(http.HandlerFunc(subjectHandler.Update))))
	mux.Handle("GET /api/v1/subjects", authMiddleware(requirePermissions(false, permissions.SubjectRead)(http.HandlerFunc(subjectHandler.GetSubjects))))
	// Students
	mux.Handle("GET /api/v1/students/{id}", authMiddleware(requirePermissions(false, permissions.StudentReadOther)(http.HandlerFunc(studentHandler.GetStudentByID))))
	// Parents
	mux.Handle("GET /api/v1/parents/{id}", authMiddleware(requirePermissions(false, permissions.ParentReadOther)(http.HandlerFunc(parentHandler.GetParentByID))))
	// Roles
	// mux.Handle("GET /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.GetPermissions)))
	// mux.Handle("PUT /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.AssignPermissions)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}
