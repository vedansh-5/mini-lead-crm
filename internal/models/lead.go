package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// custom type to ensure typesafety
type LeadStatus string

const (
	StatusNew       LeadStatus = "NEW"
	StatusContacted LeadStatus = "CONTACTED"
	StatusQualified LeadStatus = "QUALIFIED"
	StatusConverted LeadStatus = "CONVERTED"
	StatusLost      LeadStatus = "LOST"
)

type Lead struct {
	// gorm:"type:uuid;primary_key" tells GORM this is a UUID and the primary key
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name      string     `gorm:"not null" json:"name" binding:"required"` // binding req for gin validation
	Email     string     `gorm:"not null;unique" json:"email" binding:"required,email"`
	Phone     string     `json:"phone"`
	Status    LeadStatus `gorm:"type:varchar(20);defaut:'NEW'" json:"status"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// beforeCreate generates a new UUID for the lead

func (lead *Lead) BeforeCreate(tx *gorm.DB) (err error) {
	lead.ID = uuid.New()
	return
}

func IsValidTransition(current, next LeadStatus) bool {
	if next == StatusLost {
		return current != StatusConverted
	}

	switch current {
	case StatusNew:
		return next == StatusContacted
	case StatusContacted:
		return next == StatusQualified
	case StatusQualified:
		return next == StatusConverted
	default: // lost and converted can't be changed
		return false
	}
}
