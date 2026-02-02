package model

import (
	"github.com/google/uuid"
	"time"
)

// Student is a table (model), that contains info, related only to students.
// This table extends the "users" table
type Student struct {
	// one-to-one (student-to-user)
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	StudentGroupID uint16
	// many-to-many (student-to-parent)
	Parents *[]Parent `gorm:"many2many:parent_students;foreignKey:UserID;joinForeignKey:StudentID;references:UserID;joinReferences:ParentId"`
}
