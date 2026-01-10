package model

import (
	"time"
)

type StaffPosition struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(100);unique;check:length(trim(name)) >= 4"`
	// one-to-many (position-to-staff)
	Staff []Staff
}
