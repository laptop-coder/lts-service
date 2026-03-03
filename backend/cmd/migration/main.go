package main

import (
	"backend/internal/database"
	"backend/internal/model"
	"backend/internal/permissions"
	"backend/pkg/env"
	"backend/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	// Logger
	log := logger.New()
	log.Info("MIGRATION | Starting...")

	// Database
	log.Info("MIGRATION | Initializing database...")
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
		log.Error("MIGRATION | Cannot initialize database")
		panic("Cannot initialize database")
	}
	defer database.Close(db)
	log.Info("MIGRATION | Database connected successfully")

	if err := database.Migrate(
		db,
		[]any{
			&model.Role{},
			&model.Permission{},
			&model.AdministratorPosition{},
			&model.StaffPosition{},
			&model.Room{},
			&model.Subject{},
			&model.User{},
			&model.InstitutionAdministrator{},
			&model.Staff{},
			&model.Teacher{},
			&model.StudentGroup{},
			&model.Student{},
			&model.Parent{},
			&model.Post{},
		},
	); err != nil {
		log.Error("MIGRATION | Cannot make migration", err)
		panic(err)
	}

	// Add constraints (e.g., foreign key ON DELETE)
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := model.AddConstraintsRolePermissions(tx); err != nil {
			return err
		}
		if err := model.AddConstraintsUserRoles(tx); err != nil {
			return err
		}
		if err := model.AddConstraintsTeacherSubjects(tx); err != nil {
			return err
		}
		if err := model.AddConstraintsParentStudents(tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Error("MIGRATION | Cannot add constraints")
		panic("Cannot add constraints")
	}

	// Create roles
	roles := []model.Role{
		{ID: 1, Name: "superadmin"},
		{ID: 2, Name: "admin"},
		{ID: 3, Name: "institution_administrator"},
		{ID: 4, Name: "staff"},
		{ID: 5, Name: "teacher"},
		{ID: 6, Name: "parent"},
		{ID: 7, Name: "student"},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{ID: role.ID}).Error; err != nil {
			log.Error("MIGRATION | Failed to create roles", err)
			panic(err)
		}
	}
	log.Info("MIGRATION | Roles created successfully")

	// Create permissions
	permissions := []model.Permission{
		{ID: 1, Name: permissions.PostCreate},
		{ID: 2, Name: permissions.PostRead},
		{ID: 3, Name: permissions.PostUpdate},
		{ID: 4, Name: permissions.PostUpdateOwn},
		{ID: 5, Name: permissions.PostDelete},
		{ID: 6, Name: permissions.PostDeleteOwn},
		{ID: 7, Name: permissions.PostVerify},
		{ID: 8, Name: permissions.UserCreate},
		{ID: 9, Name: permissions.UserRead},
		{ID: 10, Name: permissions.UserUpdate},
		{ID: 11, Name: permissions.UserUpdateOwn},
		{ID: 12, Name: permissions.UserDelete},
		{ID: 13, Name: permissions.UserDeleteOwn},
		{ID: 14, Name: permissions.RoomCreate},
		{ID: 15, Name: permissions.RoomRead},
		{ID: 16, Name: permissions.RoomUpdate},
		{ID: 17, Name: permissions.RoomDelete},
		{ID: 18, Name: permissions.SubjectCreate},
		{ID: 19, Name: permissions.SubjectRead},
		{ID: 20, Name: permissions.SubjectUpdate},
		{ID: 21, Name: permissions.SubjectDelete},
		{ID: 22, Name: permissions.GroupCreate},
		{ID: 23, Name: permissions.GroupRead},
		{ID: 24, Name: permissions.GroupUpdate},
		{ID: 25, Name: permissions.GroupDelete},
		{ID: 26, Name: permissions.GroupAssignAdvisor},
		{ID: 27, Name: permissions.TeacherAssignSubject},
		{ID: 28, Name: permissions.TeacherAssignClassroom},
		{ID: 29, Name: permissions.StudentAssignParent},
		{ID: 30, Name: permissions.ParentViewStudents},
		{ID: 31, Name: permissions.RoleAssign},
		{ID: 32, Name: permissions.RoleCreate},
		{ID: 33, Name: permissions.RoleUpdate},
		{ID: 34, Name: permissions.RoleDelete},
	}
	for _, permission := range permissions {
		if err := db.FirstOrCreate(&permission, model.Permission{ID: permission.ID}).Error; err != nil {
			log.Error("MIGRATION | Failed to create permissions", err)
			panic(err)
		}
	}
	log.Info("MIGRATION | Permissions created successfully")
}
