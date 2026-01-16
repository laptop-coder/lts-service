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
	// 1. We consider that one group can have only one advisor (mentor, e.g.), but
	// one advisor can manage many groups, so:
	// many-to-one (group-to-advisor, i.e. group-to-user)
	// 2. Set GroupAdvisorID null in case of removing the user (i.e. the advisor)
	GroupAdvisor *User `gorm:"foreignKey:GroupAdvisorID;references:ID;constraint:OnDelete:set null,OnUpdate:restrict"`
	// one-to-many (group-to-student)
    Students []Student
}
