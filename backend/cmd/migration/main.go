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
			&model.InstitutionAdministratorPosition{},
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
		{ID: 19, Name: permissions.RoomCreate},
		{ID: 21, Name: permissions.RoomUpdate},
		{ID: 22, Name: permissions.RoomDelete},
		{ID: 23, Name: permissions.SubjectCreate},
		{ID: 25, Name: permissions.SubjectUpdate},
		{ID: 26, Name: permissions.SubjectDelete},
		{ID: 27, Name: permissions.StudentGroupCreate},
		{ID: 29, Name: permissions.StudentGroupUpdate},
		{ID: 30, Name: permissions.StudentGroupDelete},
		{ID: 31, Name: permissions.StudentGroupAdvisorAssign},
		{ID: 32, Name: permissions.StudentGroupAdvisorUnassignAny},
		{ID: 33, Name: permissions.StudentGroupAdvisorUnassignOwn},
		{ID: 34, Name: permissions.StudentGroupAdvisorRead},
		{ID: 35, Name: permissions.TeacherSubjectReadAny},
		{ID: 36, Name: permissions.TeacherSubjectReadOwn},
		{ID: 37, Name: permissions.TeacherSubjectAddAny},
		{ID: 38, Name: permissions.TeacherSubjectAddOwn},
		{ID: 39, Name: permissions.TeacherSubjectAssignAny},
		{ID: 40, Name: permissions.TeacherSubjectAssignOwn},
		{ID: 41, Name: permissions.TeacherSubjectUnassignAny},
		{ID: 42, Name: permissions.TeacherSubjectUnassignOwn},
		{ID: 43, Name: permissions.TeacherClassroomReadAny},
		{ID: 44, Name: permissions.TeacherClassroomReadOwn},
		{ID: 45, Name: permissions.TeacherClassroomAssignAny},
		{ID: 46, Name: permissions.TeacherClassroomAssignOwn},
		{ID: 47, Name: permissions.TeacherClassroomUnassignAny},
		{ID: 48, Name: permissions.TeacherClassroomUnassignOwn},
		{ID: 49, Name: permissions.TeacherReadOther},
		{ID: 50, Name: permissions.TeacherReadOwn},
		{ID: 51, Name: permissions.TeacherStudentGroupReadOwn},
		{ID: 52, Name: permissions.ParentStudentReadAny},
		{ID: 53, Name: permissions.ParentStudentReadOwn},
		{ID: 54, Name: permissions.ParentStudentAddAny},
		{ID: 55, Name: permissions.ParentStudentAddOwn},
		{ID: 56, Name: permissions.ParentStudentUnassignAny},
		{ID: 57, Name: permissions.ParentStudentUnassignOwn},
		{ID: 58, Name: permissions.ParentReadOther},
		{ID: 59, Name: permissions.ParentReadOwn},
		{ID: 60, Name: permissions.ParentStudentGroupReadOwn},
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
		{ID: 71, Name: permissions.StudentReadOwn},
		{ID: 72, Name: permissions.StudentClassroomReadAny},
		{ID: 73, Name: permissions.StudentClassroomReadOwn},
		{ID: 74, Name: permissions.StudentAdvisorReadAny},
		{ID: 75, Name: permissions.StudentAdvisorReadOwn},
		{ID: 76, Name: permissions.StudentParentReadAny},
		{ID: 77, Name: permissions.StudentParentReadOwn},
		{ID: 78, Name: permissions.StudentStudentGroupReadOwn},
		{ID: 79, Name: permissions.InstitutionAdministratorReadOther},
		{ID: 80, Name: permissions.InstitutionAdministratorReadOwn},
		{ID: 81, Name: permissions.InstitutionAdministratorPositionAssign},
		{ID: 82, Name: permissions.InstitutionAdministratorPositionRead},
		{ID: 83, Name: permissions.StaffReadOther},
		{ID: 84, Name: permissions.StaffReadOwn},
		{ID: 85, Name: permissions.StaffPositionAssign},
		{ID: 86, Name: permissions.StaffPositionRead},
		{ID: 87, Name: permissions.PositionInstitutionAdministratorCreate},
		{ID: 89, Name: permissions.PositionInstitutionAdministratorUpdate},
		{ID: 90, Name: permissions.PositionInstitutionAdministratorDelete},
		{ID: 91, Name: permissions.PositionStaffCreate},
		{ID: 93, Name: permissions.PositionStaffUpdate},
		{ID: 94, Name: permissions.PositionStaffDelete},
	}
	for _, permission := range permissions {
		if err := db.FirstOrCreate(&permission, model.Permission{ID: permission.ID}).Error; err != nil {
			log.Error("MIGRATION | Failed to create permissions", err)
			panic(err)
		}
	}
	log.Info("MIGRATION | Permissions created successfully")

	// Lists of permissions by roles
	superadminPermissions := []string{
		"token.invite.admin.create", "token.invite.admin.delete",
	}

	adminPermissions := []string{
		"post.create", "post.read.any", "post.read.own", "post.update.any", "post.update.own", "post.delete.any", "post.delete.own", "post.photo.delete.any", "post.photo.delete.own", "post.verify", "post.mark.returned.any", "post.mark.returned.own", "user.read.own", "user.read.other", "user.read.all", "user.update.own", "user.delete.any", "user.delete.own", "room.create", "room.update", "room.delete", "subject.create", "subject.update", "subject.delete", "student_group.create", "student_group.update", "student_group.delete", "student_group.advisor.assign", "student_group.advisor.unassign.any", "student_group.advisor.read", "teacher.subject.read.any", "teacher.subject.add.any", "teacher.subject.assign.any", "teacher.subject.unassign.any", "teacher.classroom.read.any", "teacher.classroom.assign.any", "teacher.classroom.unassign.any", "teacher.read.other", "parent.student.read.any", "parent.student.add.any", "parent.student.unassign.any", "parent.read.other", "role.assign", "role.add", "role.delete", "role.read.any", "role.read.own", "token.invite.user.create", "token.invite.user.delete", "student.read.other", "student.classroom.read.any", "student.advisor.read.any", "student.parent.read.any", "institution_administrator.read.other", "institution_administrator.position.assign", "institution_administrator.position.read", "staff.read.other", "staff.position.assign", "staff.position.read", "position.institution_administrator.create", "position.institution_administrator.update", "position.institution_administrator.delete", "position.staff.create", "position.staff.update", "position.staff.delete",
	}

	institutionAdministratorPermissions := []string{
		"post.create", "post.read.own", "post.update.own", "post.delete.own", "post.photo.delete.own", "post.mark.returned.own", "user.read.own", "user.read.other", "user.update.own", "user.delete.own", "student_group.advisor.read", "teacher.subject.read.any", "teacher.classroom.read.any", "teacher.read.other", "parent.read.other", "role.read.any", "role.read.own", "student.read.other", "student.classroom.read.any", "student.advisor.read.any", "student.parent.read.any", "institution_administrator.read.other", "institution_administrator.read.own", "institution_administrator.position.read", "staff.read.other", "staff.position.read",
	}

	staffPermissions := []string{
		"post.create", "post.read.own", "post.update.own", "post.delete.own", "post.photo.delete.own", "post.mark.returned.own", "user.read.own", "user.read.other", "user.update.own", "user.delete.own", "student_group.advisor.read", "teacher.subject.read.any", "teacher.classroom.read.any", "teacher.read.other", "parent.read.other", "role.read.any", "role.read.own", "student.read.other", "student.classroom.read.any", "student.advisor.read.any", "student.parent.read.any", "institution_administrator.read.other", "institution_administrator.position.read", "staff.read.other", "staff.read.own", "staff.position.read",
	}

	teacherPermissions := []string{
		"post.create", "post.read.own", "post.update.own", "post.delete.own", "post.photo.delete.own", "post.mark.returned.own", "user.read.own", "user.read.other", "user.update.own", "user.delete.own", "student_group.advisor.assign", "student_group.advisor.unassign.own", "student_group.advisor.read", "teacher.subject.read.any", "teacher.subject.read.own", "teacher.subject.add.own", "teacher.subject.assign.own", "teacher.subject.unassign.own", "teacher.classroom.read.any", "teacher.classroom.read.own", "teacher.classroom.assign.own", "teacher.classroom.unassign.own", "teacher.read.other", "teacher.read.own", "parent.read.other", "role.read.any", "role.read.own", "student.read.other", "student.classroom.read.any", "student.advisor.read.any", "student.parent.read.any", "institution_administrator.read.other", "institution_administrator.position.read", "staff.read.other", "staff.position.read", "teacher.student_group.read.own",
	}

	parentPermissions := []string{
		"post.create", "post.read.own", "post.update.own", "post.delete.own", "post.photo.delete.own", "post.mark.returned.own", "user.read.own", "user.read.other", "user.update.own", "user.delete.own", "student_group.advisor.read", "teacher.subject.read.any", "teacher.classroom.read.any", "teacher.read.other", "parent.student.read.own", "parent.student.add.own", "parent.student.unassign.own", "parent.read.other", "parent.read.own", "role.read.any", "role.read.own", "student.read.other", "student.classroom.read.any", "student.advisor.read.any", "student.parent.read.any", "institution_administrator.read.other", "institution_administrator.position.read", "staff.read.other", "staff.position.read", "parent.student_group.read.own",
	}

	studentPermissions := []string{
		"post.create", "post.read.own", "post.update.own", "post.delete.own", "post.photo.delete.own", "post.mark.returned.own", "user.read.own", "user.read.other", "user.update.own", "user.delete.own", "student_group.advisor.read", "teacher.subject.read.any", "teacher.classroom.read.any", "teacher.read.other", "parent.read.other", "role.read.any", "role.read.own", "student.read.other", "student.read.own", "student.classroom.read.any", "student.advisor.read.any", "student.advisor.read.own", "student.parent.read.any", "student.parent.read.own", "student.classroom.read.own", "institution_administrator.read.other", "institution_administrator.position.read", "staff.read.other", "staff.position.read", "student.student_group.read.own",
	}

	// Assign permissions to roles
	for _, name := range superadminPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 1, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to superadmin role", err)
			panic(err)
		}
	}

	for _, name := range adminPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 2, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to admin role", err)
			panic(err)
		}
	}

	for _, name := range institutionAdministratorPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 3, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to institution administrator role", err)
			panic(err)
		}
	}

	for _, name := range staffPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 4, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to staff role", err)
			panic(err)
		}
	}

	for _, name := range teacherPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 5, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to teacher role", err)
			panic(err)
		}
	}

	for _, name := range parentPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 6, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to parent role", err)
			panic(err)
		}
	}

	for _, name := range studentPermissions {
		if err := db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT 7, id FROM permissions WHERE name = ?
		ON CONFLICT DO NOTHING;
		`, name).Error; err != nil {
			log.Error("MIGRATION | Failed to assign permissions to student role", err)
			panic(err)
		}
	}
}
