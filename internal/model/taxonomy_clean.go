package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxonomyTypeClean string

const (
	TaxonomyCleanCategory TaxonomyTypeClean = "category"
	TaxonomyCleanKeyword  TaxonomyTypeClean = "keyword"
	TaxonomyCleanTag      TaxonomyTypeClean = "tag"
)

type Taxonomy struct {
	ID        uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey"`
	Type      TaxonomyTypeClean `json:"type" gorm:"size:20;not null;index"`
	CreatedAt time.Time         `json:"created_at"`
}

func (t *Taxonomy) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TaxonomyTranslation struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	TaxonomyID uuid.UUID  `json:"taxonomy_id" gorm:"type:uuid;not null;uniqueIndex:idx_taxonomy_locale,priority:1"`
	Locale     LocaleCode `json:"locale" gorm:"size:5;not null;uniqueIndex:idx_taxonomy_locale,priority:2"`
	Name       string     `json:"name" gorm:"size:120;not null"`
}

func (tt *TaxonomyTranslation) BeforeCreate(_ *gorm.DB) error {
	if tt.ID == uuid.Nil {
		tt.ID = uuid.New()
	}
	return nil
}

type ContentTaxonomy struct {
	ContentID  uuid.UUID `json:"content_id" gorm:"type:uuid;primaryKey"`
	TaxonomyID uuid.UUID `json:"taxonomy_id" gorm:"type:uuid;primaryKey"`
}

func (ContentTaxonomy) TableName() string { return "content_taxonomies" }
