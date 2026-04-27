package repository

import (
	"strings"

	"cms/internal/model"

	"gorm.io/gorm"
)

type TaxonomyRepository struct{ db *gorm.DB }

func NewTaxonomyRepository(db *gorm.DB) *TaxonomyRepository { return &TaxonomyRepository{db: db} }

func (r *TaxonomyRepository) List(termType model.TaxonomyType, search string, limit int) ([]model.TaxonomyTerm, error) {
	q := r.db.Model(&model.TaxonomyTerm{}).Where("type = ?", termType)
	if strings.TrimSpace(search) != "" {
		q = q.Where("name LIKE ?", "%"+strings.TrimSpace(search)+"%")
	}
	if limit <= 0 {
		limit = 20
	}
	var items []model.TaxonomyTerm
	err := q.Order("name asc").Limit(limit).Find(&items).Error
	return items, err
}

func (r *TaxonomyRepository) UpsertNames(termType model.TaxonomyType, names []string) error {
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		var existing model.TaxonomyTerm
		err := r.db.Where("type = ? AND name = ?", termType, name).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		term := model.TaxonomyTerm{Type: termType, Name: name}
		if err := r.db.Create(&term).Error; err != nil {
			return err
		}
	}
	return nil
}
