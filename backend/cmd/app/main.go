// Package main is the entrypoint of the backend app.
package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/valkey"
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

	// Valkey
	log.Info("Initializing Valkey...")
	jwtClient := valkey.NewClient(valkey.ClientDBs.JWT, log)
	defer valkey.Close(jwtClient)
	log.Info("Valkey client(-s) connected successfully")

	// Repositories
	log.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(db, log)
	jwtRepo := repository.NewJWTRepository(jwtClient, log)
	studentGroupRepo := repository.NewStudentGroupRepository(db, log)
	postRepo := repository.NewPostRepository(db, log)
	roomRepo := repository.NewRoomRepository(db, log)
	subjectRepo := repository.NewSubjectRepository(db, log)
	studentRepo := repository.NewStudentRepository(db, log)
	parentRepo := repository.NewParentRepository(db, log)

	// Services
	log.Info("Creating service configurations...")
	serviceConfigs := config.NewServiceConfigs(sharedConfig)
	log.Info("Initializing services...")
	authService := service.NewAuthService(userRepo, jwtRepo, db, serviceConfigs.Auth, log)
	userService := service.NewUserService(userRepo, studentRepo, roomRepo, db, serviceConfigs.User, log)
	postService := service.NewPostService(postRepo, db, serviceConfigs.Post, log)
	studentGroupService := service.NewStudentGroupService(studentGroupRepo, db, log)
	roomService := service.NewRoomService(roomRepo, db, log)
	subjectService := service.NewSubjectService(subjectRepo, db, log)
	studentService := service.NewStudentService(studentRepo, userRepo, db, log)
	parentService := service.NewParentService(parentRepo, userRepo, db, log)

	// Handlers
	log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService, userService, serviceConfigs.Auth, log)
	userHandler := handler.NewUserHandler(userService, log)
	postHandler := handler.NewPostHandler(postService, log)
	studentGroupHandler := handler.NewStudentGroupHandler(studentGroupService, log)
	roomHandler := handler.NewRoomHandler(roomService, log)
	subjectHandler := handler.NewSubjectHandler(subjectService, log)
	studentHandler := handler.NewStudentHandler(studentService, log)
	parentHandler := handler.NewParentHandler(parentService, log)

	mux := http.NewServeMux()
	authMiddleware := middleware.Auth(authService, db)

	SetupRoutes(
		mux,
		log,
		authMiddleware,
		authHandler,
		userHandler,
		postHandler,
		studentGroupHandler,
		roomHandler,
		subjectHandler,
		studentHandler,
		parentHandler,
	)

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
