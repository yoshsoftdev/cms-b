package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NavigationTranslation struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	NavigationID uuid.UUID  `json:"navigation_id" gorm:"type:uuid;not null;uniqueIndex:idx_navigation_locale,priority:1"`
	Locale       LocaleCode `json:"locale" gorm:"size:5;not null;uniqueIndex:idx_navigation_locale,priority:2"`
	Label        string     `json:"label" gorm:"size:255;not null"`
	URL          string     `json:"url" gorm:"size:500;not null"`
}

func (nt *NavigationTranslation) BeforeCreate(_ *gorm.DB) error {
	if nt.ID == uuid.Nil {
		nt.ID = uuid.New()
	}
	return nil
}
