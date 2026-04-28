package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Media struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"size:255;not null;default:''"`
	URL       string    `json:"url" gorm:"size:500;not null"`
	AltText   string    `json:"alt_text" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Media) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type ContentMedia struct {
	ContentID uuid.UUID `json:"content_id" gorm:"type:uuid;primaryKey"`
	MediaID   uuid.UUID `json:"media_id" gorm:"type:uuid;primaryKey"`
	Type      string    `json:"type" gorm:"size:50;primaryKey"`
}

func (ContentMedia) TableName() string { return "content_media" }
