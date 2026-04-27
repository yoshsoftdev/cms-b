package service

import (
	"fmt"
	"html"

	"cms/internal/model"
	"cms/internal/repository"
	"cms/internal/utils"
)

type PageService struct{ repo *repository.PageRepository }

func NewPageService(repo *repository.PageRepository) *PageService { return &PageService{repo: repo} }

func (s *PageService) Create(page *model.Page) error {
	page.Content = html.EscapeString(page.Content)
	base := utils.Slugify(page.Title)
	for i := 0; i < 100; i++ {
		slug := base
		if i > 0 {
			slug = fmt.Sprintf("%s-%d", base, i)
		}
		exists, err := s.repo.ExistsBySlug(slug, "")
		if err != nil {
			return err
		}
		if !exists {
			page.Slug = slug
			return s.repo.Create(page)
		}
	}
	return fmt.Errorf("failed to generate unique slug")
}
func (s *PageService) Update(id string, page *model.Page) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	page.Content = html.EscapeString(page.Content)
	base := utils.Slugify(page.Title)
	for i := 0; i < 100; i++ {
		slug := base
		if i > 0 {
			slug = fmt.Sprintf("%s-%d", base, i)
		}
		exists, err := s.repo.ExistsBySlug(slug, id)
		if err != nil {
			return err
		}
		if !exists {
			page.Slug = slug
			return s.repo.Update(page)
		}
	}
	return fmt.Errorf("failed to generate unique slug")
}
func (s *PageService) Delete(id string) error                     { return s.repo.Delete(id) }
func (s *PageService) GetByID(id string) (*model.Page, error)     { return s.repo.GetByID(id) }
func (s *PageService) GetBySlug(slug string) (*model.Page, error) { return s.repo.GetBySlug(slug) }
func (s *PageService) List(filter repository.ContentFilter) ([]model.Page, int64, error) {
	return s.repo.List(filter)
}
