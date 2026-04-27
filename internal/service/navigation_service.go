package service

import (
	"cms/internal/model"
	"cms/internal/repository"
	"cms/internal/utils"
)

type NavigationService struct {
	repo repository.NavigationRepository
}

func NewNavigationService(repo repository.NavigationRepository) *NavigationService {
	return &NavigationService{repo: repo}
}

func (s *NavigationService) Create(nav *model.Navigation) error           { return s.repo.Create(nav) }
func (s *NavigationService) Update(nav *model.Navigation) error           { return s.repo.Update(nav) }
func (s *NavigationService) Delete(id string) error                       { return s.repo.Delete(id) }
func (s *NavigationService) GetByID(id string) (*model.Navigation, error) { return s.repo.GetByID(id) }
func (s *NavigationService) Reorder(location string, ids []string) error {
	return s.repo.Reorder(location, ids)
}

func (s *NavigationService) Tree(location string, onlyActive bool) ([]utils.NavigationTree, error) {
	items, err := s.repo.List(location, onlyActive)
	if err != nil {
		return nil, err
	}
	return utils.BuildNavigationTree(items), nil
}
