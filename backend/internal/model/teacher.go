package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Teacher is a table (model), that contains info, related only to teachers.
// This table extends the "users" table
type Teacher struct {
	// one-to-one (teacher-to-user)
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// ClassroomID *uint8 `gorm:"default:null"` // TODO: remove if not necessary

	// 1. Classroom may not be specified
	// 2. Can't remove teacher if there are room, assigned to him. To remove the
	// teacher you need to unbind him or reassign to other classroom at first
	// (actual for schools)
	// 3. one-to-one (teacher-to-room)
	Classroom *Room `gorm:"foreignKey:TeacherID;references:UserID;constraint:OnDelete:restrict,OnUpdate:restrict"`
	// many-to-many (teacher-to-subject)
	Subjects []Subject `gorm:"many2many:teacher_subjects;foreignKey:UserID;joinForeignKey:TeacherId;references:ID;joinReferences:SubjectID"`
	// 1. This is a list of student groups for which this user is the advisor
	// (classroom teacher, for example)
	// 2. Set GroupAdvisorID null in case of removing the user (i.e. the
	// advisor)
	// 3. We consider that one group can have only one advisor (mentor, e.g.),
	// but one advisor can manage many groups, so: one-to-many (teacher-to-group,
	// i.e. advisor-to-group)
	StudentGroups *[]StudentGroup `gorm:"foreignKey:GroupAdvisorID;references:UserID"`
}

func AddConstraintsTeacherSubjects(db *gorm.DB) error {
	return db.Exec(`
        ALTER TABLE teacher_subjects 
		    DROP CONSTRAINT IF EXISTS fk_teacher_subjects_teacher,
		    DROP CONSTRAINT IF EXISTS fk_teacher_subjects_subject,
			ADD CONSTRAINT fk_teacher_subjects_teacher
				FOREIGN KEY (teacher_id) REFERENCES teachers(user_id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT,
			ADD CONSTRAINT fk_teacher_subjects_subject
				FOREIGN KEY (subject_id) REFERENCES subjects(id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT;
    `).Error
}
