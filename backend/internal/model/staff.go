package model

import (
	"github.com/google/uuid"
	"time"
)

// Staff is a table (model), that contains info, related only to institution
// staff. This table extends the "users" table
type Staff struct {
	// one-to-one (staff-to-user)
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PositionID uint16
	Position   StaffPosition `gorm:"foreignKey:PositionID;references:ID"`
}
