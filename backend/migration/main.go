package migration

import (
	"gorm.io/gorm"
	"backend/internal/database"
	"backend/internal/model"
	"backend/pkg/logger"
	"fmt"
	"backend/pkg/env"
)

func main(){
	// Logger
	log := logger.New()
	log.Info("Starting application...")

	// Database
	log.Info("Initializing database...")
	db, err := database.Connect(
		database.Config{
			DBName:   env.GetStringRequired("POSTGRES_DB"),
			Host:     env.GetStringRequired("POSTGRES_HOST"),
			Password: env.GetStringRequired("POSTGRES_PASSWORD"),
			Port:     env.GetIntRequired("POSTGRES_PORT"),
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

	if err := db.AutoMigrate(
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
	); err != nil {
		msg := fmt.Sprintf(
				"cannot make migration: %s", err.Error(),
			)
		log.Error(msg)
		panic(msg)
	}

	// Add constraints (e.g., foreign key ON DELETE)
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := model.AddConstraintsRolePermissions(tx); err != nil {
			return err
		}
		if err:=model.AddConstraintsUserRoles(tx); err != nil {
			return err
		}
		if err:=model.AddConstraintsTeacherSubjects(tx); err != nil {
			return err
		}
		if err := model.AddConstraintsParentStudents(tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		msg := "cannot add constraints"
		log.Error(msg)
		panic(msg)
	}
}
