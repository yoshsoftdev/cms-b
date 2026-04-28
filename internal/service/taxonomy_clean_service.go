package service

import (
	"strings"

	"cms/internal/repository"
)

type TaxonomyCleanService struct {
	repo *repository.TaxonomyCleanRepository
}

func NewTaxonomyCleanService(repo *repository.TaxonomyCleanRepository) *TaxonomyCleanService {
	return &TaxonomyCleanService{repo: repo}
}

func (s *TaxonomyCleanService) List(locale, termType, search string, limit int) ([]repository.TaxonomyCleanItem, error) {
	return s.repo.List(strings.TrimSpace(locale), strings.TrimSpace(termType), strings.TrimSpace(search), limit)
}
