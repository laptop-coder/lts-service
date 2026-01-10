package model

import (
	"time"
)

// Subject provides model of table with list of subjects ("Русский язык",
// "Литература", etc.)
type Subject struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(100);unique;check:length(trim(name)) >= 3"`
	// many-to-many (teacher-to-subject)
	Teachers []Teacher `gorm:"many2many:teacher_subjects"`
}
