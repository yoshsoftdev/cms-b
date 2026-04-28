package service

import (
	"cms/internal/model"
	"cms/internal/repository"
	"strings"
)

type AuthorService struct{ repo *repository.AuthorRepository }

func NewAuthorService(repo *repository.AuthorRepository) *AuthorService { return &AuthorService{repo: repo} }

func (s *AuthorService) List(search string) ([]model.Author, error) {
	return s.repo.List(search)
}

type UpsertAuthorRequest struct {
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	AltText string `json:"alt_text"`
}

func (s *AuthorService) Create(req UpsertAuthorRequest) (*model.Author, error) {
	input := model.Author{
		Name:    strings.TrimSpace(req.Name),
		Photo:   strings.TrimSpace(req.Photo),
		AltText: strings.TrimSpace(req.AltText),
	}
	return s.repo.Create(input)
}

func (s *AuthorService) Update(id string, req UpsertAuthorRequest) (*model.Author, error) {
	input := model.Author{
		Name:    strings.TrimSpace(req.Name),
		Photo:   strings.TrimSpace(req.Photo),
		AltText: strings.TrimSpace(req.AltText),
	}
	return s.repo.Update(id, input)
}

func (s *AuthorService) Delete(id string) error {
	return s.repo.Delete(id)
}
