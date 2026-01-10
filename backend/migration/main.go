// Package migration provides function to make migrations.
package migration

import (
	"gorm.io/gorm"
	"backend/internal/database"
	"backend/internal/model"
	log "backend/pkg/logger"
	"fmt"
)

func main(){
	db := database.SetUpDatabase()
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
