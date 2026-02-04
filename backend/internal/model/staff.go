package model

import (
	"github.com/google/uuid"
	"time"
)

// Staff is a table (model), that contains info, related only to institution
// staff. This table extends the "users" table
type Staff struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	// one-to-one (staff-to-user)
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	PositionID uint8
}
