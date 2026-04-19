package model

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Who sent a message (author or requester)
	SenderID uuid.UUID `gorm:"type:uuid;not null;index"`
	Sender   User      `gorm:"foreignKey:SenderID;constraint:OnDelete:cascade"`

	Content   string `gorm:"type:text;not null"`
	IsRead    bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
