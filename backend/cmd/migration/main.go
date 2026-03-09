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
		{ID: 2, Name: permissions.PostReadAny},
		{ID: 3, Name: permissions.PostReadOwn},
		{ID: 4, Name: permissions.PostUpdateAny},
		{ID: 5, Name: permissions.PostUpdateOwn},
		{ID: 6, Name: permissions.PostDeleteAny},
		{ID: 7, Name: permissions.PostDeleteOwn},
		{ID: 8, Name: permissions.PostPhotoDeleteAny},
		{ID: 9, Name: permissions.PostPhotoDeleteOwn},
		{ID: 10, Name: permissions.PostVerify},
		{ID: 11, Name: permissions.PostMarkReturnedAny},
		{ID: 12, Name: permissions.PostMarkReturnedOwn},
		{ID: 13, Name: permissions.UserReadOwn},
		{ID: 14, Name: permissions.UserReadOther},
		{ID: 15, Name: permissions.UserReadAll},
		{ID: 16, Name: permissions.UserUpdateOwn},
		{ID: 17, Name: permissions.UserDeleteAny},
		{ID: 18, Name: permissions.UserDeleteOwn},
		{ID: 19, Name: permissions.UserReadOwnSubjects},
		{ID: 20, Name: permissions.TeacherReadClassroom},
		{ID: 21, Name: permissions.StudentReadClassroom},
		{ID: 22, Name: permissions.TeacherStudentsRead},
		{ID: 23, Name: permissions.ParentStudentsRead},
		{ID: 24, Name: permissions.StudentTeacherRead},
		{ID: 25, Name: permissions.StudentParentsRead},
		{ID: 26, Name: permissions.RoomCreate},
		{ID: 27, Name: permissions.RoomRead},
		{ID: 28, Name: permissions.RoomUpdate},
		{ID: 29, Name: permissions.RoomDelete},
		{ID: 30, Name: permissions.SubjectCreate},
		{ID: 31, Name: permissions.SubjectRead},
		{ID: 32, Name: permissions.SubjectUpdate},
		{ID: 33, Name: permissions.SubjectDelete},
		{ID: 34, Name: permissions.StudentGroupCreate},
		{ID: 35, Name: permissions.StudentGroupReadAny},
		{ID: 36, Name: permissions.StudentGroupReadOwn},
		{ID: 37, Name: permissions.StudentGroupUpdate},
		{ID: 38, Name: permissions.StudentGroupDelete},
		{ID: 39, Name: permissions.StudentGroupAdvisorAssign},
		{ID: 40, Name: permissions.StudentGroupAdvisorUnassign},
		{ID: 41, Name: permissions.StudentGroupAdvisorRead},
		{ID: 42, Name: permissions.TeacherSubjectReadAny},
		{ID: 43, Name: permissions.TeacherSubjectReadOwn},
		{ID: 44, Name: permissions.TeacherSubjectAssignAny},
		{ID: 45, Name: permissions.TeacherSubjectAssignOwn},
		{ID: 46, Name: permissions.TeacherSubjectUnassignAny},
		{ID: 47, Name: permissions.TeacherSubjectUnassignOwn},
		{ID: 48, Name: permissions.TeacherClassroomReadAny},
		{ID: 49, Name: permissions.TeacherClassroomReadOwn},
		{ID: 50, Name: permissions.TeacherClassroomAssignAny},
		{ID: 51, Name: permissions.TeacherClassroomAssignOwn},
		{ID: 52, Name: permissions.TeacherClassroomUnassignAny},
		{ID: 53, Name: permissions.TeacherClassroomUnassignOwn},
		{ID: 54, Name: permissions.ParentStudentReadAny},
		{ID: 55, Name: permissions.ParentStudentReadOwn},
		{ID: 56, Name: permissions.ParentStudentAssignAny},
		{ID: 57, Name: permissions.ParentStudentAssignOwn},
		{ID: 58, Name: permissions.ParentStudentUnassignAny},
		{ID: 59, Name: permissions.ParentStudentUnassignOwn},
		{ID: 60, Name: permissions.ParentReadOther},
		{ID: 61, Name: permissions.RoleAssign},
		{ID: 62, Name: permissions.RoleAdd},
		{ID: 63, Name: permissions.RoleDelete},
		{ID: 64, Name: permissions.RoleReadAny},
		{ID: 65, Name: permissions.RoleReadOwn},
		{ID: 66, Name: permissions.TokenInviteAdminCreate},
		{ID: 67, Name: permissions.TokenInviteUserCreate},
		{ID: 68, Name: permissions.TokenInviteAdminDelete},
		{ID: 69, Name: permissions.TokenInviteUserDelete},
		{ID: 70, Name: permissions.StudentReadOther},
	}
	for _, permission := range permissions {
		if err := db.FirstOrCreate(&permission, model.Permission{ID: permission.ID}).Error; err != nil {
			log.Error("MIGRATION | Failed to create permissions", err)
			panic(err)
		}
	}
	log.Info("MIGRATION | Permissions created successfully")
}
