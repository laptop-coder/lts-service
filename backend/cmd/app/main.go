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
	"backend/pkg/imghash"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	valkeyGo "github.com/valkey-io/valkey-go"
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
	appMode := config.ParseAppMode(env.GetStringRequired("APP_MODE"))
	appConfig := config.LoadAppConfig(appMode)
	sharedConfig := config.LoadSharedConfig()

	// Create directories
	log.Info("Creating directories (if not exist)...")
	if err := os.MkdirAll(sharedConfig.Storage.Avatar.UploadPath, 0755); err != nil {
		log.Error("failed to create upload directory for avatars", "error", err.Error())
		panic(fmt.Errorf("failed to create upload directory for avatars: %w", err))
	}
	log.Info("Avatars upload directory... OK")
	if err := os.MkdirAll(sharedConfig.Storage.PostPhoto.UploadPath, 0755); err != nil {
		log.Error("failed to create upload directory for post photos", "error", err.Error())
		panic(fmt.Errorf("failed to create upload directory for post photos: %w", err))
	}
	log.Info("Post photos upload directory... OK")
	if err := os.MkdirAll(sharedConfig.Storage.Avatar.DeletePath, 0755); err != nil {
		log.Error("failed to create delete directory (trash) for avatars", "error", err.Error())
		panic(fmt.Errorf("failed to create delete directory (trash) for avatars: %w", err))
	}
	log.Info("Avatars trash... OK")
	if err := os.MkdirAll(sharedConfig.Storage.PostPhoto.DeletePath, 0755); err != nil {
		log.Error("failed to create delete directory (trash) for post photos", "error", err.Error())
		panic(fmt.Errorf("failed to create delete directory (trash) for post photos: %w", err))
	}
	log.Info("Post photos trash... OK")

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
			AppMode:  appConfig.AppMode,
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
	businessClient := valkey.NewClient(valkey.ClientDBs.Business, log)
	defer valkey.Close(businessClient)
	log.Info("Valkey client(-s) connected successfully")

	// Packages (pkg)
	log.Info("Initializing packages (pkg)...")
	hashCalc := imghash.NewHashCalculator()

	// Repositories
	log.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(db, log)
	jwtRepo := repository.NewJWTRepository(jwtClient, log)
	studentGroupRepo := repository.NewStudentGroupRepository(db, log)
	postRepo := repository.NewPostRepository(db, businessClient, log)
	postModerationRepo := repository.NewPostModerationRepository(db, log)
	msgRepo := repository.NewMessageRepository(db, log)
	convRepo := repository.NewConversationRepository(db, log)
	roomRepo := repository.NewRoomRepository(db, log)
	subjectRepo := repository.NewSubjectRepository(db, log)
	studentRepo := repository.NewStudentRepository(db, log)
	teacherRepo := repository.NewTeacherRepository(db, log)
	parentRepo := repository.NewParentRepository(db, log)
	staffRepo := repository.NewStaffRepository(db, log)
	roleRepo := repository.NewRoleRepository(db, log)
	institutionAdministratorRepo := repository.NewInstitutionAdministratorRepository(db, log)
	institutionAdministratorPositionRepo := repository.NewInstitutionAdministratorPositionRepository(db, log)
	staffPositionRepo := repository.NewStaffPositionRepository(db, log)

	// Services
	log.Info("Creating service configurations...")
	serviceConfigs := config.NewServiceConfigs(sharedConfig, appConfig)
	log.Info("Initializing services...")
	emailService, err := service.NewEmailService(serviceConfigs.Email, log)
	if err != nil {
		panic(err)
	}
	authService := service.NewAuthService(emailService, userRepo, jwtRepo, db, businessClient, serviceConfigs.Auth, log)
	userService := service.NewUserService(userRepo, studentRepo, roomRepo, db, serviceConfigs.User, log)
	postService := service.NewPostService(postRepo, postModerationRepo, userService, hashCalc, db, businessClient, serviceConfigs.Post, log)
	conversationService := service.NewConversationService(convRepo, msgRepo, postRepo, userRepo, emailService, db, log)
	studentGroupService := service.NewStudentGroupService(userRepo, studentGroupRepo, db, log)
	roomService := service.NewRoomService(roomRepo, db, log)
	subjectService := service.NewSubjectService(subjectRepo, db, log)
	studentService := service.NewStudentService(studentRepo, studentGroupRepo, userRepo, teacherRepo, db, log)
	teacherService := service.NewTeacherService(teacherRepo, userRepo, db, log)
	parentService := service.NewParentService(parentRepo, userRepo, db, log)
	staffService := service.NewStaffService(staffRepo, userRepo, db, log)
	institutionAdministratorService := service.NewInstitutionAdministratorService(institutionAdministratorRepo, userRepo, db, log)
	inviteService := service.NewInviteService(emailService, jwtRepo, userRepo, roleRepo, db, serviceConfigs.Invite, log)
	institutionAdministratorPositionService := service.NewInstitutionAdministratorPositionService(institutionAdministratorPositionRepo, db, log)
	staffPositionService := service.NewStaffPositionService(staffPositionRepo, db, log)

	// Handlers
	log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService, userService, inviteService, serviceConfigs.Auth, log)
	userHandler := handler.NewUserHandler(userService, log)
	postHandler := handler.NewPostHandler(postService, userService, teacherService, parentService, studentGroupService, studentService, log)
	conversationHandler := handler.NewConversationHandler(conversationService, log)
	studentGroupHandler := handler.NewStudentGroupHandler(teacherService, studentGroupService, log)
	roomHandler := handler.NewRoomHandler(roomService, log)
	subjectHandler := handler.NewSubjectHandler(subjectService, log)
	studentHandler := handler.NewStudentHandler(studentService, log)
	teacherHandler := handler.NewTeacherHandler(teacherService, log)
	parentHandler := handler.NewParentHandler(parentService, log)
	staffHandler := handler.NewStaffHandler(staffService, log)
	institutionAdministratorHandler := handler.NewInstitutionAdministratorHandler(institutionAdministratorService, log)
	inviteHandler := handler.NewInviteHandler(inviteService, serviceConfigs.Invite, log)
	institutionAdministratorPositionHandler := handler.NewInstitutionAdministratorPositionHandler(institutionAdministratorPositionService, log)
	staffPositionHandler := handler.NewStaffPositionHandler(staffPositionService, log)
	documentHandler := handler.NewDocumentHandler(serviceConfigs.Document, log)

	mux := http.NewServeMux()
	authMiddleware := func(allowUnauthorized bool) func(http.Handler) http.Handler {
		return middleware.Auth(context.Background(), authService, serviceConfigs.Auth, jwtRepo, db, log, allowUnauthorized)
	}
	requireRoles := middleware.RequireRoles
	requirePermissions := middleware.RequirePermissions

	SetupRoutes(
		mux,
		log,
		authMiddleware,
		requireRoles,
		requirePermissions,
		authHandler,
		userHandler,
		postHandler,
		conversationHandler,
		studentGroupHandler,
		roomHandler,
		subjectHandler,
		studentHandler,
		teacherHandler,
		parentHandler,
		staffHandler,
		institutionAdministratorHandler,
		inviteHandler,
		staffPositionHandler,
		institutionAdministratorPositionHandler,
		documentHandler,
	)

	// Middleware
	var handler http.Handler = mux
	handler = middleware.Logging(log, handler)
	handler = middleware.CORS(handler)

	// Server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(appConfig.Port),
		Handler: handler,
	}
	go func() {
		log.Info("Starting server...", "port", strconv.Itoa(appConfig.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err.Error())
			panic(err)
		}
	}()

	// Moderation worker
	go func() {
		counter := 0
		log.Info("starting moderation worker...")
		for {
			// Check if moderation worker exists
			if _, err := userService.GetPostsModeratorBot(context.Background()); err != nil {
				log.Error(
					"failed to check posts moderator bot existence, maybe it does not exist. Waiting for 30 seconds to check one more time...",
					"error",
					err.Error(),
				)
				time.Sleep(30 * time.Second)
				continue
			}
			// Get post ID. If this is the 64th request, take ID from the DB
			// instead of queue
			var postID uuid.UUID
			if counter == 64 {
				counter = 0
				post, err := postService.GetOldestPendingPost(context.Background())
				if err != nil {
					log.Error("moderation worker: failed to get the oldest post with pending moderation status", "error", err.Error())
					continue
				}
				if post == nil {
					log.Error("moderation worker: post is nil")
					continue
				}
				postID = post.ID
			} else {
				// Else get post ID from the queue
				res := businessClient.Do(
					context.Background(),
					businessClient.
						B().
						Brpop().
						Key("moderation:posts:queue").
						Timeout(5).
						Build(),
				)
				if errors.Is(res.Error(), valkeyGo.Nil) {
					log.Info("moderation worker: there are no posts in queue to moderate. Looking at the database...")
					if err := postService.ModerateAllPosts(context.Background()); err != nil {
						log.Error("moderation worker: failed to moderate all posts", "error", err.Error())
						continue
					}
					log.Info("moderation worker: all posts were moderated. Waiting for 30 seconds to check one more time...")
					time.Sleep(30 * time.Second)
					continue
				}
				if res.Error() != nil {
					log.Error("moderation worker error: failed to get post to moderate from queue. Waiting for 5 seconds to retry...", "error", res.Error().Error())
					time.Sleep(5 * time.Second)
					continue
				}
				arr, err := res.AsStrSlice()
				if err != nil {
					log.Error("moderation worker error: failed to represent response as array. Waiting for 5 seconds to retry...")
					time.Sleep(5 * time.Second)
					continue
				}
				postIDParsed, err := uuid.Parse(arr[1])
				if err != nil {
					log.Error("moderation worker error: cannot convert post id to uuid. Waiting for 5 seconds to retry...")
					time.Sleep(5 * time.Second)
					continue
				}
				postID = postIDParsed
			}
			// Send post to the moderation service and change status in the DB
			if err := postService.ModeratePost(context.Background(), postID); err != nil {
				log.Error("moderation worker error: failed to moderate post. Waiting for 5 seconds to retry...", "error", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			// Increment counter
			counter += 1
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
		log.Error("Server forced to shutdown", "error", err.Error())
	}
	log.Info("Server exited properly")
}
