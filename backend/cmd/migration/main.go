package main

import (
	"backend/internal/database"
	"backend/internal/model"
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
}
