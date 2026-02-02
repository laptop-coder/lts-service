package model

import (
	"github.com/google/uuid"
	"time"
)

// Post contains a model of a table with info about lost things in the format of
// posts
type Post struct {
	// UUID4; it will be used in URIs to see statuses of the posts
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(50);check:length(trim(name)) >= 2"`
	Description      string `gorm:"type:varchar(1000)"`
	// was the post verified by moderator (user with the role of service
	// administrator)? (true/false)
	Verified bool `gorm:"type:boolean;default:false"`
	// was the thing found, i.e. returned to owner? (true/false)
	ThingReturnedToOwner bool `gorm:"type:boolean;default:false;check:(thing_returned_to_owner=true and verified=true) or thing_returned_to_owner=false"`
	// the logic is the same as for user's avatar
	HasPhoto bool `gorm:"type:boolean;default:false"`
	AuthorID uuid.UUID `gorm:"type:uuid"`
}
