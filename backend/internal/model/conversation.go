package model

import (
	"github.com/google/uuid"
	"time"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool      `gorm:"default:true"`
	Messages  []Message `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:cascade,OnUpdate:restrict"`

	// Post
	PostID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_conversation_unique"`
	Post   Post      `gorm:"foreignKey:PostID"`

	// Participants
	RequesterID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_conversation_unique"` // user who pressed button "contact"
	Requester   User      `gorm:"foreignKey:RequesterID"`

	AuthorID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_conversation_unique"` // post author
	Author   User      `gorm:"foreignKey:AuthorID"`
}
