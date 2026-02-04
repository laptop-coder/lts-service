// Package model provides types (models) for using with the GORM.
package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	// UUID4; maybe it will be used in URIs to see users' profiles, e.g., or
	// their posts.
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// Email will be used as login
	Email string `gorm:"type:varchar(320);unique;check:length(trim(email)) >= 5 and email like '_%@_%._%'"`
	// 60 bytes is the bcrypt hash length
	Password   string  `gorm:"type:varchar(60);check:length(trim(password)) = 60"`
	FirstName  string  `gorm:"type:varchar(100);check:length(trim(first_name)) >= 2"`
	MiddleName *string `gorm:"type:varchar(100);default:null;check:(middle_name is null) or (length(trim(middle_name)) >= 2)"`
	LastName   string  `gorm:"type:varchar(100);check:length(trim(last_name)) >= 2"`
	// Name of the avatar is ID of the user. It is located in the root of the
	// site (public dir) and saved in the JPEG format: "/<user_id>.jpeg"
	HasAvatar bool `gorm:"type:boolean;default:false"`
	// many-to-many (user-to-role)
	Roles []Role `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoleID"`
	// one-to-one (user-to-administrator)
	InstitutionAdministrator *InstitutionAdministrator `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`
	// one-to-one (user-to-staff)
	Staff *Staff `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`
	// one-to-one (user-to-teacher)
	Teacher *Teacher `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`
	// one-to-one (user-to-parent)
	Parent *Parent `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`
	// one-to-one (user-to-student)
	Student *Student `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`
	// 1. This is a list of student groups for which this user is the advisor
	// (classroom teacher, for example)
	// 2. Set GroupAdvisorID null in case of removing the user (i.e. the
	// advisor)
	// 3. We consider that one group can have only one advisor (mentor, e.g.),
	// but one advisor can manage many groups, so: one-to-many (user-to-group,
	// i.e. advisor-to-group)
	StudentGroups *[]StudentGroup `gorm:"foreignKey:GroupAdvisorID;references:ID;constraint:OnDelete:set null,OnUpdate:restrict"`
	// 1. Can be null if the user hasn't created any posts yet
	// 2. Removing of the user will cause removing all of his posts
	// 3. one-to-many (author-to-post, i.e. user-to-post)
	Posts *[]Post `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`
}

func AddConstraintsUserRoles(db *gorm.DB) error {
	return db.Exec(`
        ALTER TABLE user_roles 
		    DROP CONSTRAINT IF EXISTS fk_user_roles_user,
		    DROP CONSTRAINT IF EXISTS fk_user_roles_role,
			ADD CONSTRAINT fk_user_roles_user
				FOREIGN KEY (user_id) REFERENCES users(id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT,
			ADD CONSTRAINT fk_user_roles_role
				FOREIGN KEY (role_id) REFERENCES roles(id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT;
    `).Error
}
