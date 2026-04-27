package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentStatus string

const (
	StatusDraft     ContentStatus = "draft"
	StatusPublished ContentStatus = "published"
)

type Post struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Title           string         `json:"title" gorm:"size:255;not null"`
	TitleID         string         `json:"title_id" gorm:"size:255"`
	Slug            string         `json:"slug" gorm:"size:255;uniqueIndex;not null"`
	Content         string         `json:"content" gorm:"type:text;not null"`
	ContentID       string         `json:"content_id" gorm:"type:text"`
	ContentJSON     string         `json:"content_json" gorm:"type:text"`
	ContentJSONID   string         `json:"content_json_id" gorm:"type:text"`
	Excerpt         string         `json:"excerpt" gorm:"type:text"`
	CoverImage      string         `json:"cover_image" gorm:"size:500"`
	Status          ContentStatus  `json:"status" gorm:"size:20;index;not null;default:draft"`
	MetaTitle       string         `json:"meta_title" gorm:"size:255"`
	MetaTitleID     string         `json:"meta_title_id" gorm:"size:255"`
	MetaDescription string         `json:"meta_description" gorm:"type:text"`
	MetaDescriptionID string       `json:"meta_description_id" gorm:"type:text"`
	OGImage         string         `json:"og_image" gorm:"size:500"`
	CanonicalURL    string         `json:"canonical_url" gorm:"size:500"`
	Keywords        string         `json:"keywords" gorm:"type:text"`
	KeywordsID      string         `json:"keywords_id" gorm:"type:text"`
	Category        string         `json:"category" gorm:"size:100"`
	PublishDate     string         `json:"publish_date" gorm:"size:40"`
	Author          string         `json:"author" gorm:"size:255"`
	ReadTime        string         `json:"read_time" gorm:"size:40"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Post) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
