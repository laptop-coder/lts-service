package model

import (
	"github.com/google/uuid"
	"time"
)

// InstitutionAdministrator is a table (model), that contains info, related only
// to institution administrators. This table extends the "users" table
type InstitutionAdministrator struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	// one-to-one (administrator-to-user)
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	PositionID uint8
}
