package model

import (
	"github.com/google/uuid"
	"time"
)

// Student is a table (model), that contains info, related only to students.
// This table extends the "users" table
type Student struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// one-to-one (student-to-user)
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`

	CreatedAt time.Time
	UpdatedAt time.Time

	StudentGroupID uint8
    // 1. Can't remove student group if there are at least one student in it. To 
    // remove the group you need to reassign all students to another group at
    // first
	// 2. many-to-one (student-to-group)
	StudentGroup StudentGroup `gorm:"foreignKey:StudentGroupID;constraint:OnDelete:restrict,OnUpdate:restrict"`
	// many-to-many (student-to-parent)
	Parents *[]Parent `gorm:"many2many:parent_students"`
}
