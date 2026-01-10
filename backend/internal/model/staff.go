package model

import (
	"github.com/google/uuid"
	"time"
)

// Staff is a table (model), that contains info, related only to institution
// staff. This table extends the "users" table
type Staff struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// one-to-one (staff-to-user)
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PositionID uint8
	// 1. Can't remove position if there are at least one person with it
	// 2. many-to-one (staff-to-position)
	Position StaffPosition `gorm:"foreignKey:PositionID;constraint:OnDelete:restrict,OnUpdate:restrict"`
}
