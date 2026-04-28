package service

import (
	"cms/internal/model"
	"cms/internal/repository"
	"strings"
)

type MediaService struct{ repo *repository.MediaRepository }

func NewMediaService(repo *repository.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

type UpsertMediaRequest struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	AltText string `json:"alt_text"`
}

func (s *MediaService) List(search string) ([]model.Media, error) {
	return s.repo.List(search)
}

func (s *MediaService) Create(req UpsertMediaRequest) (*model.Media, error) {
	input := model.Media{
		Name:    strings.TrimSpace(req.Name),
		URL:     strings.TrimSpace(req.URL),
		AltText: strings.TrimSpace(req.AltText),
	}
	return s.repo.Create(input)
}

func (s *MediaService) Update(id string, req UpsertMediaRequest) (*model.Media, error) {
	input := model.Media{
		Name:    strings.TrimSpace(req.Name),
		URL:     strings.TrimSpace(req.URL),
		AltText: strings.TrimSpace(req.AltText),
	}
	return s.repo.Update(id, input)
}

func (s *MediaService) Delete(id string) error {
	return s.repo.Delete(id)
}
