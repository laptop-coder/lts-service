package model

import (
	"github.com/google/uuid"
	"time"
)

// StudentGroup is a model of the table contains info about groups of students
// ("1А" grade, e.g.)
type StudentGroup struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(20);unique;check:length(trim(name)) >= 1"`

	GroupAdvisorID *uuid.UUID `gorm:"type:uuid;default:null"`
	// 1. Can't remove student group if there are at least one student in it. To
	// remove the group you need to reassign all students to another group at
	// first
	// 2. one-to-many (group-to-student)
	Students []Student `gorm:"foreignKey:StudentGroupID;references:ID;constraint:OnDelete:restrict,OnUpdate:restrict"`
}
