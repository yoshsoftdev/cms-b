package repository

import (
	"cms/internal/model"

	"gorm.io/gorm"
)

type PageRepository struct{ db *gorm.DB }

func NewPageRepository(db *gorm.DB) *PageRepository { return &PageRepository{db: db} }

func (r *PageRepository) Create(page *model.Page) error { return r.db.Create(page).Error }
func (r *PageRepository) Update(page *model.Page) error { return r.db.Save(page).Error }
func (r *PageRepository) Delete(id string) error {
	return r.db.Delete(&model.Page{}, "id = ?", id).Error
}

func (r *PageRepository) GetByID(id string) (*model.Page, error) {
	var page model.Page
	if err := r.db.First(&page, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &page, nil
}
func (r *PageRepository) GetBySlug(slug string) (*model.Page, error) {
	var page model.Page
	if err := r.db.First(&page, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &page, nil
}
func (r *PageRepository) ExistsBySlug(slug, excludeID string) (bool, error) {
	q := r.db.Model(&model.Page{}).Where("slug = ?", slug)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *PageRepository) List(filter ContentFilter) ([]model.Page, int64, error) {
	q := r.db.Model(&model.Page{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		q = q.Where("title LIKE ?", "%"+filter.Search+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	var items []model.Page
	err := q.Order("created_at desc").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Find(&items).Error
	return items, total, err
}
