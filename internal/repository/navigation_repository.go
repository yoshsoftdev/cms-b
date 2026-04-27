package repository

import (
	"cms/internal/model"

	"gorm.io/gorm"
)

type NavigationRepository interface {
	Create(nav *model.Navigation) error
	Update(nav *model.Navigation) error
	Delete(id string) error
	GetByID(id string) (*model.Navigation, error)
	List(location string, onlyActive bool) ([]model.Navigation, error)
	Reorder(location string, ids []string) error
}

type GormNavigationRepository struct {
	db *gorm.DB
}

func NewGormNavigationRepository(db *gorm.DB) *GormNavigationRepository {
	return &GormNavigationRepository{db: db}
}

func (r *GormNavigationRepository) Create(nav *model.Navigation) error {
	return r.db.Create(nav).Error
}

func (r *GormNavigationRepository) Update(nav *model.Navigation) error {
	return r.db.Save(nav).Error
}

func (r *GormNavigationRepository) Delete(id string) error {
	return r.db.Delete(&model.Navigation{}, "id = ?", id).Error
}

func (r *GormNavigationRepository) GetByID(id string) (*model.Navigation, error) {
	var nav model.Navigation
	if err := r.db.First(&nav, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &nav, nil
}

func (r *GormNavigationRepository) List(location string, onlyActive bool) ([]model.Navigation, error) {
	query := r.db.Model(&model.Navigation{})
	if location != "" {
		query = query.Where("location = ?", location)
	}
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	var navs []model.Navigation
	if err := query.Order("`order` asc, created_at asc").Find(&navs).Error; err != nil {
		return nil, err
	}
	return navs, nil
}

func (r *GormNavigationRepository) Reorder(location string, ids []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.Navigation{}).Where("id = ? AND location = ?", id, location).Update("order", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
