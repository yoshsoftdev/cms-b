package repository

import (
	"cms/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct{ db *gorm.DB }

func NewPostRepository(db *gorm.DB) *PostRepository { return &PostRepository{db: db} }

func (r *PostRepository) Create(post *model.Post) error { return r.db.Create(post).Error }
func (r *PostRepository) Update(post *model.Post) error { return r.db.Save(post).Error }
func (r *PostRepository) Delete(id string) error {
	return r.db.Delete(&model.Post{}, "id = ?", id).Error
}

func (r *PostRepository) GetByID(id string) (*model.Post, error) {
	var post model.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
func (r *PostRepository) GetBySlug(slug string) (*model.Post, error) {
	var post model.Post
	if err := r.db.First(&post, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
func (r *PostRepository) ExistsBySlug(slug, excludeID string) (bool, error) {
	q := r.db.Model(&model.Post{}).Where("slug = ?", slug)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *PostRepository) List(filter ContentFilter) ([]model.Post, int64, error) {
	q := r.db.Model(&model.Post{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		q = q.Where("title LIKE ? OR title_id LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
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
	var items []model.Post
	err := q.Order("created_at desc").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Find(&items).Error
	return items, total, err
}
