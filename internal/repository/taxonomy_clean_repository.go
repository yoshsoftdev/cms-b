package repository

import (
	"strings"

	"gorm.io/gorm"
)

type TaxonomyCleanRepository struct{ db *gorm.DB }

func NewTaxonomyCleanRepository(db *gorm.DB) *TaxonomyCleanRepository {
	return &TaxonomyCleanRepository{db: db}
}

type TaxonomyCleanItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (r *TaxonomyCleanRepository) List(locale, termType, search string, limit int) ([]TaxonomyCleanItem, error) {
	if locale == "" {
		locale = "en"
	}
	if limit <= 0 {
		limit = 20
	}
	q := r.db.Table("taxonomies t").
		Select("t.id, t.type, tt.name").
		Joins("JOIN taxonomy_translations tt ON tt.taxonomy_id = t.id AND tt.locale = ?", locale)
	if termType != "" {
		q = q.Where("t.type = ?", termType)
	}
	if strings.TrimSpace(search) != "" {
		q = q.Where("tt.name LIKE ?", "%"+strings.TrimSpace(search)+"%")
	}
	var out []TaxonomyCleanItem
	if err := q.Order("tt.name ASC").Limit(limit).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
