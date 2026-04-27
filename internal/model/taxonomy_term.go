package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxonomyType string

const (
	TaxonomyKeyword  TaxonomyType = "keyword"
	TaxonomyCategory TaxonomyType = "category"
)

type TaxonomyTerm struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Type      TaxonomyType   `json:"type" gorm:"size:30;uniqueIndex:idx_taxonomy_type_name,priority:1;not null"`
	Name      string         `json:"name" gorm:"size:120;uniqueIndex:idx_taxonomy_type_name,priority:2;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (t *TaxonomyTerm) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
