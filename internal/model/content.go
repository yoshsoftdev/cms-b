package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentType string

const (
	ContentTypePost ContentType = "post"
	ContentTypePage ContentType = "page"
)

type ContentRecord struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Type      ContentType    `json:"type" gorm:"size:20;index;not null"`
	Slug      string         `json:"slug" gorm:"size:255;uniqueIndex;not null"`
	Status    ContentStatus  `json:"status" gorm:"size:20;index;not null;default:draft"`
	AuthorID  *uuid.UUID     `json:"author_id" gorm:"type:uuid;index"`
	PublishAt *time.Time     `json:"publish_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ContentRecord) TableName() string { return "contents" }

func (c *ContentRecord) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
