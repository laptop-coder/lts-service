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
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// many-to-one (student-to-group)
	StudentGroupID uint16
	StudentGroup   StudentGroup `gorm:"foreignKey:StudentGroupID;references:ID"`
	// many-to-many (student-to-parent)
	Parents *[]Parent `gorm:"many2many:parent_students;foreignKey:UserID;joinForeignKey:StudentID;references:UserID;joinReferences:ParentId"`
}
