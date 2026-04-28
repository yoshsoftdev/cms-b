package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Author struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Photo     string    `json:"photo" gorm:"size:512"`
	AltText   string    `json:"alt_text" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Author) TableName() string { return "authors" }

func (a *Author) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
