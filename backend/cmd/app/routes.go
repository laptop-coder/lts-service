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
	teacherHandler *handler.TeacherHandler,
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
	mux.Handle("GET /api/v1/users/me", authMiddleware(requirePermissions(false, permissions.UserReadOwn)(http.HandlerFunc(userHandler.GetOwnUser))))
	// User roles
	mux.Handle("PUT /api/v1/users/{id}/roles", authMiddleware(requirePermissions(false, permissions.RoleAssign)(http.HandlerFunc(userHandler.AssignRoles))))
	mux.Handle("PATCH /api/v1/users/{id}/roles", authMiddleware(requirePermissions(false, permissions.RoleAdd)(http.HandlerFunc(userHandler.AddRoles))))
	mux.Handle("DELETE /api/v1/users/{userId}/roles/{roleId}", authMiddleware(requirePermissions(false, permissions.RoleDelete)(http.HandlerFunc(userHandler.RemoveRole))))
	mux.Handle("GET /api/v1/users/{id}/roles", authMiddleware(requirePermissions(false, permissions.RoleReadAny)(http.HandlerFunc(userHandler.GetRoles))))
	mux.Handle("GET /api/v1/users/me/roles", authMiddleware(requirePermissions(false, permissions.RoleReadOwn)(http.HandlerFunc(userHandler.GetOwnRoles))))
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
	mux.Handle("GET /api/v1/students/me", authMiddleware(requirePermissions(false, permissions.StudentReadOwn)(http.HandlerFunc(studentHandler.GetOwn))))
	// Teacher
	mux.Handle("GET /api/v1/teachers/{id}", authMiddleware(requirePermissions(false, permissions.TeacherReadOther)(http.HandlerFunc(teacherHandler.GetTeacherByID))))
	mux.Handle("GET /api/v1/teachers/me", authMiddleware(requirePermissions(false, permissions.TeacherReadOwn)(http.HandlerFunc(teacherHandler.GetOwn))))
	// Parents
	mux.Handle("GET /api/v1/parents/{id}", authMiddleware(requirePermissions(false, permissions.ParentReadOther)(http.HandlerFunc(parentHandler.GetParentByID))))
	mux.Handle("GET /api/v1/parents/me", authMiddleware(requirePermissions(false, permissions.ParentReadOwn)(http.HandlerFunc(parentHandler.GetOwn))))
	// Roles
	// mux.Handle("GET /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.GetPermissions)))
	// mux.Handle("PUT /api/v1/roles/{id}/permissions", authMiddleware(http.HandlerFunc(roleHandler.AssignPermissions)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}
