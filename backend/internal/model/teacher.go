package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Teacher is a table (model), that contains info, related only to teachers.
// This table extends the "users" table
type Teacher struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// one-to-one (teacher-to-user)
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// 1. Classroom may not be specified
	ClassroomID *uint8 `gorm:"default:null"`
	// 2. Can't remove room if there are at least one teacher, assigned to it. To
	// remove the room you need to reassign the teacher to other classroom at
	// first (actual for schools)
	// 3. one-to-one (teacher-to-room)
	Classroom *Room `gorm:"foreignKey:ClassroomID;constraint:OnDelete:restrict,OnUpdate:restrict"`
	// many-to-many (teacher-to-subject)
	Subjects []Subject `gorm:"many2many:teacher_subjects"`
}

func AddConstraintsTeacherSubjects(db *gorm.DB) error {
	return db.Exec(`
        ALTER TABLE teacher_subjects 
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
