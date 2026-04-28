package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocaleCode string

const (
	LocaleEN LocaleCode = "en"
	LocaleID LocaleCode = "id"
)

type ContentTranslation struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ContentID       uuid.UUID  `json:"content_id" gorm:"type:uuid;not null;uniqueIndex:idx_content_locale,priority:1"`
	Locale          LocaleCode `json:"locale" gorm:"size:5;not null;uniqueIndex:idx_content_locale,priority:2"`
	Title           string     `json:"title" gorm:"size:255;not null"`
	Content         string     `json:"content" gorm:"type:text"`
	ContentJSON     string     `json:"content_json" gorm:"type:text"`
	Excerpt         string     `json:"excerpt" gorm:"type:text"`
	MetaTitle       string     `json:"meta_title" gorm:"size:255"`
	MetaDescription string     `json:"meta_description" gorm:"type:text"`
	Keywords        string     `json:"keywords" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (ct *ContentTranslation) BeforeCreate(_ *gorm.DB) error {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	return nil
}
