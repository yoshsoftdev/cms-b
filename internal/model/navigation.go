package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Navigation struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Label     string         `json:"label" gorm:"size:255;not null"`
	URL       string         `json:"url" gorm:"size:500;not null"`
	ParentID  *uuid.UUID     `json:"parent_id" gorm:"type:uuid;index"`
	Order     int            `json:"order" gorm:"not null;default:0"`
	Type      string         `json:"type" gorm:"size:20;not null;default:internal"`
	Target    string         `json:"target" gorm:"size:20;not null;default:_self"`
	IsActive  bool           `json:"is_active" gorm:"not null;default:true"`
	Location  string         `json:"location" gorm:"size:20;index;not null;default:header"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (n *Navigation) BeforeCreate(_ *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
