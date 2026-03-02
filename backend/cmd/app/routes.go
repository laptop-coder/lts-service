package main

import (
	"backend/internal/handler"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

func SetupRoutes(
	mux *http.ServeMux,
	log logger.Logger,
	authMiddleware func(http.Handler) http.Handler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	studentGroupHandler *handler.StudentGroupHandler,
	roomHandler *handler.RoomHandler,
	subjectHandler *handler.SubjectHandler,
	studentHandler *handler.StudentHandler,
	parentHandler *handler.ParentHandler,
) {
	// Public routes (no auth required)
	mux.HandleFunc("POST /api/v1/users", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/health", healthHandler)

	// Secure routes (auth required)

	// User
	mux.Handle("PATCH /api/v1/users/{id}", authMiddleware(http.HandlerFunc(userHandler.UpdateProfile)))
	mux.Handle("DELETE /api/v1/users/{id}/avatar", authMiddleware(http.HandlerFunc(userHandler.RemoveAvatar)))
	mux.Handle("PUT /api/v1/users/{id}/avatar", authMiddleware(http.HandlerFunc(userHandler.UpdateAvatar)))
	mux.Handle("GET /api/v1/users/{id}", authMiddleware(http.HandlerFunc(userHandler.GetUserByID)))
	// Student groups
	mux.Handle("GET /api/v1/student_groups/{id}/advisor", authMiddleware(http.HandlerFunc(studentGroupHandler.GetAdvisorByGroupID)))
	// Auth
	mux.Handle("DELETE /api/v1/users/{id}", authMiddleware(http.HandlerFunc(authHandler.DeleteAccount)))
	mux.Handle("POST /api/v1/auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))
	// Posts
	mux.Handle("POST /api/v1/posts", authMiddleware(http.HandlerFunc(postHandler.Create)))
	mux.Handle("DELETE /api/v1/posts/{id}", authMiddleware(http.HandlerFunc(postHandler.Delete)))
	mux.Handle("DELETE /api/v1/posts/{id}/photo", authMiddleware(http.HandlerFunc(postHandler.RemovePhoto)))
	mux.Handle("PATCH /api/v1/posts/{id}", authMiddleware(http.HandlerFunc(postHandler.Update)))
	// Rooms
	mux.Handle("POST /api/v1/rooms", authMiddleware(http.HandlerFunc(roomHandler.Create)))
	mux.Handle("DELETE /api/v1/rooms/{id}", authMiddleware(http.HandlerFunc(roomHandler.Delete)))
	mux.Handle("PATCH /api/v1/rooms/{id}", authMiddleware(http.HandlerFunc(roomHandler.Update)))
	// Subjects
	mux.Handle("POST /api/v1/subjects", authMiddleware(http.HandlerFunc(subjectHandler.Create)))
	mux.Handle("DELETE /api/v1/subjects/{id}", authMiddleware(http.HandlerFunc(subjectHandler.Delete)))
	mux.Handle("PATCH /api/v1/subjects/{id}", authMiddleware(http.HandlerFunc(subjectHandler.Update)))
	// Students
	mux.Handle("GET /api/v1/students/{id}", authMiddleware(http.HandlerFunc(studentHandler.GetStudentByID)))
	// Parents
	mux.Handle("GET /api/v1/parents/{id}", authMiddleware(http.HandlerFunc(parentHandler.GetParentByID)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}
