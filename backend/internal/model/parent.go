package model

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
	"time"
)

// Parent is a table (model), that contains info, related only to parents.
// This table extends the "users" table
type Parent struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// one-to-one (parent-to-user)
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// many-to-many (parent-to-student)
	Students *[]Student `gorm:"many2many:parent_students"`
}

func AddConstraintsParentStudents(db *gorm.DB) error {
	return db.Exec(`
        ALTER TABLE parent_students 
			ADD CONSTRAINT fk_parent_students_parent
				FOREIGN KEY (parent_id) REFERENCES parents(user_id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT,
			ADD CONSTRAINT fk_parent_students_student
				FOREIGN KEY (student_id) REFERENCES students(user_id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT;
    `).Error
}
