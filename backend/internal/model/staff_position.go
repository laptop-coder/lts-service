package model

import (
	"time"
)

type StaffPosition struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(100);unique;check:length(trim(name)) >= 4"`
	// 1. Can't remove position if there are at least one person with it
	// 2. one-to-many (position-to-staff)
	Staff []Staff `gorm:"foreignKey:PositionID;references:ID;constraint:OnDelete:restrict,OnUpdate:restrict"`
}
