package repository

import (
	"cms/internal/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaRepository struct{ db *gorm.DB }

func NewMediaRepository(db *gorm.DB) *MediaRepository { return &MediaRepository{db: db} }

func (r *MediaRepository) List(search string) ([]model.Media, error) {
	q := r.db.Model(&model.Media{})
	if s := strings.TrimSpace(search); s != "" {
		like := "%" + s + "%"
		q = q.Where("name LIKE ? OR alt_text LIKE ? OR url LIKE ?", like, like, like)
	}
	var items []model.Media
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MediaRepository) Create(input model.Media) (*model.Media, error) {
	if err := r.db.Create(&input).Error; err != nil {
		return nil, err
	}
	return &input, nil
}

func (r *MediaRepository) Update(id string, input model.Media) (*model.Media, error) {
	var existing model.Media
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&existing).Updates(map[string]any{
		"name":     input.Name,
		"url":      input.URL,
		"alt_text": input.AltText,
	}).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *MediaRepository) Delete(id string) error {
	if parsed, err := uuid.Parse(id); err != nil || parsed == uuid.Nil {
		return gorm.ErrRecordNotFound
	}
	return r.db.Delete(&model.Media{}, "id = ?", id).Error
}
