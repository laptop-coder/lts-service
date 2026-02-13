// Package main is the entrypoint of the backend app.
package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/env"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Logger
	log := logger.New()
	log.Info("Starting application...")

	// Configs
	log.Info("Loading configurations...")
	appConfig := config.LoadAppConfig()
	sharedConfig := config.LoadSharedConfig()

	// Database
	log.Info("Initializing database...")
	db, err := database.Connect(
		database.Config{
			DBName:   env.GetStringRequired("POSTGRES_DB"),
			Host:     env.GetStringRequired("POSTGRES_HOST"),
			Password: env.GetStringRequired("POSTGRES_PASSWORD"),
			Port:     5432,
			SSLMode: func() string {
				if env.GetBoolRequired("POSTGRES_SSL_MODE") {
					return "enable"
				}
				return "disable"
			}(),
			TimeZone: env.GetStringRequired("POSTGRES_TIME_ZONE"),
			User:     env.GetStringRequired("POSTGRES_USER"),
		},
	)
	if err != nil {
		log.Error("Cannot initialize database")
		panic("Cannot initialize database")
	}
	defer database.Close(db)
	log.Info("Database connected successfully")

	// Repositories
	log.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(db, log)
	studentGroupRepo := repository.NewStudentGroupRepository(db, log)

	// Services
	log.Info("Creating service configurations...")
	serviceConfigs := config.NewServiceConfigs(sharedConfig)

	log.Info("Initializing services...")
	userService := service.NewUserService(
		userRepo,
		db,
		serviceConfigs.User,
		log,
	)
	studentGroupService := service.NewStudentGroupService(
		studentGroupRepo,
		db,
		log,
	)

	// Handlers
	log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(userService, log)
	userHandler := handler.NewUserHandler(userService, log)
	studentGroupHandler := handler.NewStudentGroupHandler(studentGroupService, log)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/users", authHandler.Register)
	mux.HandleFunc("DELETE /api/v1/users/{id}", authHandler.DeleteAccount)
	mux.HandleFunc("PATCH /api/v1/users/{id}", userHandler.UpdateProfile)
	mux.HandleFunc("DELETE /api/v1/users/{id}/avatar", userHandler.RemoveAvatar)
	mux.HandleFunc("PUT /api/v1/users/{id}/avatar", userHandler.UpdateAvatar)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("GET /api/v1/student_groups/{id}/advisor", studentGroupHandler.GetAdvisorByGroupID)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
	})

	// Middleware
	var handler http.Handler = mux
	handler = middleware.Logging(log, handler)

	// Server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(appConfig.Port),
		Handler: handler,
	}

	go func() {
		log.Info("Starting server...", "port", strconv.Itoa(appConfig.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited properly")
}
