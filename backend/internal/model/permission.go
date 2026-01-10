package model

import (
	"time"
)

type Permission struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(150);unique;check:length(trim(name)) >= 6"`
	// many-to-many (permission-to-role)
	Roles []Role `gorm:"many2many:role_permissions"`
}
