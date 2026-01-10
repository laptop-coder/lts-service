package model

import (
	"github.com/google/uuid"
	"time"
)

// InstitutionAdministrator is a table (model), that contains info, related only
// to institution administrators. This table extends the "users" table
type InstitutionAdministrator struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// one-to-one (administrator-to-user)
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PositionID uint8
	// 1. Can't remove position if there are at least one person with it
	// 2. many-to-one (administrator-to-position)
	Position AdministratorPosition `gorm:"foreignKey:PositionID;constraint:OnDelete:restrict,OnUpdate:restrict"`
}
