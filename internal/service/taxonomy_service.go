package service

import (
	"strings"

	"cms/internal/model"
	"cms/internal/repository"
)

type TaxonomyService struct {
	repo *repository.TaxonomyRepository
}

func NewTaxonomyService(repo *repository.TaxonomyRepository) *TaxonomyService {
	return &TaxonomyService{repo: repo}
}

func (s *TaxonomyService) List(termType model.TaxonomyType, search string, limit int) ([]model.TaxonomyTerm, error) {
	return s.repo.List(termType, strings.TrimSpace(search), limit)
}

func (s *TaxonomyService) UpsertKeywordsFromCSV(csv string) error {
	items := strings.Split(csv, ",")
	return s.repo.UpsertNames(model.TaxonomyKeyword, items)
}

func (s *TaxonomyService) UpsertCategory(category string) error {
	if strings.TrimSpace(category) == "" {
		return nil
	}
	return s.repo.UpsertNames(model.TaxonomyCategory, []string{category})
}
