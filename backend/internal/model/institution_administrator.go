package model

import (
	"github.com/google/uuid"
	"time"
)

// InstitutionAdministrator is a table (model), that contains info, related only
// to institution administrators. This table extends the "users" table
type InstitutionAdministrator struct {
	// one-to-one (administrator-to-user)
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PositionID uint16
	Position   InstitutionAdministratorPosition `gorm:"foreignKey:PositionID;references:ID"`
}
