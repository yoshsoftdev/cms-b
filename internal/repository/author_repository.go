package repository

import (
	"cms/internal/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorRepository struct{ db *gorm.DB }

func NewAuthorRepository(db *gorm.DB) *AuthorRepository { return &AuthorRepository{db: db} }

func (r *AuthorRepository) List(search string) ([]model.Author, error) {
	q := r.db.Model(&model.Author{})
	if s := strings.TrimSpace(search); s != "" {
		like := "%" + s + "%"
		q = q.Where("name LIKE ?", like)
	}
	var items []model.Author
	if err := q.Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *AuthorRepository) Create(input model.Author) (*model.Author, error) {
	if err := r.db.Create(&input).Error; err != nil {
		return nil, err
	}
	return &input, nil
}

func (r *AuthorRepository) Update(id string, input model.Author) (*model.Author, error) {
	var existing model.Author
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&existing).Updates(map[string]any{
		"name":     input.Name,
		"photo":    input.Photo,
		"alt_text": input.AltText,
	}).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *AuthorRepository) Delete(id string) error {
	if parsed, err := uuid.Parse(id); err != nil || parsed == uuid.Nil {
		return gorm.ErrRecordNotFound
	}
	var used int64
	if err := r.db.Table("contents").Where("author_id = ? AND deleted_at IS NULL", id).Count(&used).Error; err != nil {
		return err
	}
	if used > 0 {
		return gorm.ErrInvalidData
	}
	return r.db.Delete(&model.Author{}, "id = ?", id).Error
}
