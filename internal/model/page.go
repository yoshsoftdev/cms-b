package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Page struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Title           string         `json:"title" gorm:"size:255;not null"`
	Slug            string         `json:"slug" gorm:"size:255;uniqueIndex;not null"`
	Content         string         `json:"content" gorm:"type:text;not null"`
	Excerpt         string         `json:"excerpt" gorm:"type:text"`
	CoverImage      string         `json:"cover_image" gorm:"size:500"`
	Status          ContentStatus  `json:"status" gorm:"size:20;index;not null;default:draft"`
	MetaTitle       string         `json:"meta_title" gorm:"size:255"`
	MetaDescription string         `json:"meta_description" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Page) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
