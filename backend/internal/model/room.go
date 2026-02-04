package model

import (
	"github.com/google/uuid"
	"time"
)

// Room provides model of table with list of rooms (cabinets, dining room, etc.)
type Room struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(20);unique;check:length(trim(name)) >= 1"`
	// one-to-one (room-to-teacher)
	TeacherID uuid.UUID `gorm:"type:uuid"`
}
